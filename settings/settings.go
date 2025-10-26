package settings

import (
	"fyne.io/fyne/v2"
)

// Settings holds all user preferences
type Settings struct {
	SubtitleLanguage string // e.g., "en", "es", "fr", "de", etc.
	AudioLanguage    string // e.g., "en", "es", "fr", "de", etc.
	AutoNext         bool   // Auto-play next episode for TV shows
	VideoPlayer      string // Video player ID (e.g., "mpv", "vlc", "mpc-hc", "potplayer")
}

var (
	app               fyne.App
	currentSettings   *Settings
)

// Initialize sets up the settings system
func Initialize(fyneApp fyne.App) {
	app = fyneApp
	currentSettings = Load()
}

// Load retrieves settings from storage
func Load() *Settings {
	if app == nil {
		// Return defaults if not initialized
		return GetDefaults()
	}

	prefs := app.Preferences()
	
	settings := &Settings{
		SubtitleLanguage: prefs.StringWithFallback("subtitle_language", "en"),
		AudioLanguage:    prefs.StringWithFallback("audio_language", "en"),
		AutoNext:         prefs.BoolWithFallback("auto_next", false),
		VideoPlayer:      prefs.StringWithFallback("video_player", "mpv"),
	}
	
	return settings
}

// Save stores settings to persistent storage
func Save(settings *Settings) {
	if app == nil {
		return
	}

	prefs := app.Preferences()
	prefs.SetString("subtitle_language", settings.SubtitleLanguage)
	prefs.SetString("audio_language", settings.AudioLanguage)
	prefs.SetBool("auto_next", settings.AutoNext)
	prefs.SetString("video_player", settings.VideoPlayer)
	
	currentSettings = settings
}

// Get returns the current settings
func Get() *Settings {
	if currentSettings == nil {
		currentSettings = Load()
	}
	return currentSettings
}

// GetDefaults returns default settings
func GetDefaults() *Settings {
	return &Settings{
		SubtitleLanguage: "en",
		AudioLanguage:    "en",
		AutoNext:         false,
		VideoPlayer:      "mpv",
	}
}

// GetLanguageOptions returns available language options
func GetLanguageOptions() []string {
	return []string{
		"en - English",
		"es - Spanish",
		"fr - French",
		"de - German",
		"it - Italian",
		"pt - Portuguese",
		"ja - Japanese",
		"ko - Korean",
		"zh - Chinese",
		"ar - Arabic",
		"ru - Russian",
		"hi - Hindi",
	}
}

// GetLanguageCode extracts the language code from the display string
func GetLanguageCode(displayString string) string {
	if len(displayString) >= 2 {
		return displayString[:2]
	}
	return "en"
}

// GetLanguageDisplayString finds the display string for a language code
func GetLanguageDisplayString(code string) string {
	options := GetLanguageOptions()
	for _, opt := range options {
		if len(opt) >= 2 && opt[:2] == code {
			return opt
		}
	}
	return "en - English"
}

