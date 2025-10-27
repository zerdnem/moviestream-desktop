package subtitles

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Addic7edClient handles communication with Addic7ed website
type Addic7edClient struct {
	baseURL   string
	userAgent string
}

// NewAddic7edClient creates a new Addic7ed client
func NewAddic7edClient() *Addic7edClient {
	return &Addic7edClient{
		baseURL:   "https://www.addic7ed.com",
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}
}

// SearchByTitle searches for subtitles on Addic7ed
func (c *Addic7edClient) SearchByTitle(title string, language string, season, episode int) ([]SubtitleResult, error) {
	// Addic7ed requires TV show format for search
	if season == 0 || episode == 0 {
		return nil, fmt.Errorf("Addic7ed only supports TV shows (season and episode required)")
	}
	
	// Search for the show first
	searchURL := fmt.Sprintf("%s/search.php?search=%s&Submit=Search", c.baseURL, url.QueryEscape(title))
	
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %v", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", c.baseURL)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read search response: %v", err)
	}
	
	// Find show ID from search results
	showID, err := c.extractShowID(string(body), title)
	if err != nil {
		return nil, fmt.Errorf("show not found: %v", err)
	}
	
	// Now get the episode page
	episodeURL := fmt.Sprintf("%s/season/%d/%d", c.baseURL, showID, season)
	
	req, err = http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create episode request: %v", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", searchURL)
	
	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("episode request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("episode page returned status %d", resp.StatusCode)
	}
	
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read episode response: %v", err)
	}
	
	// Parse subtitles from the episode page
	return c.parseSubtitles(string(body), language, episode)
}

// extractShowID extracts the show ID from search results
func (c *Addic7edClient) extractShowID(html string, title string) (int, error) {
	// Look for show links: /show/[ID]
	re := regexp.MustCompile(`/show/(\d+)`)
	matches := re.FindAllStringSubmatch(html, -1)
	
	if len(matches) == 0 {
		return 0, fmt.Errorf("no shows found")
	}
	
	// Take the first match (usually the best match)
	showID, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return 0, fmt.Errorf("invalid show ID: %v", err)
	}
	
	return showID, nil
}

// parseSubtitles extracts subtitle information from episode page HTML
func (c *Addic7edClient) parseSubtitles(html string, language string, episode int) ([]SubtitleResult, error) {
	var results []SubtitleResult
	
	// Parse table rows for the specific episode
	// Format: <tr class="epeven completed"><td>SEASON</td><td>EPISODE</td><td>TITLE</td><td>LANGUAGE</td><td>VERSION</td>...<td><a href="/updated/LANG/FILE/VER">Download</a></td>
	
	// Split by table rows
	rows := regexp.MustCompile(`<tr class="ep[^"]*">`).Split(html, -1)
	
	for _, row := range rows {
		// Check if this row is for our episode
		// Extract all <td> content
		tdRe := regexp.MustCompile(`<td[^>]*>([^<]*(?:<[^/][^>]*>[^<]*</[^>]*>)*[^<]*)</td>`)
		tdMatches := tdRe.FindAllStringSubmatch(row, -1)
		
		if len(tdMatches) < 5 {
			continue
		}
		
		// Parse columns: [0]=season, [1]=episode, [2]=title, [3]=language, [4]=version...
		episodeNum := strings.TrimSpace(stripHTML(tdMatches[1][1]))
		episodeLang := strings.TrimSpace(stripHTML(tdMatches[3][1]))
		episodeVersion := strings.TrimSpace(stripHTML(tdMatches[4][1]))
		
		// Check if this is the episode we want
		if episodeNum != fmt.Sprintf("%d", episode) {
			continue
		}
		
		// Filter by language if specified
		langCode := c.mapLanguageToCode(episodeLang)
		if language != "" && langCode != language {
			continue
		}
		
		// Extract download link from this row
		downloadRe := regexp.MustCompile(`href="(/updated/\d+/\d+/\d+)"`)
		downloadMatch := downloadRe.FindStringSubmatch(row)
		
		if len(downloadMatch) < 2 {
			continue
		}
		
		downloadPath := downloadMatch[1]
		
		result := SubtitleResult{
			ID:           downloadPath,
			Language:     langCode,
			LanguageName: episodeLang,
			FileName:     fmt.Sprintf("addic7ed_%s.srt", sanitizeFilename(episodeVersion)),
			DownloadURL:  c.baseURL + downloadPath,
		}
		
		results = append(results, result)
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("no subtitles found for episode %d", episode)
	}
	
	return results, nil
}

// stripHTML removes HTML tags from a string
func stripHTML(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(s, "")
}

// sanitizeFilename removes characters that are invalid in filenames
func sanitizeFilename(s string) string {
	s = strings.ReplaceAll(s, "+", "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	s = strings.ReplaceAll(s, ":", "_")
	return s
}

// DownloadSubtitle downloads a subtitle file
func (c *Addic7edClient) DownloadSubtitle(downloadPath string, filename string) (string, error) {
	// downloadPath is either full URL or path like /original/[id]/[number]
	downloadURL := downloadPath
	if !strings.HasPrefix(downloadPath, "http") {
		downloadURL = c.baseURL + downloadPath
	}
	
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create download request: %v", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", c.baseURL)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}
	
	// Create temp file
	tempDir := os.TempDir()
	
	// Ensure filename has .srt extension
	if !strings.HasSuffix(strings.ToLower(filename), ".srt") {
		filename += ".srt"
	}
	
	tempFile := filepath.Join(tempDir, "moviestream_addic7ed_"+filename)
	
	file, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	
	// Write content
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(tempFile)
		return "", fmt.Errorf("failed to write file: %v", err)
	}
	
	return tempFile, nil
}

// mapLanguageToCode converts language name to code
func (c *Addic7edClient) mapLanguageToCode(langName string) string {
	mapping := map[string]string{
		"English":    "en",
		"Spanish":    "es",
		"French":     "fr",
		"German":     "de",
		"Italian":    "it",
		"Portuguese": "pt",
		"Japanese":   "ja",
		"Korean":     "ko",
		"Chinese":    "zh",
		"Arabic":     "ar",
		"Russian":    "ru",
		"Hindi":      "hi",
	}
	
	for name, code := range mapping {
		if strings.Contains(strings.ToLower(langName), strings.ToLower(name)) {
			return code
		}
	}
	
	return langName
}

