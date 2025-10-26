package player

import (
	"fmt"
	"os/exec"
	"runtime"
)

// PlayWithMPV plays a stream URL using MPV player with optional subtitles
func PlayWithMPV(streamURL string, title string, subtitleURLs []string) error {
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

	// Build MPV command with useful options
	args := []string{
		streamURL,
		fmt.Sprintf("--title=%s", title),
		"--force-window=immediate",
		"--keep-open=yes",
		"--osd-level=1",
		"--slang=en,eng,english", // Prefer English subtitles
		"--sub-auto=all",          // Load all available subtitles
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
		fmt.Printf("✓ English subtitles enabled by default (Press V to toggle)\n")
	} else {
		fmt.Printf("⚠ No subtitles available for this content\n")
	}

	cmd := exec.Command(mpvPath, args...)
	
	// Start MPV in background
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start MPV: %v", err)
	}

	return nil
}

// CheckMPVInstalled checks if MPV is installed on the system
func CheckMPVInstalled() bool {
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

