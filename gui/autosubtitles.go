package gui

import (
	"fmt"
	"moviestream-gui/player"
	"moviestream-gui/subtitles"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// AutoDownloadAndPlaySubtitles automatically searches, downloads, and plays with the best subtitle
func AutoDownloadAndPlaySubtitles(
	title string,
	tmdbID int,
	season, episode int,
	streamURL string,
	onEnd player.OnPlaybackEndCallback,
) {
	// Show loading indicator with dynamic message
	progress := dialog.NewProgressInfinite("Loading Subtitles", 
		"Searching for subtitles from OpenSubtitles and Addic7ed...", 
		currentWindow)
	progress.Show()
	
	go func() {
		manager := subtitles.NewManager()
		
		// Search for subtitles
		fmt.Println("🔍 Searching for subtitles...")
		results, err := manager.SearchSubtitles(title, tmdbID, season, episode)
		
		if err != nil || len(results) == 0 {
			// No subtitles found - play without subtitles
			fyne.Do(func() {
				progress.Hide()
				fmt.Printf("⚠ No subtitles found. Playing without subtitles.\n")
				
				if err := player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd); err != nil {
					dialog.ShowError(err, currentWindow)
				} else {
					fmt.Printf("\n▶ Playing: %s\n", title)
					fmt.Printf("   ⚠ No subtitles loaded\n\n")
				}
			})
			return
		}
		
		// Separate results by source and pick best
		var bestSubtitle *subtitles.SubtitleResult
		var sourceName string
		
		// Priority: OpenSubtitles first, then Addic7ed
		for i, result := range results {
			if len(result.ID) > 0 && result.ID[0] != '/' {
				// OpenSubtitles (numeric ID)
				bestSubtitle = &results[i]
				sourceName = "OpenSubtitles"
				break
			}
		}
		
		// If no OpenSubtitles, try Addic7ed
		if bestSubtitle == nil {
			for i, result := range results {
				if len(result.ID) > 0 && result.ID[0] == '/' {
					// Addic7ed (ID starts with "/")
					bestSubtitle = &results[i]
					sourceName = "Addic7ed"
					break
				}
			}
		}
		
		if bestSubtitle == nil {
			// Shouldn't happen, but handle it
			fyne.Do(func() {
				progress.Hide()
				fmt.Printf("⚠ No valid subtitles found. Playing without subtitles.\n")
				
				if err := player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd); err != nil {
					dialog.ShowError(err, currentWindow)
				} else {
					fmt.Printf("\n▶ Playing: %s\n", title)
					fmt.Printf("   ⚠ No subtitles loaded\n\n")
				}
			})
			return
		}
		
		// Found best subtitle - download it
		fmt.Printf("✓ Auto-selected: %s from %s\n", bestSubtitle.FileName, sourceName)
		
		// Update progress message for downloading
		fyne.Do(func() {
			progress.Hide()
		})
		
		downloadProgress := dialog.NewProgressInfinite("Downloading Subtitle", 
			fmt.Sprintf("Downloading from %s...", sourceName), 
			currentWindow)
		fyne.Do(func() {
			downloadProgress.Show()
		})
		
		fmt.Printf("⬇ Downloading subtitle from %s...\n", sourceName)
		subtitlePath, downloadErr := manager.DownloadSubtitle(*bestSubtitle)
		
		fyne.Do(func() {
			downloadProgress.Hide()
			
			if downloadErr != nil {
				// Download failed - try alternative source if available
				fmt.Printf("⚠ Download failed: %v\n", downloadErr)
				
				// Check if we can try alternative source
				if sourceName == "OpenSubtitles" && season > 0 && episode > 0 {
					// Try Addic7ed as fallback
					fmt.Println("🔄 Trying Addic7ed as alternative...")
					
					addic7edSearchProgress := dialog.NewProgressInfinite("Searching Addic7ed", 
						"Searching for subtitles from Addic7ed...", 
						currentWindow)
					addic7edSearchProgress.Show()
					
					go func() {
						addic7edClient := subtitles.NewAddic7edClient()
						addic7edResults, addic7edErr := addic7edClient.SearchByTitle(title, bestSubtitle.Language, season, episode)
						
						fyne.Do(func() {
							addic7edSearchProgress.Hide()
							
							if addic7edErr != nil || len(addic7edResults) == 0 {
								// Addic7ed also failed - play without subtitles
								fmt.Println("⚠ Addic7ed also failed. Playing without subtitles.")
								
								if err := player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd); err != nil {
									dialog.ShowError(err, currentWindow)
								} else {
									fmt.Printf("\n▶ Playing: %s\n", title)
									fmt.Printf("   ⚠ No subtitles loaded\n\n")
								}
								return
							}
							
							// Try downloading from Addic7ed
							fmt.Println("⬇ Downloading subtitle from Addic7ed...")
							
							addic7edDownloadProgress := dialog.NewProgressInfinite("Downloading Subtitle", 
								"Downloading from Addic7ed...", 
								currentWindow)
							addic7edDownloadProgress.Show()
							
							go func() {
								addic7edPath, addic7edDownloadErr := manager.DownloadSubtitle(addic7edResults[0])
								
								fyne.Do(func() {
									addic7edDownloadProgress.Hide()
									
									if addic7edDownloadErr != nil {
										// Both sources failed - play without subtitles
										fmt.Printf("⚠ Both sources failed. Playing without subtitles.\n")
										
										if err := player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd); err != nil {
											dialog.ShowError(err, currentWindow)
										} else {
											fmt.Printf("\n▶ Playing: %s\n", title)
											fmt.Printf("   ⚠ No subtitles loaded\n\n")
										}
									} else {
										// Success with Addic7ed!
										if err := player.PlayWithMPVAndCallback(streamURL, title, []string{addic7edPath}, onEnd); err != nil {
											dialog.ShowError(err, currentWindow)
										} else {
											fmt.Printf("\n▶ Playing: %s\n", title)
											fmt.Printf("   ✓ Subtitle auto-loaded from Addic7ed: %s (%s)\n\n", addic7edResults[0].FileName, addic7edResults[0].LanguageName)
										}
									}
								})
							}()
						})
					}()
				} else {
					// Can't try alternative - play without subtitles
					if err := player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd); err != nil {
						dialog.ShowError(err, currentWindow)
					} else {
						fmt.Printf("\n▶ Playing: %s\n", title)
						fmt.Printf("   ⚠ No subtitles loaded\n\n")
					}
				}
				return
			}
			
			// Success! Play with subtitle
			if err := player.PlayWithMPVAndCallback(streamURL, title, []string{subtitlePath}, onEnd); err != nil {
				dialog.ShowError(err, currentWindow)
			} else {
				fmt.Printf("\n▶ Playing: %s\n", title)
				fmt.Printf("   ✓ Subtitle auto-loaded from %s: %s (%s)\n\n", sourceName, bestSubtitle.FileName, bestSubtitle.LanguageName)
			}
		})
	}()
}

