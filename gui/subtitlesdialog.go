package gui

import (
	"fmt"
	"moviestream-gui/player"
	"moviestream-gui/subtitles"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowSubtitleDownloadDialog shows a dialog to download subtitles from external sources
func ShowSubtitleDownloadDialog(
	title string,
	tmdbID int,
	season, episode int,
	streamURL string,
	onEnd player.OnPlaybackEndCallback,
) {
	// Create subtitle manager
	manager := subtitles.NewManager()
	
	// Create modal dialog
	content := container.NewVBox()
	
	// Header
	headerText := widget.NewLabel("No subtitles found in stream")
	headerText.Wrapping = fyne.TextWrapWord
	headerText.TextStyle = fyne.TextStyle{Bold: true}
	
	infoText := widget.NewLabel("Searching for subtitles from OpenSubtitles and Addic7ed...")
	infoText.Wrapping = fyne.TextWrapWord
	
	content.Add(headerText)
	content.Add(infoText)
	content.Add(widget.NewSeparator())
	
	// Results container
	resultsContainer := container.NewVBox()
	resultsScroll := container.NewVScroll(resultsContainer)
	resultsScroll.SetMinSize(fyne.NewSize(500, 200))
	
	// Loading indicator
	loadingLabel := widget.NewLabel("Searching for subtitles...")
	progressBar := widget.NewProgressBarInfinite()
	loadingContainer := container.NewVBox(loadingLabel, progressBar)
	
	content.Add(loadingContainer)
	
	var subtitleResults []subtitles.SubtitleResult
	var selectedSubtitle *subtitles.SubtitleResult
	
	// Dialog buttons
	var d dialog.Dialog
	
	playWithoutBtn := widget.NewButton("Play Without Subtitles", func() {
		d.Hide()
		// Play with empty subtitle list
		if err := player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd); err != nil {
			dialog.ShowError(err, currentWindow)
		} else {
			fmt.Printf("\n▶ Playing: %s\n", title)
			fmt.Printf("   ⚠ No subtitles loaded\n\n")
		}
	})
	
	downloadAndPlayBtn := widget.NewButton("Download & Play", func() {
		if selectedSubtitle == nil {
			dialog.ShowError(fmt.Errorf("please select a subtitle"), currentWindow)
			return
		}
		
		// Store selected subtitle for retry
		selectedSub := selectedSubtitle
		
		d.Hide()
		
		// Show downloading progress
		progress := dialog.NewProgressInfinite("Downloading Subtitle", 
			fmt.Sprintf("Downloading: %s", selectedSub.FileName), 
			currentWindow)
		progress.Show()
		
		go func() {
			// Download subtitle
			subtitlePath, err := manager.DownloadSubtitle(*selectedSub)
			
			fyne.Do(func() {
				progress.Hide()
				
				if err != nil {
					// Check if this was OpenSubtitles and if Addic7ed is available as alternative
					isOpenSubtitles := len(selectedSub.ID) > 0 && selectedSub.ID[0] != '/'
					isTVShow := season > 0 && episode > 0
					
					if isOpenSubtitles && isTVShow {
						// Offer to try Addic7ed as alternative
						errorMsg := fmt.Sprintf("OpenSubtitles download failed:\n%v\n\nWould you like to try Addic7ed instead?\n(Addic7ed is a reliable alternative for TV shows)", err)
						
						retryDialog := dialog.NewConfirm("Download Failed", errorMsg, 
							func(tryAddic7ed bool) {
								if tryAddic7ed {
									// Search Addic7ed for alternatives
									fmt.Println("\nTrying Addic7ed as alternative...")
									progress.Show()
									progress.Refresh()
									
									go func() {
										// Search specifically in Addic7ed
										addic7edClient := subtitles.NewAddic7edClient()
										addic7edResults, addic7edErr := addic7edClient.SearchByTitle(title, selectedSub.Language, season, episode)
										
										fyne.Do(func() {
											progress.Hide()
											
											if addic7edErr != nil || len(addic7edResults) == 0 {
												// Addic7ed also failed - don't auto-play
												errorText := "Addic7ed search also failed. Please try again or play without subtitles."
												if addic7edErr != nil {
													errorText = fmt.Sprintf("Addic7ed error: %v\n\nPlease try again or play without subtitles.", addic7edErr)
												}
												dialog.ShowError(fmt.Errorf("%s", errorText), currentWindow)
												return
											}
											
											// Try downloading from Addic7ed
											fmt.Printf("✓ Found %d subtitles from Addic7ed\n", len(addic7edResults))
											progress.Show()
											
											go func() {
												// Download first Addic7ed result
												addic7edPath, downloadErr := manager.DownloadSubtitle(addic7edResults[0])
												
												fyne.Do(func() {
													progress.Hide()
													
													if downloadErr != nil {
														dialog.ShowError(fmt.Errorf("Addic7ed download also failed: %v\n\nPlease try again or play without subtitles.", downloadErr), currentWindow)
													} else {
														// Success with Addic7ed!
														if err := player.PlayWithMPVAndCallback(streamURL, title, []string{addic7edPath}, onEnd); err != nil {
															dialog.ShowError(err, currentWindow)
														} else {
															fmt.Printf("\n▶ Playing: %s\n", title)
															fmt.Printf("   ✓ Subtitle loaded from Addic7ed: %s (%s)\n\n", addic7edResults[0].FileName, addic7edResults[0].LanguageName)
														}
													}
												})
											}()
										})
									}()
								} else {
									// Play without subtitles
									if err := player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd); err != nil {
										dialog.ShowError(err, currentWindow)
									} else {
										fmt.Printf("\n▶ Playing: %s\n", title)
										fmt.Printf("   ⚠ No subtitles loaded\n\n")
									}
								}
							}, currentWindow)
						retryDialog.SetDismissText("Play Without Subtitles")
						retryDialog.SetConfirmText("Try Addic7ed")
						retryDialog.Show()
					} else {
						// Either not OpenSubtitles, or not a TV show, or Addic7ed download failed
						// Offer simple retry
						errorMsg := fmt.Sprintf("Failed to download subtitle:\n%v\n\nOptions:\n• Click OK to retry\n• Or play without subtitles", err)
						
						retryDialog := dialog.NewConfirm("Download Failed", errorMsg, 
							func(retry bool) {
								if retry {
									// Retry download
									progress.Show()
									go func() {
										subtitlePath, err := manager.DownloadSubtitle(*selectedSub)
										fyne.Do(func() {
											progress.Hide()
											if err != nil {
												// Second failure - don't auto-play
												dialog.ShowError(fmt.Errorf("Download failed again: %v\n\nPlease try again or play without subtitles.", err), currentWindow)
											} else {
												// Success on retry!
												if err := player.PlayWithMPVAndCallback(streamURL, title, []string{subtitlePath}, onEnd); err != nil {
													dialog.ShowError(err, currentWindow)
												} else {
													fmt.Printf("\n▶ Playing: %s\n", title)
													fmt.Printf("   ✓ Subtitle loaded: %s (%s) [retry successful]\n\n", selectedSub.FileName, selectedSub.LanguageName)
												}
											}
										})
									}()
								} else {
									// Play without subtitles
									if err := player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd); err != nil {
										dialog.ShowError(err, currentWindow)
									} else {
										fmt.Printf("\n▶ Playing: %s\n", title)
										fmt.Printf("   ⚠ No subtitles loaded\n\n")
									}
								}
							}, currentWindow)
						retryDialog.SetDismissText("Play Without Subtitles")
						retryDialog.SetConfirmText("Retry Download")
						retryDialog.Show()
					}
					return
				}
				
				// Play with subtitle
				if err := player.PlayWithMPVAndCallback(streamURL, title, []string{subtitlePath}, onEnd); err != nil {
					dialog.ShowError(err, currentWindow)
				} else {
					fmt.Printf("\n▶ Playing: %s\n", title)
					fmt.Printf("   ✓ Subtitle loaded: %s (%s)\n\n", selectedSub.FileName, selectedSub.LanguageName)
				}
			})
		}()
	})
	
	downloadAndPlayBtn.Importance = widget.HighImportance
	downloadAndPlayBtn.Disable()
	
	cancelBtn := widget.NewButton("Cancel", func() {
		d.Hide()
	})
	
	buttons := container.NewHBox(
		cancelBtn,
		playWithoutBtn,
		downloadAndPlayBtn,
	)
	
	content.Add(resultsScroll)
	content.Add(widget.NewSeparator())
	content.Add(buttons)
	
	// Create dialog
	d = dialog.NewCustom("Download Subtitles", "", content, currentWindow)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
	
	// Search for subtitles
	go func() {
		results, err := manager.SearchSubtitles(title, tmdbID, season, episode)
		
		fyne.Do(func() {
			// Remove loading indicator
			content.Remove(loadingContainer)
			
			if err != nil {
				errorLabel := widget.NewLabel(fmt.Sprintf("Search failed: %v", err))
				errorLabel.Wrapping = fyne.TextWrapWord
				
				manualLabel := widget.NewLabel("\nYou can manually download subtitles from:\n• OpenSubtitles.org\n• Subscene.com")
				manualLabel.Wrapping = fyne.TextWrapWord
				
				resultsContainer.Add(errorLabel)
				resultsContainer.Add(manualLabel)
				return
			}
			
			if len(results) == 0 {
				noResultsLabel := widget.NewLabel("No subtitles found for this content.\n\nYou can try:\n• OpenSubtitles.org\n• Subscene.com\n• YIFY Subtitles")
				noResultsLabel.Wrapping = fyne.TextWrapWord
				resultsContainer.Add(noResultsLabel)
				return
			}
			
			subtitleResults = results
			
			// Separate results by source
			var openSubtitlesResults []subtitles.SubtitleResult
			var addic7edResults []subtitles.SubtitleResult
			
			for _, result := range results {
				if len(result.ID) > 0 && result.ID[0] == '/' {
					// Addic7ed (ID starts with "/")
					addic7edResults = append(addic7edResults, result)
				} else {
					// OpenSubtitles (numeric ID)
					openSubtitlesResults = append(openSubtitlesResults, result)
				}
			}
			
			// Show results summary
			resultsLabel := widget.NewLabel(fmt.Sprintf("Found %d subtitle(s):", len(results)))
			resultsLabel.TextStyle = fyne.TextStyle{Bold: true}
			resultsContainer.Add(resultsLabel)
			
			// Create radio group for subtitle selection
			var options []string
			var optionToResult []int // Map option index to result index
			
			// Add OpenSubtitles results
			if len(openSubtitlesResults) > 0 {
				osHeader := widget.NewLabel("OpenSubtitles:")
				osHeader.TextStyle = fyne.TextStyle{Bold: true}
				resultsContainer.Add(osHeader)
				
				for i, result := range openSubtitlesResults {
					option := fmt.Sprintf("  %s (%s) - %s", 
						result.LanguageName, 
						result.Language, 
						result.FileName)
					options = append(options, option)
					
					// Find index in original results
					for j, r := range results {
						if r.ID == result.ID && r.FileName == result.FileName {
							optionToResult = append(optionToResult, j)
							break
						}
					}
					
					// Auto-select first OpenSubtitles result
					if i == 0 && selectedSubtitle == nil {
						selectedSubtitle = &results[optionToResult[len(optionToResult)-1]]
						downloadAndPlayBtn.Enable()
					}
				}
			}
			
			// Add Addic7ed results
			if len(addic7edResults) > 0 {
				addic7edHeader := widget.NewLabel("Addic7ed:")
				addic7edHeader.TextStyle = fyne.TextStyle{Bold: true}
				resultsContainer.Add(addic7edHeader)
				
				for _, result := range addic7edResults {
					option := fmt.Sprintf("  %s (%s) - %s", 
						result.LanguageName, 
						result.Language, 
						result.FileName)
					options = append(options, option)
					
					// Find index in original results
					for j, r := range results {
						if r.ID == result.ID && r.FileName == result.FileName {
							optionToResult = append(optionToResult, j)
							break
						}
					}
					
					// Auto-select first result if no OpenSubtitles results
					if len(openSubtitlesResults) == 0 && len(optionToResult) == 1 {
						selectedSubtitle = &results[optionToResult[0]]
						downloadAndPlayBtn.Enable()
					}
				}
			}
			
			radioGroup := widget.NewRadioGroup(options, func(selected string) {
				// Find selected subtitle using optionToResult mapping
				for i, option := range options {
					if option == selected {
						selectedSubtitle = &subtitleResults[optionToResult[i]]
						downloadAndPlayBtn.Enable()
						break
					}
				}
			})
			
			if len(options) > 0 {
				radioGroup.SetSelected(options[0])
			}
			
			resultsContainer.Add(radioGroup)
			
			// Refresh layout
			d.Refresh()
		})
	}()
}
