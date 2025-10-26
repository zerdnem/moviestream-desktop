package player

import (
	"fmt"
	"io"
	"moviestream-gui/settings"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// PlayWithPlayer plays a stream using the user's selected video player
func PlayWithPlayer(streamURL string, title string, subtitleURLs []string, onEnd OnPlaybackEndCallback) error {
	userSettings := settings.Get()
	playerID := userSettings.VideoPlayer

	// Get the selected player
	player := GetPlayerByID(playerID)
	if player == nil || !player.IsInstalled {
		// Fallback to default player
		player = GetDefaultPlayer()
		if !player.IsInstalled {
			return fmt.Errorf("no video player found. Please install MPV, VLC, or another supported player")
		}
	}

	// Launch the appropriate player
	switch player.ID {
	case "mpv":
		return launchMPV(player.Executable, streamURL, title, subtitleURLs, onEnd)
	case "vlc":
		return launchVLC(player.Executable, streamURL, title, subtitleURLs, onEnd)
	case "mpc-hc":
		return launchMPCHC(player.Executable, streamURL, title, subtitleURLs, onEnd)
	case "potplayer":
		return launchPotPlayer(player.Executable, streamURL, title, subtitleURLs, onEnd)
	default:
		return fmt.Errorf("unsupported video player: %s", player.ID)
	}
}

// launchMPV launches MPV player
func launchMPV(exePath, streamURL, title string, subtitleURLs []string, onEnd OnPlaybackEndCallback) error {
	userSettings := settings.Get()

	args := []string{
		streamURL,
		fmt.Sprintf("--title=%s", title),
		"--force-window=immediate",
		"--keep-open=yes",
		"--osd-level=1",
		fmt.Sprintf("--slang=%s,%s,eng,english", userSettings.SubtitleLanguage, userSettings.SubtitleLanguage+"g"),
		fmt.Sprintf("--alang=%s,%s,eng,english", userSettings.AudioLanguage, userSettings.AudioLanguage+"g"),
		"--sub-auto=all",
	}

	// Add subtitle files if provided
	if len(subtitleURLs) > 0 {
		fmt.Printf("✓ Loading %d subtitle track(s) into MPV\n", len(subtitleURLs))
		for i, subURL := range subtitleURLs {
			args = append(args, fmt.Sprintf("--sub-file=%s", subURL))
			fmt.Printf("  Subtitle %d: %s\n", i+1, subURL)
		}
		args = append(args, "--sid=1")
		langName := getLanguageName(userSettings.SubtitleLanguage)
		fmt.Printf("✓ Subtitles enabled (preferred: %s) - Press V to toggle\n", langName)
	} else {
		fmt.Printf("⚠ No subtitles available for this content\n")
	}

	cmd := exec.Command(exePath, args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start MPV: %v", err)
	}

	if onEnd != nil {
		go func() {
			cmd.Wait()
			onEnd()
		}()
	}

	return nil
}

// launchVLC launches VLC player
func launchVLC(exePath, streamURL, title string, subtitleURLs []string, onEnd OnPlaybackEndCallback) error {
	userSettings := settings.Get()

	args := []string{
		streamURL,
		fmt.Sprintf("--meta-title=%s", title),
		fmt.Sprintf("--sub-language=%s", userSettings.SubtitleLanguage),
		fmt.Sprintf("--audio-language=%s", userSettings.AudioLanguage),
	}

	// Download subtitle files to temp directory if provided
	var tempSubFiles []string
	if len(subtitleURLs) > 0 {
		fmt.Printf("✓ Loading %d subtitle track(s) into VLC\n", len(subtitleURLs))
		for i, subURL := range subtitleURLs {
			// Download subtitle to temp file
			tempFile, err := downloadSubtitleToTemp(subURL, i)
			if err != nil {
				fmt.Printf("  ⚠ Warning: Failed to download subtitle %d: %v\n", i+1, err)
				continue
			}
			tempSubFiles = append(tempSubFiles, tempFile)
			args = append(args, fmt.Sprintf("--sub-file=%s", tempFile))
			fmt.Printf("  Subtitle %d: %s (downloaded)\n", i+1, subURL)
		}
		langName := getLanguageName(userSettings.SubtitleLanguage)
		fmt.Printf("✓ Subtitles enabled (preferred: %s) - Press V to toggle\n", langName)
	}

	cmd := exec.Command(exePath, args...)
	if err := cmd.Start(); err != nil {
		// Clean up temp files on error
		cleanupTempSubtitles(tempSubFiles)
		return fmt.Errorf("failed to start VLC: %v", err)
	}

	// Clean up temp subtitle files after player exits
	if onEnd != nil {
		go func() {
			cmd.Wait()
			cleanupTempSubtitles(tempSubFiles)
			onEnd()
		}()
	} else {
		go func() {
			cmd.Wait()
			cleanupTempSubtitles(tempSubFiles)
		}()
	}

	return nil
}

// launchMPCHC launches MPC-HC player
func launchMPCHC(exePath, streamURL, title string, subtitleURLs []string, onEnd OnPlaybackEndCallback) error {
	args := []string{streamURL}

	// Download subtitle file if provided (MPC-HC can only load one subtitle file via command line)
	var tempSubFiles []string
	if len(subtitleURLs) > 0 {
		fmt.Printf("✓ Loading subtitle into MPC-HC\n")
		// Download first subtitle to temp file
		tempFile, err := downloadSubtitleToTemp(subtitleURLs[0], 0)
		if err != nil {
			fmt.Printf("  ⚠ Warning: Failed to download subtitle: %v\n", err)
		} else {
			tempSubFiles = append(tempSubFiles, tempFile)
			args = append(args, fmt.Sprintf("/sub %s", tempFile))
		}
		if len(subtitleURLs) > 1 {
			fmt.Printf("⚠ Note: MPC-HC only supports one subtitle file via command line\n")
		}
	}

	cmd := exec.Command(exePath, args...)
	if err := cmd.Start(); err != nil {
		cleanupTempSubtitles(tempSubFiles)
		return fmt.Errorf("failed to start MPC-HC: %v", err)
	}

	// Clean up temp subtitle files after player exits
	if onEnd != nil {
		go func() {
			cmd.Wait()
			cleanupTempSubtitles(tempSubFiles)
			onEnd()
		}()
	} else {
		go func() {
			cmd.Wait()
			cleanupTempSubtitles(tempSubFiles)
		}()
	}

	return nil
}

// launchPotPlayer launches PotPlayer
func launchPotPlayer(exePath, streamURL, title string, subtitleURLs []string, onEnd OnPlaybackEndCallback) error {
	args := []string{streamURL}

	// Download subtitle files if provided
	var tempSubFiles []string
	if len(subtitleURLs) > 0 {
		fmt.Printf("✓ Loading %d subtitle track(s) into PotPlayer\n", len(subtitleURLs))
		for i, subURL := range subtitleURLs {
			// Download subtitle to temp file
			tempFile, err := downloadSubtitleToTemp(subURL, i)
			if err != nil {
				fmt.Printf("  ⚠ Warning: Failed to download subtitle %d: %v\n", i+1, err)
				continue
			}
			tempSubFiles = append(tempSubFiles, tempFile)
			args = append(args, fmt.Sprintf("/sub=%s", tempFile))
		}
	}

	cmd := exec.Command(exePath, args...)
	if err := cmd.Start(); err != nil {
		cleanupTempSubtitles(tempSubFiles)
		return fmt.Errorf("failed to start PotPlayer: %v", err)
	}

	// Clean up temp subtitle files after player exits
	if onEnd != nil {
		go func() {
			cmd.Wait()
			cleanupTempSubtitles(tempSubFiles)
			onEnd()
		}()
	} else {
		go func() {
			cmd.Wait()
			cleanupTempSubtitles(tempSubFiles)
		}()
	}

	return nil
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

// downloadSubtitleToTemp downloads a subtitle file to a temporary location
func downloadSubtitleToTemp(url string, index int) (string, error) {
	// Download subtitle content
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download subtitle: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to download subtitle: HTTP %d", resp.StatusCode)
	}

	// Create temp file
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, fmt.Sprintf("moviestream_subtitle_%d.vtt", index))
	
	file, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer file.Close()

	// Write content to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(tempFile)
		return "", fmt.Errorf("failed to write subtitle file: %v", err)
	}

	return tempFile, nil
}

// cleanupTempSubtitles removes temporary subtitle files
func cleanupTempSubtitles(files []string) {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			// Silently ignore cleanup errors
			fmt.Printf("Warning: Failed to cleanup temp subtitle %s: %v\n", file, err)
		}
	}
}

