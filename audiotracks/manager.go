package audiotracks

import (
	"fmt"
	"io"
	"moviestream-gui/settings"
	"net/http"
	"os"
	"path/filepath"
)

// Manager handles audio track search and download from multiple sources
type Manager struct {
	// Future: Add audio track API clients here
	// For now, we'll support manual audio track URLs
}

// NewManager creates a new audio track manager
func NewManager() *Manager {
	return &Manager{}
}

// SearchAudioTracks searches for audio tracks across available sources
// Note: Unlike subtitles, audio tracks are less commonly available from public APIs
// This function provides a framework for future integration with audio track sources
func (m *Manager) SearchAudioTracks(title string, tmdbID int, season, episode int) ([]AudioTrackResult, error) {
	userSettings := settings.Get()
	language := userSettings.AudioLanguage
	
	fmt.Printf("Searching for audio tracks (language: %s)...\n", language)
	
	// TODO: Implement audio track search from various sources
	// Possible sources:
	// 1. Stream metadata (some streams have multiple audio tracks embedded)
	// 2. External audio track repositories (if they exist)
	// 3. User-provided URLs
	
	// For now, return empty results
	// Users can manually add audio tracks via URL
	return nil, fmt.Errorf("automatic audio track search not yet implemented. Please add audio tracks manually via URL")
}

// DownloadAudioTrack downloads an audio track from a URL
func (m *Manager) DownloadAudioTrack(url string, filename string) (string, error) {
	fmt.Printf("Downloading audio track from: %s\n", url)
	
	// Download the audio file
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download audio track: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}
	
	// Create temp directory if it doesn't exist
	tempDir := filepath.Join(os.TempDir(), "moviestream_audio")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}
	
	// Ensure filename has proper extension
	if filename == "" {
		filename = "audio_track.aac"
	}
	
	tempFile := filepath.Join(tempDir, filename)
	
	// Create the file
	file, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	
	// Copy content
	written, err := io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(tempFile)
		return "", fmt.Errorf("failed to write file: %v", err)
	}
	
	fmt.Printf("✓ Downloaded audio track: %s (%.2f MB)\n", filename, float64(written)/(1024*1024))
	
	return tempFile, nil
}

// ValidateAudioFile checks if a file exists and is readable
func (m *Manager) ValidateAudioFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("audio file not found: %v", err)
	}
	
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file")
	}
	
	if info.Size() == 0 {
		return fmt.Errorf("audio file is empty")
	}
	
	// Check if file is readable
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot read audio file: %v", err)
	}
	file.Close()
	
	return nil
}

// GetLanguageName returns the full language name for a code
func GetLanguageName(code string) string {
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

// CleanupTempAudioFiles removes temporary audio files
func CleanupTempAudioFiles() error {
	tempDir := filepath.Join(os.TempDir(), "moviestream_audio")
	if err := os.RemoveAll(tempDir); err != nil {
		return fmt.Errorf("failed to cleanup temp audio files: %v", err)
	}
	return nil
}

