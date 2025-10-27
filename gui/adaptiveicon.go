package gui

import (
	"os"
	"path/filepath"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
)

// GetAdaptiveIcon returns an icon based on the system theme
// Dark icon for light mode, light icon for dark mode
func GetAdaptiveIcon() fyne.Resource {
	isSystemDark := IsSystemDarkMode()
	
	if isSystemDark {
		// System is in dark mode, use a light/white icon
		icon, err := loadIconFromFile("Icon-dark.png")
		if err != nil {
			// Fallback to default icon
			return theme.MediaVideoIcon()
		}
		return icon
	} else {
		// System is in light mode, use a dark icon
		icon, err := loadIconFromFile("Icon-light.png")
		if err != nil {
			// Fallback to default icon
			return theme.MediaVideoIcon()
		}
		return icon
	}
}

// loadIconFromFile loads an icon from the filesystem
func loadIconFromFile(filename string) (fyne.Resource, error) {
	// Get the executable directory
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)
	
	// Try loading from executable directory first
	iconPath := filepath.Join(exeDir, filename)
	if _, err := os.Stat(iconPath); err == nil {
		uri := storage.NewFileURI(iconPath)
		return storage.LoadResourceFromURI(uri)
	}
	
	// Try loading from current working directory
	if _, err := os.Stat(filename); err == nil {
		uri := storage.NewFileURI(filename)
		return storage.LoadResourceFromURI(uri)
	}
	
	return nil, os.ErrNotExist
}

