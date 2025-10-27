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
func PlayWithPlayer(streamURL string, title string, subtitleURLs []string, audioTrackURLs []string, onEnd OnPlaybackEndCallback) error {
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
		return launchMPV(player.Executable, streamURL, title, subtitleURLs, audioTrackURLs, onEnd)
	case "vlc":
		return launchVLC(player.Executable, streamURL, title, subtitleURLs, audioTrackURLs, onEnd)
	case "mpc-hc":
		return launchMPCHC(player.Executable, streamURL, title, subtitleURLs, audioTrackURLs, onEnd)
	case "potplayer":
		return launchPotPlayer(player.Executable, streamURL, title, subtitleURLs, audioTrackURLs, onEnd)
	default:
		return fmt.Errorf("unsupported video player: %s", player.ID)
	}
}

// launchMPV launches MPV player
func launchMPV(exePath, streamURL, title string, subtitleURLs []string, audioTrackURLs []string, onEnd OnPlaybackEndCallback) error {
	userSettings := settings.Get()

	args := []string{
		streamURL,
		fmt.Sprintf("--title=%s", title),
		"--force-window=immediate",
		"--osd-level=1",
		fmt.Sprintf("--slang=%s,%s,eng,english", userSettings.SubtitleLanguage, userSettings.SubtitleLanguage+"g"),
		fmt.Sprintf("--alang=%s,%s,eng,english", userSettings.AudioLanguage, userSettings.AudioLanguage+"g"),
		"--sub-auto=all",
	}

	// Add keep-open option based on user setting
	if userSettings.AutoClosePlayer {
		args = append(args, "--keep-open=no")
	} else {
		args = append(args, "--keep-open=yes")
	}

	// Add fullscreen option if enabled
	if userSettings.Fullscreen {
		args = append(args, "--fullscreen")
	}

	// Add external audio tracks if provided
	if len(audioTrackURLs) > 0 {
		fmt.Printf("✓ Loading %d external audio track(s) into MPV\n", len(audioTrackURLs))
		for i, audioURL := range audioTrackURLs {
			args = append(args, fmt.Sprintf("--audio-file=%s", audioURL))
			fmt.Printf("  Audio track %d: %s\n", i+1, audioURL)
		}
		langName := getLanguageName(userSettings.AudioLanguage)
		fmt.Printf("✓ External audio tracks loaded (preferred: %s) - Press # to cycle audio tracks\n", langName)
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
func launchVLC(exePath, streamURL, title string, subtitleURLs []string, audioTrackURLs []string, onEnd OnPlaybackEndCallback) error {
	userSettings := settings.Get()

	args := []string{
		streamURL,
		fmt.Sprintf("--meta-title=%s", title),
		fmt.Sprintf("--sub-language=%s", userSettings.SubtitleLanguage),
		fmt.Sprintf("--audio-language=%s", userSettings.AudioLanguage),
	}

	// Add auto-close flag if enabled
	if userSettings.AutoClosePlayer {
		args = append(args, "--play-and-exit")
	}

	// Add fullscreen option if enabled
	if userSettings.Fullscreen {
		args = append(args, "--fullscreen")
	}

	var tempFiles []string

	// Download audio track files to temp directory if provided
	if len(audioTrackURLs) > 0 {
		fmt.Printf("✓ Loading %d external audio track(s) into VLC\n", len(audioTrackURLs))
		for i, audioURL := range audioTrackURLs {
			// Download audio to temp file
			tempFile, err := downloadAudioToTemp(audioURL, i)
			if err != nil {
				fmt.Printf("  ⚠ Warning: Failed to download audio track %d: %v\n", i+1, err)
				continue
			}
			tempFiles = append(tempFiles, tempFile)
			// VLC uses --input-slave for additional audio tracks
			args = append(args, fmt.Sprintf("--input-slave=%s", tempFile))
			fmt.Printf("  Audio track %d: %s (downloaded)\n", i+1, audioURL)
		}
		langName := getLanguageName(userSettings.AudioLanguage)
		fmt.Printf("✓ External audio tracks loaded (preferred: %s) - Press B to cycle audio tracks\n", langName)
	}

	// Download subtitle files to temp directory if provided
	if len(subtitleURLs) > 0 {
		fmt.Printf("✓ Loading %d subtitle track(s) into VLC\n", len(subtitleURLs))
		for i, subURL := range subtitleURLs {
			// Download subtitle to temp file
			tempFile, err := downloadSubtitleToTemp(subURL, i)
			if err != nil {
				fmt.Printf("  ⚠ Warning: Failed to download subtitle %d: %v\n", i+1, err)
				continue
			}
			tempFiles = append(tempFiles, tempFile)
			args = append(args, fmt.Sprintf("--sub-file=%s", tempFile))
			fmt.Printf("  Subtitle %d: %s (downloaded)\n", i+1, subURL)
		}
		langName := getLanguageName(userSettings.SubtitleLanguage)
		fmt.Printf("✓ Subtitles enabled (preferred: %s) - Press V to toggle\n", langName)
	}

	cmd := exec.Command(exePath, args...)
	if err := cmd.Start(); err != nil {
		// Clean up temp files on error
		cleanupTempFiles(tempFiles)
		return fmt.Errorf("failed to start VLC: %v", err)
	}

	// Clean up temp files after player exits
	if onEnd != nil {
		go func() {
			cmd.Wait()
			cleanupTempFiles(tempFiles)
			onEnd()
		}()
	} else {
		go func() {
			cmd.Wait()
			cleanupTempFiles(tempFiles)
		}()
	}

	return nil
}

// launchMPCHC launches MPC-HC player
func launchMPCHC(exePath, streamURL, title string, subtitleURLs []string, audioTrackURLs []string, onEnd OnPlaybackEndCallback) error {
	userSettings := settings.Get()
	args := []string{streamURL}

	// Add auto-close flag if enabled
	if userSettings.AutoClosePlayer {
		args = append(args, "/close")
	}

	// Add fullscreen option if enabled
	if userSettings.Fullscreen {
		args = append(args, "/fullscreen")
	}

	var tempFiles []string
	
	// Note: MPC-HC has limited support for external audio tracks via command line
	if len(audioTrackURLs) > 0 {
		fmt.Printf("⚠ Note: MPC-HC has limited support for external audio tracks\n")
		fmt.Printf("   Please load audio tracks manually in the player if needed\n")
	}

	// Download subtitle file if provided (MPC-HC can only load one subtitle file via command line)
	if len(subtitleURLs) > 0 {
		fmt.Printf("✓ Loading subtitle into MPC-HC\n")
		// Download first subtitle to temp file
		tempFile, err := downloadSubtitleToTemp(subtitleURLs[0], 0)
		if err != nil {
			fmt.Printf("  ⚠ Warning: Failed to download subtitle: %v\n", err)
		} else {
			tempFiles = append(tempFiles, tempFile)
			args = append(args, fmt.Sprintf("/sub %s", tempFile))
		}
		if len(subtitleURLs) > 1 {
			fmt.Printf("⚠ Note: MPC-HC only supports one subtitle file via command line\n")
		}
	}

	cmd := exec.Command(exePath, args...)
	if err := cmd.Start(); err != nil {
		cleanupTempFiles(tempFiles)
		return fmt.Errorf("failed to start MPC-HC: %v", err)
	}

	// Clean up temp files after player exits
	if onEnd != nil {
		go func() {
			cmd.Wait()
			cleanupTempFiles(tempFiles)
			onEnd()
		}()
	} else {
		go func() {
			cmd.Wait()
			cleanupTempFiles(tempFiles)
		}()
	}

	return nil
}

// launchPotPlayer launches PotPlayer
func launchPotPlayer(exePath, streamURL, title string, subtitleURLs []string, audioTrackURLs []string, onEnd OnPlaybackEndCallback) error {
	userSettings := settings.Get()
	args := []string{streamURL}

	// Add auto-close flag if enabled
	if userSettings.AutoClosePlayer {
		args = append(args, "/close")
	}

	// Add fullscreen option if enabled
	if userSettings.Fullscreen {
		args = append(args, "/fullscreen")
	}

	var tempFiles []string

	// Download audio track files if provided
	if len(audioTrackURLs) > 0 {
		fmt.Printf("✓ Loading %d external audio track(s) into PotPlayer\n", len(audioTrackURLs))
		for i, audioURL := range audioTrackURLs {
			// Download audio to temp file
			tempFile, err := downloadAudioToTemp(audioURL, i)
			if err != nil {
				fmt.Printf("  ⚠ Warning: Failed to download audio track %d: %v\n", i+1, err)
				continue
			}
			tempFiles = append(tempFiles, tempFile)
			// PotPlayer can load external audio tracks
			args = append(args, fmt.Sprintf("/add=%s", tempFile))
			fmt.Printf("  Audio track %d: %s (downloaded)\n", i+1, audioURL)
		}
		langName := getLanguageName(userSettings.AudioLanguage)
		fmt.Printf("✓ External audio tracks loaded (preferred: %s)\n", langName)
	}

	// Download subtitle files if provided
	if len(subtitleURLs) > 0 {
		fmt.Printf("✓ Loading %d subtitle track(s) into PotPlayer\n", len(subtitleURLs))
		for i, subURL := range subtitleURLs {
			// Download subtitle to temp file
			tempFile, err := downloadSubtitleToTemp(subURL, i)
			if err != nil {
				fmt.Printf("  ⚠ Warning: Failed to download subtitle %d: %v\n", i+1, err)
				continue
			}
			tempFiles = append(tempFiles, tempFile)
			args = append(args, fmt.Sprintf("/sub=%s", tempFile))
		}
	}

	cmd := exec.Command(exePath, args...)
	if err := cmd.Start(); err != nil {
		cleanupTempFiles(tempFiles)
		return fmt.Errorf("failed to start PotPlayer: %v", err)
	}

	// Clean up temp files after player exits
	if onEnd != nil {
		go func() {
			cmd.Wait()
			cleanupTempFiles(tempFiles)
			onEnd()
		}()
	} else {
		go func() {
			cmd.Wait()
			cleanupTempFiles(tempFiles)
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

// downloadAudioToTemp downloads an audio file to a temporary location
func downloadAudioToTemp(url string, index int) (string, error) {
	// Download audio content
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download audio: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to download audio: HTTP %d", resp.StatusCode)
	}

	// Create temp directory for audio files
	tempDir := filepath.Join(os.TempDir(), "moviestream_audio")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Determine file extension from URL or content-type
	ext := ".aac"
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		switch contentType {
		case "audio/mpeg", "audio/mp3":
			ext = ".mp3"
		case "audio/ogg":
			ext = ".ogg"
		case "audio/opus":
			ext = ".opus"
		case "audio/aac":
			ext = ".aac"
		}
	}

	tempFile := filepath.Join(tempDir, fmt.Sprintf("moviestream_audio_%d%s", index, ext))
	
	file, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer file.Close()

	// Write content to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(tempFile)
		return "", fmt.Errorf("failed to write audio file: %v", err)
	}

	return tempFile, nil
}

// cleanupTempSubtitles removes temporary subtitle files (legacy, kept for compatibility)
func cleanupTempSubtitles(files []string) {
	cleanupTempFiles(files)
}

// cleanupTempFiles removes temporary files (subtitles, audio, etc.)
func cleanupTempFiles(files []string) {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			// Silently ignore cleanup errors
			fmt.Printf("Warning: Failed to cleanup temp file %s: %v\n", file, err)
		}
	}
}

