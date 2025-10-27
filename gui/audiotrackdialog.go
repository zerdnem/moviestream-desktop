package gui

import (
	"fmt"
	"moviestream-gui/api"
	"moviestream-gui/audiotracks"
	"moviestream-gui/player"
	"strings"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowAudioTrackDialog shows a dialog to select from available audio tracks or add external ones
func ShowAudioTrackDialog(
	title string,
	tmdbID int,
	season, episode int,
	streamURL string,
	subtitleURLs []string,
	availableAudioTracks []api.AudioTrack,
	onEnd player.OnPlaybackEndCallback,
) {
	// Create audio track manager
	manager := audiotracks.NewManager()
	
	// Store added audio tracks
	var audioTrackURLs []string
	
	// Create modal dialog
	content := container.NewVBox()
	
	// Header
	headerText := widget.NewLabel("Select Audio Track")
	headerText.Wrapping = fyne.TextWrapWord
	headerText.TextStyle = fyne.TextStyle{Bold: true}
	
	var infoText *widget.Label
	if len(availableAudioTracks) > 0 {
		infoText = widget.NewLabel(fmt.Sprintf("Found %d audio track(s) available from the streaming service. Select one to play:", len(availableAudioTracks)))
	} else {
		infoText = widget.NewLabel("No audio tracks available from streaming service. You can add external audio tracks manually.")
	}
	infoText.Wrapping = fyne.TextWrapWord
	
	content.Add(headerText)
	content.Add(infoText)
	content.Add(widget.NewSeparator())
	
	// Dialog buttons
	var d dialog.Dialog
	var selectedAPITrack *api.AudioTrack
	
	// Show available API audio tracks if any
	if len(availableAudioTracks) > 0 {
		apiTracksLabel := widget.NewLabel("Available Audio Tracks:")
		apiTracksLabel.TextStyle = fyne.TextStyle{Bold: true}
		content.Add(apiTracksLabel)
		
		// Create radio group for audio track selection
		var trackOptions []string
		for _, track := range availableAudioTracks {
			option := fmt.Sprintf("%s - %s", track.Name, track.Description)
			trackOptions = append(trackOptions, option)
		}
		
		trackRadio := widget.NewRadioGroup(trackOptions, func(selected string) {
			// Find selected track
			for i, option := range trackOptions {
				if option == selected {
					selectedAPITrack = &availableAudioTracks[i]
					break
				}
			}
		})
		
		// Auto-select first track
		if len(trackOptions) > 0 {
			trackRadio.SetSelected(trackOptions[0])
			selectedAPITrack = &availableAudioTracks[0]
		}
		
		content.Add(trackRadio)
		content.Add(widget.NewLabel(""))
		
		// Play button for API tracks
		playAPITrackBtn := widget.NewButton("Play with Selected Track", func() {
			if selectedAPITrack == nil {
				dialog.ShowError(fmt.Errorf("please select an audio track"), currentWindow)
				return
			}
			
			d.Hide()
			
			// Show loading progress
			progress := dialog.NewProgressInfinite("Loading Stream", 
				fmt.Sprintf("Loading stream with: %s", selectedAPITrack.Name), 
				currentWindow)
			progress.Show()
			
			go func() {
				// Get stream URL with selected audio track
				contentType := "movie"
				if season > 0 && episode > 0 {
					contentType = "tv"
				}
				
				newStreamURL, err := api.GetStreamURLWithAudioTrack(tmdbID, contentType, season, episode, selectedAPITrack.Data)
				
				fyne.Do(func() {
					progress.Hide()
					
					if err != nil {
						dialog.ShowError(fmt.Errorf("failed to load stream with selected audio: %v", err), currentWindow)
						return
					}
					
					// Play with selected audio track
					if err := player.PlayWithMPVAndCallback(newStreamURL, title, subtitleURLs, onEnd); err != nil {
						dialog.ShowError(err, currentWindow)
					} else {
						fmt.Printf("\n▶ Playing: %s\n", title)
						fmt.Printf("   🎵 Audio: %s - %s\n", selectedAPITrack.Name, selectedAPITrack.Description)
						if len(subtitleURLs) > 0 {
							fmt.Printf("   ✓ %d subtitle(s) loaded\n", len(subtitleURLs))
						}
						fmt.Println()
					}
				})
			}()
		})
		playAPITrackBtn.Importance = widget.HighImportance
		content.Add(playAPITrackBtn)
		content.Add(widget.NewSeparator())
	}
	
	// Instructions for manual tracks
	instructionsText := widget.NewLabel("Or add external audio tracks:\n• Direct URLs (e.g., https://example.com/audio.aac)\n• Local files (e.g., C:\\path\\to\\audio.mp3)\n\nSupported formats: AAC, MP3, OGG, Opus")
	instructionsText.Wrapping = fyne.TextWrapWord
	
	content.Add(instructionsText)
	content.Add(widget.NewLabel(""))
	
	// Audio track list
	tracksLabel := widget.NewLabel("Added Audio Tracks:")
	tracksLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	tracksListContainer := container.NewVBox()
	noTracksLabel := widget.NewLabel("No audio tracks added yet")
	noTracksLabel.TextStyle = fyne.TextStyle{Italic: true}
	tracksListContainer.Add(noTracksLabel)
	
	tracksScroll := container.NewVScroll(tracksListContainer)
	tracksScroll.SetMinSize(fyne.NewSize(500, 150))
	
	content.Add(tracksLabel)
	content.Add(tracksScroll)
	content.Add(widget.NewLabel(""))
	
	// Add audio track section
	addTrackLabel := widget.NewLabel("Add Audio Track:")
	addTrackLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	audioURLEntry := widget.NewEntry()
	audioURLEntry.SetPlaceHolder("Enter URL or local file path...")
	
	// Language selection
	languageOptions := []string{
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
	
	languageSelect := widget.NewSelect(languageOptions, nil)
	languageSelect.SetSelected("en - English")
	
	// Function to update tracks list display (declare as var first for recursion)
	var updateTracksList func()
	updateTracksList = func() {
		tracksListContainer.Objects = nil
		
		if len(audioTrackURLs) == 0 {
			noTracksLabel := widget.NewLabel("No audio tracks added yet")
			noTracksLabel.TextStyle = fyne.TextStyle{Italic: true}
			tracksListContainer.Add(noTracksLabel)
		} else {
			for i, trackURL := range audioTrackURLs {
				trackIndex := i
				trackText := fmt.Sprintf("%d. %s", i+1, trackURL)
				if len(trackURL) > 60 {
					trackText = fmt.Sprintf("%d. %s...", i+1, trackURL[:60])
				}
				
				trackLabel := widget.NewLabel(trackText)
				
				removeBtn := widget.NewButton("Remove", func() {
					// Remove track at trackIndex
					audioTrackURLs = append(audioTrackURLs[:trackIndex], audioTrackURLs[trackIndex+1:]...)
					updateTracksList()
				})
				removeBtn.Importance = widget.DangerImportance
				
				trackRow := container.NewBorder(nil, nil, nil, removeBtn, trackLabel)
				tracksListContainer.Add(trackRow)
			}
		}
		
		tracksListContainer.Refresh()
	}
	
	addTrackBtn := widget.NewButton("Add Track", func() {
		url := strings.TrimSpace(audioURLEntry.Text)
		if url == "" {
			dialog.ShowError(fmt.Errorf("please enter a URL or file path"), currentWindow)
			return
		}
		
		// Validate URL/file
		if err := manager.ValidateAudioFile(url); err != nil {
			// If validation fails, it might still be a valid URL
			// Check if it looks like a URL
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				dialog.ShowError(fmt.Errorf("invalid file path or URL: %v", err), currentWindow)
				return
			}
		}
		
		// Add to list
		audioTrackURLs = append(audioTrackURLs, url)
		audioURLEntry.SetText("")
		updateTracksList()
	})
	addTrackBtn.Importance = widget.HighImportance
	
	addTrackSection := container.NewVBox(
		addTrackLabel,
		audioURLEntry,
		container.NewBorder(nil, nil, widget.NewLabel("Language:"), nil, languageSelect),
		addTrackBtn,
	)
	
	content.Add(addTrackSection)
	content.Add(widget.NewSeparator())
	
	playBtn := widget.NewButton("Play with Audio Tracks", func() {
		d.Hide()
		
		// Play with audio tracks
		if err := player.PlayWithMPVAndAudio(streamURL, title, subtitleURLs, audioTrackURLs, onEnd); err != nil {
			dialog.ShowError(err, currentWindow)
		} else {
			fmt.Printf("\n▶ Playing: %s\n", title)
			if len(audioTrackURLs) > 0 {
				fmt.Printf("   ✓ %d external audio track(s) loaded\n", len(audioTrackURLs))
			}
			if len(subtitleURLs) > 0 {
				fmt.Printf("   ✓ %d subtitle(s) loaded\n", len(subtitleURLs))
			}
			fmt.Println()
		}
	})
	
	playWithoutAudioBtn := widget.NewButton("Play Without Extra Audio", func() {
		d.Hide()
		// Play without extra audio tracks
		if err := player.PlayWithMPVAndCallback(streamURL, title, subtitleURLs, onEnd); err != nil {
			dialog.ShowError(err, currentWindow)
		} else {
			fmt.Printf("\n▶ Playing: %s\n", title)
			fmt.Printf("   ⚠ Using original audio only\n\n")
		}
	})
	
	cancelBtn := widget.NewButton("Cancel", func() {
		d.Hide()
	})
	
	playBtn.Importance = widget.HighImportance
	
	buttons := container.NewHBox(
		cancelBtn,
		playWithoutAudioBtn,
		playBtn,
	)
	
	content.Add(buttons)
	
	// Create dialog
	d = dialog.NewCustom("Audio Tracks", "", content, currentWindow)
	d.Resize(fyne.NewSize(600, 600))
	d.Show()
}

// ShowAudioTrackSelectionButton shows a button to open audio track dialog
func ShowAudioTrackSelectionButton(
	title string,
	tmdbID int,
	season, episode int,
	streamURL string,
	subtitleURLs []string,
	availableAudioTracks []api.AudioTrack,
	onEnd player.OnPlaybackEndCallback,
) *widget.Button {
	btn := CreateIconButton("Add Audio Tracks", IconAdd, func() {
		ShowAudioTrackDialog(title, tmdbID, season, episode, streamURL, subtitleURLs, availableAudioTracks, onEnd)
	})
	btn.Importance = widget.LowImportance
	return btn
}

