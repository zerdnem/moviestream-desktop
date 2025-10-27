package gui

import (
	"runtime"
	
	"golang.org/x/sys/windows/registry"
)

// IsSystemDarkMode detects if the Windows system is using dark mode
// Returns true if dark mode, false if light mode
func IsSystemDarkMode() bool {
	if runtime.GOOS != "windows" {
		// Default to dark mode for non-Windows systems
		return true
	}
	
	// Check Windows registry for theme preference
	// HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Themes\Personalize
	// Key: AppsUseLightTheme (0 = dark, 1 = light)
	k, err := registry.OpenKey(registry.CURRENT_USER, 
		`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, 
		registry.QUERY_VALUE)
	if err != nil {
		// If we can't read the registry, assume dark mode
		return true
	}
	defer k.Close()
	
	value, _, err := k.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		// If the key doesn't exist, assume dark mode
		return true
	}
	
	// 0 = dark mode, 1 = light mode
	return value == 0
}

