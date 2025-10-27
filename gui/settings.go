package gui

import (
	"fmt"
	"moviestream-gui/player"
	"moviestream-gui/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowSettingsDialog displays the settings dialog
func ShowSettingsDialog() {
	currentSettings := settings.Get()
	
	// Modern video player section header
	videoPlayerLabel := CreateHeader("Video Player")
	
	// Get all available players
	allPlayers := player.GetAvailablePlayers()
	installedPlayers := player.GetInstalledPlayers()
	
	var playerOptions []string
	var playerIDs []string
	var selectedIndex int
	
	// Add installed players first
	if len(installedPlayers) > 0 {
		for i, p := range installedPlayers {
			playerOptions = append(playerOptions, fmt.Sprintf("✓ %s", p.Name))
			playerIDs = append(playerIDs, p.ID)
			if p.ID == currentSettings.VideoPlayer {
				selectedIndex = i
			}
		}
	}
	
	// Add separator if we have both installed and not installed
	if len(installedPlayers) > 0 && len(installedPlayers) < len(allPlayers) {
		playerOptions = append(playerOptions, "──────────────────────────")
		playerIDs = append(playerIDs, "")
	}
	
	// Add not installed players
	for _, p := range allPlayers {
		if !p.IsInstalled {
			playerOptions = append(playerOptions, fmt.Sprintf("✗ %s (Not Installed)", p.Name))
			playerIDs = append(playerIDs, p.ID)
		}
	}
	
	playerSelect := widget.NewSelect(playerOptions, nil)
	if len(playerOptions) > 0 && selectedIndex >= 0 && selectedIndex < len(playerOptions) {
		playerSelect.SetSelected(playerOptions[selectedIndex])
	}
	
	// Compact player info text
	var playerInfoText string
	if len(installedPlayers) == 0 {
		playerInfoText = "⚠ No video players detected"
	} else {
		detectedPlayerNames := ""
		for i, p := range installedPlayers {
			if i > 0 {
				detectedPlayerNames += ", "
			}
			detectedPlayerNames += p.Name
		}
		playerInfoText = fmt.Sprintf("✓ Detected: %s", detectedPlayerNames)
	}
	playerInfo := widget.NewLabel(playerInfoText)
	playerInfo.Wrapping = fyne.TextWrapWord
	
	// Modern subtitle language section
	subtitleLabel := CreateHeader("Subtitle Language")
	languageOptions := settings.GetLanguageOptions()
	subtitleSelect := widget.NewSelect(languageOptions, nil)
	subtitleSelect.SetSelected(settings.GetLanguageDisplayString(currentSettings.SubtitleLanguage))
	
	// Modern audio language section
	audioLabel := CreateHeader("Audio Language")
	audioSelect := widget.NewSelect(languageOptions, nil)
	audioSelect.SetSelected(settings.GetLanguageDisplayString(currentSettings.AudioLanguage))
	
	// Modern auto-next toggle
	autoNextLabel := CreateHeader("Auto-play Next Episode")
	autoNextCheck := widget.NewCheck("Enable auto-play", nil)
	autoNextCheck.SetChecked(currentSettings.AutoNext)
	
	// Auto-close player toggle
	autoCloseLabel := CreateHeader("Player Behavior")
	autoCloseCheck := widget.NewCheck("Auto-close player when finished", nil)
	autoCloseCheck.SetChecked(currentSettings.AutoClosePlayer)
	
	// Fullscreen toggle
	fullscreenCheck := widget.NewCheck("Start in fullscreen mode", nil)
	fullscreenCheck.SetChecked(currentSettings.Fullscreen)
	
	// Compact info text
	infoText := widget.NewLabel("Settings apply when playing content")
	infoText.Wrapping = fyne.TextWrapWord
	
	// Modern save button
	saveBtn := CreateIconButtonWithImportance("Save Settings", "", widget.HighImportance, func() {
		// Find selected player ID
		selectedPlayerID := currentSettings.VideoPlayer
		selectedIdx := -1
		for i, opt := range playerOptions {
			if opt == playerSelect.Selected {
				selectedIdx = i
				break
			}
		}
		if selectedIdx >= 0 && selectedIdx < len(playerIDs) && playerIDs[selectedIdx] != "" {
			selectedPlayerID = playerIDs[selectedIdx]
		}
		
		newSettings := &settings.Settings{
			SubtitleLanguage: settings.GetLanguageCode(subtitleSelect.Selected),
			AudioLanguage:    settings.GetLanguageCode(audioSelect.Selected),
			AutoNext:         autoNextCheck.Checked,
			VideoPlayer:      selectedPlayerID,
			AutoClosePlayer:  autoCloseCheck.Checked,
			Fullscreen:       fullscreenCheck.Checked,
		}
		
		settings.Save(newSettings)
		fmt.Println("✓ Settings saved successfully!")
	})
	
	// Compact form layout
	formWithButtons := container.NewVBox(
		videoPlayerLabel,
		playerSelect,
		playerInfo,
		widget.NewSeparator(),
		subtitleLabel,
		subtitleSelect,
		widget.NewSeparator(),
		audioLabel,
		audioSelect,
		widget.NewSeparator(),
		autoNextLabel,
		autoNextCheck,
		widget.NewSeparator(),
		autoCloseLabel,
		autoCloseCheck,
		fullscreenCheck,
		widget.NewSeparator(),
		infoText,
		widget.NewSeparator(),
		saveBtn,
	)
	
	// Create modern custom dialog
	customDialog := dialog.NewCustom("Settings", "Close", formWithButtons, currentWindow)
	customDialog.Resize(fyne.NewSize(420, 540))
	customDialog.Show()
}

