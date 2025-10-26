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
	
	// Video Player selector
	videoPlayerLabel := widget.NewLabelWithStyle("Video Player:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	
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
	
	// Player info text
	var playerInfoText string
	if len(installedPlayers) == 0 {
		playerInfoText = "⚠ No video players detected! Please install MPV, VLC, or another supported player."
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
	
	// Subtitle language selector
	subtitleLabel := widget.NewLabel("Subtitle Language:")
	languageOptions := settings.GetLanguageOptions()
	subtitleSelect := widget.NewSelect(languageOptions, nil)
	subtitleSelect.SetSelected(settings.GetLanguageDisplayString(currentSettings.SubtitleLanguage))
	
	// Audio language selector
	audioLabel := widget.NewLabel("Audio Language:")
	audioSelect := widget.NewSelect(languageOptions, nil)
	audioSelect.SetSelected(settings.GetLanguageDisplayString(currentSettings.AudioLanguage))
	
	// Auto-next toggle
	autoNextLabel := widget.NewLabel("Auto-play Next Episode:")
	autoNextCheck := widget.NewCheck("", nil)
	autoNextCheck.SetChecked(currentSettings.AutoNext)
	
	// Info text
	infoText := widget.NewLabel("These settings will be applied when playing content.")
	infoText.Wrapping = fyne.TextWrapWord
	
	// Save button
	saveBtn := widget.NewButton("Save", func() {
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
		}
		
		settings.Save(newSettings)
		dialog.ShowInformation("Success", "Settings saved successfully!", currentWindow)
	})
	
	// Update form with buttons
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
		container.NewHBox(autoNextLabel, autoNextCheck),
		widget.NewSeparator(),
		infoText,
		widget.NewSeparator(),
		saveBtn,
	)
	
	// Create custom dialog
	customDialog := dialog.NewCustom("Settings", "Close", formWithButtons, currentWindow)
	customDialog.Resize(fyne.NewSize(450, 500))
	customDialog.Show()
}

