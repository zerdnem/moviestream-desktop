package subtitles

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// SubtitleResult represents a subtitle search result
type SubtitleResult struct {
	ID           string
	Language     string
	LanguageName string
	MovieName    string
	FileName     string
	DownloadURL  string
	Rating       float64
}

// OpenSubtitlesClient handles communication with OpenSubtitles API
type OpenSubtitlesClient struct {
	apiKey    string
	userAgent string
	baseURL   string
}

// NewOpenSubtitlesClient creates a new OpenSubtitles API client
func NewOpenSubtitlesClient() *OpenSubtitlesClient {
	return &OpenSubtitlesClient{
		apiKey:    "5CZlDmqIhLoRcalZHXItm5Thwq57MDE2",
		userAgent: "MovieStream v1.0",
		baseURL:   "https://api.opensubtitles.com/api/v1/subtitles",
	}
}

// SearchByTitle searches for subtitles by movie/show title
func (c *OpenSubtitlesClient) SearchByTitle(title string, language string, season, episode int) ([]SubtitleResult, error) {
	// Build search query
	params := url.Values{}
	params.Add("query", title)
	
	if language != "" {
		params.Add("languages", mapLanguageCode(language))
	} else {
		params.Add("languages", "en")
	}
	
	// Add TV show specific parameters
	if season > 0 && episode > 0 {
		params.Add("season_number", fmt.Sprintf("%d", season))
		params.Add("episode_number", fmt.Sprintf("%d", episode))
		params.Add("type", "episode")
	} else {
		params.Add("type", "movie")
	}
	
	searchURL := c.baseURL + "?" + params.Encode()
	
	// Create request
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	
	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse response - new API structure
	var response struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Language   string `json:"language"`
				Files      []struct {
					FileName string `json:"file_name"`
					FileID   int    `json:"file_id"`
				} `json:"files"`
				FeatureDetails struct {
					MovieName string `json:"movie_name"`
					Title     string `json:"title"`
				} `json:"feature_details"`
			} `json:"attributes"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	
	// Convert to our format
	var results []SubtitleResult
	for _, item := range response.Data {
		if len(item.Attributes.Files) == 0 {
			continue
		}
		
		fileName := item.Attributes.Files[0].FileName
		fileID := item.Attributes.Files[0].FileID
		movieName := item.Attributes.FeatureDetails.MovieName
		if movieName == "" {
			movieName = item.Attributes.FeatureDetails.Title
		}
		
		result := SubtitleResult{
			ID:           fmt.Sprintf("%d", fileID),
			Language:     item.Attributes.Language,
			LanguageName: getLanguageName(item.Attributes.Language),
			MovieName:    movieName,
			FileName:     fileName,
			DownloadURL:  fmt.Sprintf("%d", fileID), // Store just the file ID
		}
		
		results = append(results, result)
	}
	
	return results, nil
}

// SearchByIMDB searches for subtitles by IMDB ID
func (c *OpenSubtitlesClient) SearchByIMDB(imdbID string, language string, season, episode int) ([]SubtitleResult, error) {
	// Build search query
	params := url.Values{}
	params.Add("imdbid", strings.TrimPrefix(imdbID, "tt"))
	
	if language != "" {
		params.Add("sublanguageid", mapLanguageCode(language))
	} else {
		params.Add("sublanguageid", "eng")
	}
	
	// Add TV show specific parameters
	if season > 0 && episode > 0 {
		params.Add("season", fmt.Sprintf("%d", season))
		params.Add("episode", fmt.Sprintf("%d", episode))
	}
	
	searchURL := c.baseURL + "/?" + params.Encode()
	
	// Create request
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	
	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	
	// Parse response
	var rawResults []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawResults); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	
	// Convert to our format
	var results []SubtitleResult
	for _, raw := range rawResults {
		result := SubtitleResult{
			ID:           getString(raw, "IDSubtitleFile"),
			Language:     getString(raw, "SubLanguageID"),
			LanguageName: getString(raw, "LanguageName"),
			MovieName:    getString(raw, "MovieName"),
			FileName:     getString(raw, "SubFileName"),
			DownloadURL:  getString(raw, "SubDownloadLink"),
		}
		
		// Only add if we have essential fields
		if result.DownloadURL != "" && result.FileName != "" {
			results = append(results, result)
		}
	}
	
	return results, nil
}

// DownloadSubtitle downloads a subtitle file to the temp directory
func (c *OpenSubtitlesClient) DownloadSubtitle(fileID string, filename string) (string, error) {
	// First, request the download link
	downloadURL := fmt.Sprintf("https://api.opensubtitles.com/api/v1/download")
	
	// Create JSON body
	body := fmt.Sprintf(`{"file_id": %s}`, fileID)
	
	req, err := http.NewRequest("POST", downloadURL, strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Api-Key", c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")
	
	// Download subtitle
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download subtitle: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		if resp.StatusCode == 503 {
			return "", fmt.Errorf("OpenSubtitles server is temporarily busy (503). Please try again in a few moments")
		}
		if resp.StatusCode == 429 {
			return "", fmt.Errorf("rate limit reached. Please wait a moment before trying again")
		}
		return "", fmt.Errorf("download request failed with status %d", resp.StatusCode)
	}
	
	// Parse download response
	var downloadResp struct {
		Link     string `json:"link"`
		FileName string `json:"file_name"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&downloadResp); err != nil {
		return "", fmt.Errorf("failed to parse download response: %v", err)
	}
	
	// Download the actual file from the link
	fileResp, err := http.Get(downloadResp.Link)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer fileResp.Body.Close()
	
	if fileResp.StatusCode != 200 {
		return "", fmt.Errorf("file download failed with status %d", fileResp.StatusCode)
	}
	
	// Create temp file
	tempDir := os.TempDir()
	
	// Use response filename if provided, otherwise use parameter
	if downloadResp.FileName != "" {
		filename = downloadResp.FileName
	}
	
	// Ensure filename has proper extension
	if !strings.HasSuffix(strings.ToLower(filename), ".srt") &&
		!strings.HasSuffix(strings.ToLower(filename), ".vtt") {
		filename += ".srt"
	}
	
	tempFile := filepath.Join(tempDir, "moviestream_"+filename)
	
	file, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	
	// Write content
	_, err = io.Copy(file, fileResp.Body)
	if err != nil {
		os.Remove(tempFile)
		return "", fmt.Errorf("failed to write file: %v", err)
	}
	
	return tempFile, nil
}

// mapLanguageCode converts our language codes to OpenSubtitles format
func mapLanguageCode(code string) string {
	// New API uses 2-letter codes
	mapping := map[string]string{
		"en": "en",
		"es": "es",
		"fr": "fr",
		"de": "de",
		"it": "it",
		"pt": "pt",
		"ja": "ja",
		"ko": "ko",
		"zh": "zh",
		"ar": "ar",
		"ru": "ru",
		"hi": "hi",
	}
	
	if mapped, ok := mapping[code]; ok {
		return mapped
	}
	return code
}

// getLanguageName returns the full language name for a code
func getLanguageName(code string) string {
	langNames := map[string]string{
		"en": "English",
		"es": "Spanish",
		"fr": "French",
		"de": "German",
		"it": "Italian",
		"pt": "Portuguese",
		"ja": "Japanese",
		"ko": "Korean",
		"zh": "Chinese",
		"ar": "Arabic",
		"ru": "Russian",
		"hi": "Hindi",
	}
	if name, ok := langNames[code]; ok {
		return name
	}
	return code
}

// getString safely extracts a string from a map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

