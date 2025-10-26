package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	VeloraBaseURL = "https://veloratv.ru/api"
	MoviesBaseURL = "https://111movies.com"
)

// StreamInfo contains stream URL and subtitle URLs
type StreamInfo struct {
	StreamURL    string
	SubtitleURLs []SubtitleTrack
}

type SubtitleTrack struct {
	Label string
	URL   string
}

type IntroSkipData struct {
	IntroSkippable bool    `json:"introSkippable"`
	IntroEnd       float64 `json:"introEnd"`
}

type AudioTrack struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Data        string `json:"data"`
}

type StreamData struct {
	URL    string     `json:"url"`
	Tracks []VTTTrack `json:"tracks"`
}

type VTTTrack struct {
	Label string `json:"label"`
	File  string `json:"file"`
}

// GetIntroSkip gets intro skip timing information
func GetIntroSkip(tmdbID, season, episode int) (*IntroSkipData, error) {
	params := url.Values{}
	params.Add("tmdbId", fmt.Sprintf("%d", tmdbID))
	params.Add("season", fmt.Sprintf("%d", season))
	params.Add("episode", fmt.Sprintf("%d", episode))

	resp, err := http.Get(fmt.Sprintf("%s/intro-end/confirmed?%s", VeloraBaseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get intro skip data: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data IntroSkipData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// GetStreamURL extracts the stream URL and subtitles for a movie or TV episode
func GetStreamURL(tmdbID int, contentType string, season, episode int) (*StreamInfo, error) {
	var embedURL string
	
	if contentType == "tv" {
		embedURL = fmt.Sprintf("%s/tv/%d/%d/%d", MoviesBaseURL, tmdbID, season, episode)
	} else {
		embedURL = fmt.Sprintf("%s/movie/%d", MoviesBaseURL, tmdbID)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", embedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	htmlContent := string(body)
	
	fmt.Printf("DEBUG: Fetched page from %s (status: %d, size: %d bytes)\n", embedURL, resp.StatusCode, len(htmlContent))

	// Extract /to/ path from HTML
	toPath := extractToPath(htmlContent)
	if toPath == "" {
		fmt.Printf("DEBUG: Failed to extract /to/ path from HTML\n")
		fmt.Printf("DEBUG: No '/to/' path found - falling back to browser automation\n")
		
		// Fallback to browser automation (like Python's Selenium)
		streamInfo, err := GetStreamURLWithBrowser(tmdbID, contentType, season, episode)
		if err != nil {
			return nil, fmt.Errorf("failed to extract stream with browser: %v", err)
		}
		return streamInfo, nil
	}
	
	fmt.Printf("DEBUG: Successfully extracted /to/ path: %s\n", toPath)

	// Get audio tracks
	audioTracks, err := getAudioTracks(toPath)
	if err != nil {
		fmt.Printf("DEBUG: Error getting audio tracks: %v\n", err)
		return nil, fmt.Errorf("could not get audio tracks: %v", err)
	}
	if len(audioTracks) == 0 {
		fmt.Printf("DEBUG: No audio tracks returned\n")
		return nil, fmt.Errorf("no audio tracks available")
	}
	
	fmt.Printf("DEBUG: Found %d audio tracks\n", len(audioTracks))
	fmt.Printf("DEBUG: Using first track: %s - %s\n", audioTracks[0].Name, audioTracks[0].Description)

	// Get stream URL with first audio track
	streamURL, err := getStreamAndVTTs(toPath, audioTracks[0].Data)
	if err != nil {
		fmt.Printf("DEBUG: Error getting stream URL: %v\n", err)
		return nil, err
	}
	
	fmt.Printf("DEBUG: Successfully got stream URL: %s\n", streamURL)

	// Return stream info (no subtitles from /to/ method yet)
	return &StreamInfo{
		StreamURL:    streamURL,
		SubtitleURLs: []SubtitleTrack{},
	}, nil
}

// extractToPath extracts the obfuscated /to/ path from HTML
func extractToPath(html string) string {
	pattern := `/to/\d+/[0-9a-f]{8}/[0-9a-f]{40}/y/[0-9a-f]{64}/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/[A-Za-z0-9_-]+/sr`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 0 {
		fmt.Printf("DEBUG: Extracted /to/ path from HTML\n")
		return matches[0]
	}
	return ""
}

// getAudioTracks fetches available audio tracks
func getAudioTracks(toPath string) ([]AudioTrack, error) {
	url := MoviesBaseURL + toPath
	
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tracks []AudioTrack
	if err := json.Unmarshal(body, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}

// getStreamAndVTTs fetches the m3u8 stream URL
func getStreamAndVTTs(toPath, audioData string) (string, error) {
	fullPath := toPath + "/" + audioData
	url := MoviesBaseURL + fullPath

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var streamData StreamData
	if err := json.Unmarshal(body, &streamData); err != nil {
		return "", err
	}

	// Clean the URL if it has extra quotes or spaces
	streamURL := strings.TrimSpace(streamData.URL)
	streamURL = strings.Trim(streamURL, "\"")

	return streamURL, nil
}

