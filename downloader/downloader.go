package downloader

import (
	"fmt"
	"os"
	"path/filepath"
)

// DownloadStream saves the stream URL to a file for use with download tools
func DownloadStream(streamURL string, filename string, progressCallback func(downloaded, total int64)) error {
	// Create downloads directory if it doesn't exist
	downloadDir := filepath.Join(".", "downloads")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("failed to create downloads directory: %v", err)
	}

	// Save the stream URL to a text file
	urlFilename := filename + ".url.txt"
	outputPath := filepath.Join(downloadDir, urlFilename)
	
	content := fmt.Sprintf("Stream URL for %s\n\n%s\n\n", filename, streamURL)
	content += "How to download:\n\n"
	content += "Option 1 - Using ffmpeg (Recommended):\n"
	content += fmt.Sprintf("ffmpeg -i \"%s\" -c copy \"%s.mp4\"\n\n", streamURL, filename[:len(filename)-5])
	content += "Option 2 - Using yt-dlp:\n"
	content += fmt.Sprintf("yt-dlp \"%s\" -o \"%s.mp4\"\n\n", streamURL, filename[:len(filename)-5])
	content += "Option 3 - Play directly in MPV:\n"
	content += fmt.Sprintf("mpv \"%s\"\n", streamURL)
	
	err := os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to save URL: %v", err)
	}

	return nil
}

// GetDownloadPath returns the path where files are downloaded
func GetDownloadPath() string {
	return filepath.Join(".", "downloads")
}

