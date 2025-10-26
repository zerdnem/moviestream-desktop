package player

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// VideoPlayer represents a video player application
type VideoPlayer struct {
	Name        string // Display name
	ID          string // Unique identifier
	Executable  string // Path to executable or command name
	IsInstalled bool   // Whether the player is detected on the system
	CanStream   bool   // Whether the player can handle streaming URLs
}

// GetAvailablePlayers detects and returns all available video players on the system
func GetAvailablePlayers() []VideoPlayer {
	players := []VideoPlayer{
		{
			Name:       "MPV Player",
			ID:         "mpv",
			Executable: "mpv",
			CanStream:  true,
		},
		{
			Name:       "VLC Media Player",
			ID:         "vlc",
			Executable: "vlc",
			CanStream:  true,
		},
		{
			Name:       "MPC-HC (Media Player Classic)",
			ID:         "mpc-hc",
			Executable: "mpc-hc64.exe",
			CanStream:  true,
		},
		{
			Name:       "PotPlayer",
			ID:         "potplayer",
			Executable: "PotPlayerMini64.exe",
			CanStream:  true,
		},
	}

	// Detect which players are installed
	for i := range players {
		players[i].Executable = detectPlayer(players[i].ID, players[i].Executable)
		players[i].IsInstalled = players[i].Executable != ""
	}

	return players
}

// detectPlayer finds the executable path for a given player
func detectPlayer(playerID, defaultExe string) string {
	// First, try to find in PATH
	if path, err := exec.LookPath(defaultExe); err == nil {
		return path
	}

	// If not in PATH, check common installation directories (Windows specific)
	if runtime.GOOS == "windows" {
		return detectWindowsPlayer(playerID, defaultExe)
	}

	// For Linux/Mac, rely on PATH
	return ""
}

// detectWindowsPlayer checks common Windows installation paths
func detectWindowsPlayer(playerID, defaultExe string) string {
	var searchPaths []string

	switch playerID {
	case "mpv":
		searchPaths = []string{
			`C:\Program Files\mpv\mpv.exe`,
			`C:\Program Files (x86)\mpv\mpv.exe`,
			`C:\mpv\mpv.exe`,
			filepath.Join(os.Getenv("LOCALAPPDATA"), "mpv", "mpv.exe"),
		}
	case "vlc":
		searchPaths = []string{
			`C:\Program Files\VideoLAN\VLC\vlc.exe`,
			`C:\Program Files (x86)\VideoLAN\VLC\vlc.exe`,
			filepath.Join(os.Getenv("PROGRAMFILES"), "VideoLAN", "VLC", "vlc.exe"),
			filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "VideoLAN", "VLC", "vlc.exe"),
		}
	case "mpc-hc":
		searchPaths = []string{
			`C:\Program Files\MPC-HC\mpc-hc64.exe`,
			`C:\Program Files (x86)\MPC-HC\mpc-hc.exe`,
			`C:\Program Files\K-Lite Codec Pack\MPC-HC64\mpc-hc64.exe`,
			`C:\Program Files (x86)\K-Lite Codec Pack\MPC-HC\mpc-hc.exe`,
		}
	case "potplayer":
		searchPaths = []string{
			`C:\Program Files\DAUM\PotPlayer\PotPlayerMini64.exe`,
			`C:\Program Files (x86)\DAUM\PotPlayer\PotPlayerMini.exe`,
			filepath.Join(os.Getenv("PROGRAMFILES"), "DAUM", "PotPlayer", "PotPlayerMini64.exe"),
			filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "DAUM", "PotPlayer", "PotPlayerMini.exe"),
		}
	}

	// Check each path
	for _, path := range searchPaths {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// GetPlayerByID returns a specific player by its ID
func GetPlayerByID(playerID string) *VideoPlayer {
	players := GetAvailablePlayers()
	for _, player := range players {
		if player.ID == playerID {
			return &player
		}
	}
	return nil
}

// GetInstalledPlayers returns only the installed players
func GetInstalledPlayers() []VideoPlayer {
	allPlayers := GetAvailablePlayers()
	var installed []VideoPlayer
	for _, player := range allPlayers {
		if player.IsInstalled {
			installed = append(installed, player)
		}
	}
	return installed
}

// GetDefaultPlayer returns the first installed player (preferring MPV)
func GetDefaultPlayer() *VideoPlayer {
	// First, try to get MPV
	mpv := GetPlayerByID("mpv")
	if mpv != nil && mpv.IsInstalled {
		return mpv
	}

	// Otherwise, return first installed player
	installed := GetInstalledPlayers()
	if len(installed) > 0 {
		return &installed[0]
	}

	// Return MPV as fallback (even if not installed)
	return &VideoPlayer{
		Name:       "MPV Player",
		ID:         "mpv",
		Executable: "mpv",
		CanStream:  true,
		IsInstalled: false,
	}
}

