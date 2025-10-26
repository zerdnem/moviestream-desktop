package player

import (
	"fmt"
	"moviestream-gui/settings"
	"os/exec"
	"runtime"
)

// OnPlaybackEndCallback is a function that gets called when playback ends
type OnPlaybackEndCallback func()

// PlayWithMPV plays a stream URL using the user's selected video player with optional subtitles
func PlayWithMPV(streamURL string, title string, subtitleURLs []string) error {
	return PlayWithMPVAndCallback(streamURL, title, subtitleURLs, nil)
}

// PlayWithMPVAndCallback plays a stream URL using the user's selected video player with optional callback when playback ends
func PlayWithMPVAndCallback(streamURL string, title string, subtitleURLs []string, onEnd OnPlaybackEndCallback) error {
	// Use the new player launcher that respects user's player choice
	return PlayWithPlayer(streamURL, title, subtitleURLs, onEnd)
}

// PlayWithMPVLegacy is the original MPV-specific implementation (kept for reference)
func PlayWithMPVLegacy(streamURL string, title string, subtitleURLs []string, onEnd OnPlaybackEndCallback) error {
	// Check if MPV is available
	mpvPath := "mpv"
	
	// On Windows, try common installation paths
	if runtime.GOOS == "windows" {
		// Try to find mpv in PATH first
		_, err := exec.LookPath("mpv")
		if err != nil {
			// Try common installation paths
			commonPaths := []string{
				"C:\\Program Files\\mpv\\mpv.exe",
				"C:\\Program Files (x86)\\mpv\\mpv.exe",
				"C:\\mpv\\mpv.exe",
			}
			
			found := false
			for _, path := range commonPaths {
				if _, err := exec.LookPath(path); err == nil {
					mpvPath = path
					found = true
					break
				}
			}
			
			if !found {
				return fmt.Errorf("MPV player not found. Please install MPV from https://mpv.io/")
			}
		}
	}

	// Get user settings
	userSettings := settings.Get()
	
	// Build MPV command with useful options
	args := []string{
		streamURL,
		fmt.Sprintf("--title=%s", title),
		"--force-window=immediate",
		"--keep-open=yes",
		"--osd-level=1",
		fmt.Sprintf("--slang=%s,%s,eng,english", userSettings.SubtitleLanguage, userSettings.SubtitleLanguage+"g"), // Prefer user's subtitle language
		fmt.Sprintf("--alang=%s,%s,eng,english", userSettings.AudioLanguage, userSettings.AudioLanguage+"g"),       // Prefer user's audio language
		"--sub-auto=all",                                                                                            // Load all available subtitles
	}

	// Add subtitle files if provided
	if len(subtitleURLs) > 0 {
		fmt.Printf("✓ Loading %d subtitle track(s) into MPV\n", len(subtitleURLs))
		for i, subURL := range subtitleURLs {
			args = append(args, fmt.Sprintf("--sub-file=%s", subURL))
			fmt.Printf("  Subtitle %d: %s\n", i+1, subURL)
		}
		// Enable subtitles by default
		args = append(args, "--sid=1") // Select first subtitle track
		langName := map[string]string{"en": "English", "es": "Spanish", "fr": "French", "de": "German", "it": "Italian", "pt": "Portuguese", "ja": "Japanese", "ko": "Korean", "zh": "Chinese", "ar": "Arabic", "ru": "Russian", "hi": "Hindi"}
		preferredLang := langName[userSettings.SubtitleLanguage]
		if preferredLang == "" {
			preferredLang = userSettings.SubtitleLanguage
		}
		fmt.Printf("✓ Subtitles enabled (preferred: %s) - Press V to toggle\n", preferredLang)
	} else {
		fmt.Printf("⚠ No subtitles available for this content\n")
	}

	cmd := exec.Command(mpvPath, args...)
	
	// Start MPV in background
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start MPV: %v", err)
	}

	// If callback is provided, wait for MPV to finish in a goroutine
	if onEnd != nil {
		go func() {
			cmd.Wait() // Wait for MPV to exit
			onEnd()    // Call the callback
		}()
	}

	return nil
}

// CheckMPVInstalled checks if any supported video player is installed on the system
func CheckMPVInstalled() bool {
	installedPlayers := GetInstalledPlayers()
	return len(installedPlayers) > 0
}

// CheckMPVInstalledLegacy checks if MPV specifically is installed
func CheckMPVInstalledLegacy() bool {
	_, err := exec.LookPath("mpv")
	if err == nil {
		return true
	}

	// On Windows, check common paths
	if runtime.GOOS == "windows" {
		commonPaths := []string{
			"C:\\Program Files\\mpv\\mpv.exe",
			"C:\\Program Files (x86)\\mpv\\mpv.exe",
			"C:\\mpv\\mpv.exe",
		}
		
		for _, path := range commonPaths {
			if _, err := exec.LookPath(path); err == nil {
				return true
			}
		}
	}

	return false
}

