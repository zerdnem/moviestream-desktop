package gui

import (
	"fmt"
	"image/color"
	"moviestream-gui/api"
	"moviestream-gui/downloader"
	"moviestream-gui/history"
	"moviestream-gui/player"
	"moviestream-gui/queue"
	"moviestream-gui/settings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Episode tracking for auto-next functionality
var (
	currentTVShowID      int
	currentShowName      string
	currentSeason        int
	currentEpisode       int
	totalEpisodes        int
	seasonDetailsCache   map[string]*api.Season // Cache season details for auto-next
	autoNextCancelled    bool                   // Flag to cancel auto-next for current session
	autoNextInProgress   bool                   // Flag to prevent multiple auto-next triggers
)

func init() {
	seasonDetailsCache = make(map[string]*api.Season)
	autoNextCancelled = false
	autoNextInProgress = false
}

// showTVDetails shows detailed view for a TV show with episode list
func showTVDetails(tvShow api.TVShow) {
	// Reset auto-next cancellation when viewing a new show
	autoNextCancelled = false
	autoNextInProgress = false
	
	progress := dialog.NewProgressInfinite("Loading", "Fetching TV show details...", currentWindow)
	progress.Show()

	go func() {
		details, err := api.GetTVDetails(tvShow.ID)
		
		// Always hide progress
		fyne.Do(func() {
			progress.Hide()
		})

		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to get TV details: %v", err), currentWindow)
			return
		}

		// Update UI on main thread
		fyne.Do(func() {
			showTVDetailsUI(details)
		})
	}()
}

// showTVDetailsUI displays the TV show details with seasons and episodes
func showTVDetailsUI(tvDetails *api.TVDetails) {
	// Title
	titleText := canvas.NewText(tvDetails.Name, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	titleText.TextSize = 20
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.Alignment = fyne.TextAlignCenter

	// Info
	infoText := fmt.Sprintf("First Air Date: %s\nRating: %.1f/10\nSeasons: %d",
		tvDetails.FirstAirDate, tvDetails.VoteAverage, len(tvDetails.Seasons))
	infoLabel := widget.NewLabel(infoText)

	// Overview
	overviewLabel := widget.NewLabel(tvDetails.Overview)
	overviewLabel.Wrapping = fyne.TextWrapWord

	// Back button
	backBtn := widget.NewButton("← Back to Search", func() {
		GoBackToSearch()
	})

	// Season selector
	var seasonOptions []string
	validSeasons := []api.Season{}
	
	for _, season := range tvDetails.Seasons {
		if season.SeasonNumber >= 0 { // Skip specials (season 0) unless you want them
			seasonOptions = append(seasonOptions, fmt.Sprintf("Season %d", season.SeasonNumber))
			validSeasons = append(validSeasons, season)
		}
	}

	if len(seasonOptions) == 0 {
		content := container.NewVBox(
			backBtn,
			widget.NewSeparator(),
			container.NewCenter(titleText),
			infoLabel,
			widget.NewSeparator(),
			overviewLabel,
			widget.NewSeparator(),
			widget.NewLabel("No seasons available"),
		)
		currentWindow.SetContent(container.NewVScroll(content))
		return
	}

	episodesList := container.NewVBox()
	
	seasonSelect := widget.NewSelect(seasonOptions, func(selected string) {
		// Find the selected season
		for i, opt := range seasonOptions {
			if opt == selected {
				loadEpisodes(tvDetails.ID, validSeasons[i].SeasonNumber, episodesList, tvDetails.Name)
				break
			}
		}
	})

	if len(seasonOptions) > 0 {
		seasonSelect.SetSelected(seasonOptions[0])
	}

	// Episodes scroll container
	episodesScroll := container.NewVScroll(episodesList)
	episodesScroll.SetMinSize(fyne.NewSize(400, 300))

	content := container.NewVBox(
		backBtn,
		widget.NewSeparator(),
		container.NewCenter(titleText),
		infoLabel,
		widget.NewSeparator(),
		overviewLabel,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Select Season:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		seasonSelect,
		widget.NewSeparator(),
		episodesScroll,
	)

	currentWindow.SetContent(container.NewVScroll(content))
}

// loadEpisodes loads and displays episodes for a specific season
func loadEpisodes(tvID, seasonNum int, episodesList *fyne.Container, showName string) {
	fyne.Do(func() {
		episodesList.Objects = []fyne.CanvasObject{
			widget.NewLabel("Loading episodes..."),
		}
		episodesList.Refresh()
	})

	go func() {
		season, err := api.GetSeasonDetails(tvID, seasonNum)
		if err != nil {
			fyne.Do(func() {
				episodesList.Objects = []fyne.CanvasObject{
					widget.NewLabel(fmt.Sprintf("Failed to load episodes: %v", err)),
				}
				episodesList.Refresh()
			})
			return
		}

		// Cache season details for auto-next functionality
		cacheKey := fmt.Sprintf("%d-%d", tvID, seasonNum)
		seasonDetailsCache[cacheKey] = season

		// Build episode widgets
		var episodeWidgets []fyne.CanvasObject

		if len(season.Episodes) == 0 {
			episodeWidgets = append(episodeWidgets, widget.NewLabel("No episodes found"))
		} else {
			for _, episode := range season.Episodes {
				ep := episode // Capture for closure
				
				// Episode card
				episodeTitle := fmt.Sprintf("E%d: %s", ep.EpisodeNumber, ep.Name)
				titleLabel := widget.NewLabelWithStyle(episodeTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
				
				overviewLabel := widget.NewLabel(ep.Overview)
				overviewLabel.Wrapping = fyne.TextWrapWord

				// Watch button
				watchBtn := widget.NewButton("▶ Watch", func() {
					watchEpisodeWithAutoNext(tvID, seasonNum, ep.EpisodeNumber, showName, ep.Name, len(season.Episodes))
				})

				// Download button
				downloadBtn := widget.NewButton("⬇ Download", func() {
					downloadEpisode(tvID, seasonNum, ep.EpisodeNumber, showName, ep.Name)
				})

				// Add to Queue button
				addToQueueBtn := widget.NewButton("➕ Add to Queue", func() {
					addEpisodeToQueue(tvID, showName, seasonNum, ep.EpisodeNumber, ep.Name)
				})

				buttonContainer := container.NewVBox(
					container.NewGridWithColumns(2, watchBtn, downloadBtn),
					addToQueueBtn,
				)

				episodeCard := container.NewVBox(
					titleLabel,
					overviewLabel,
					buttonContainer,
					widget.NewSeparator(),
				)

				episodeWidgets = append(episodeWidgets, episodeCard)
			}
		}

		// Update UI on main thread
		fyne.Do(func() {
			episodesList.Objects = episodeWidgets
			episodesList.Refresh()
		})
	}()
}

// addEpisodeToQueue adds an episode to the playback queue
func addEpisodeToQueue(tvID int, showName string, season, episode int, episodeName string) {
	q := queue.Get()
	q.AddEpisode(tvID, showName, season, episode, episodeName)
	dialog.ShowInformation("Added to Queue", 
		fmt.Sprintf("'%s - S%dE%d: %s' has been added to the playback queue.", showName, season, episode, episodeName), 
		currentWindow)
}

// watchEpisodeWithAutoNext starts playing an episode with auto-next support
func watchEpisodeWithAutoNext(tvID, season, episode int, showName, episodeName string, totalEps int) {
	// If totalEps is 0, we need to fetch it
	if totalEps == 0 {
		go func() {
			// Try to get from cache first
			cacheKey := fmt.Sprintf("%d-%d", tvID, season)
			if seasonDetails, exists := seasonDetailsCache[cacheKey]; exists && len(seasonDetails.Episodes) > 0 {
				// Use cached data
				watchEpisodeWithAutoNextInternal(tvID, season, episode, showName, episodeName, len(seasonDetails.Episodes))
			} else {
				// Fetch season details
				seasonDetails, err := api.GetSeasonDetails(tvID, season)
				if err != nil || len(seasonDetails.Episodes) == 0 {
					// If we can't get episode count, just play without proper auto-next tracking
					fmt.Printf("⚠ Warning: Could not fetch season details for auto-next\n")
					watchEpisodeWithAutoNextInternal(tvID, season, episode, showName, episodeName, 0)
				} else {
					// Cache it for future use
					seasonDetailsCache[cacheKey] = seasonDetails
					watchEpisodeWithAutoNextInternal(tvID, season, episode, showName, episodeName, len(seasonDetails.Episodes))
				}
			}
		}()
	} else {
		watchEpisodeWithAutoNextInternal(tvID, season, episode, showName, episodeName, totalEps)
	}
}

// watchEpisodeWithAutoNextInternal is the internal implementation with known episode count
func watchEpisodeWithAutoNextInternal(tvID, season, episode int, showName, episodeName string, totalEps int) {
	// Update current episode tracking
	currentTVShowID = tvID
	currentShowName = showName
	currentSeason = season
	currentEpisode = episode
	totalEpisodes = totalEps
	autoNextInProgress = false // Reset the flag for new episode

	fmt.Printf("📺 Auto-next tracking: S%dE%d (Total episodes in season: %d)\n", season, episode, totalEps)

	watchEpisode(tvID, season, episode, showName, episodeName)
}

// watchEpisode starts playing an episode in MPV
func watchEpisode(tvID, season, episode int, showName, episodeName string) {
	watchEpisodeInternal(tvID, season, episode, showName, episodeName, true)
}

// watchEpisodeWithoutDialog starts playing an episode without showing the playback started dialog
func watchEpisodeWithoutDialog(tvID, season, episode int, showName, episodeName string) {
	watchEpisodeInternal(tvID, season, episode, showName, episodeName, false)
}

// watchEpisodeInternal is the internal implementation for playing episodes
func watchEpisodeInternal(tvID, season, episode int, showName, episodeName string, showDialog bool) {
	progress := dialog.NewProgressInfinite("Loading Stream", 
		"Fetching stream URL...\nThis may take 10-15 seconds for browser automation", 
		currentWindow)
	progress.Show()

	go func() {
		streamInfo, err := api.GetStreamURL(tvID, "tv", season, episode)
		
		// Always hide progress
		fyne.Do(func() {
			progress.Hide()
		})

		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to get stream: %v", err), currentWindow)
			return
		}

		// Extract subtitle URLs
		var subtitleURLs []string
		for _, sub := range streamInfo.SubtitleURLs {
			subtitleURLs = append(subtitleURLs, sub.URL)
		}

		title := fmt.Sprintf("%s - S%dE%d - %s", showName, season, episode, episodeName)
		
		// Create callback for auto-next and queue
		var onEndCallback player.OnPlaybackEndCallback
		userSettings := settings.Get()
		q := queue.Get()
		
		// First try auto-next, then queue
		if userSettings.AutoNext && tvID == currentTVShowID && season == currentSeason && episode == currentEpisode {
			onEndCallback = func() {
				playNextEpisode()
			}
		} else if !q.IsEmpty() {
			// If auto-next is not enabled but queue has items, play from queue
			onEndCallback = func() {
				playNextInQueue()
			}
		}

		if err := player.PlayWithMPVAndCallback(streamInfo.StreamURL, title, subtitleURLs, onEndCallback); err != nil {
			dialog.ShowError(err, currentWindow)
			return
		}

		// Record in watch history
		h := history.Get()
		h.AddEpisode(tvID, showName, season, episode, episodeName)

		// Only show playback dialog if requested (not for auto-next triggered episodes)
		if showDialog {
			// Show detailed playback status
			var statusMsg string
			episodeInfo := fmt.Sprintf("✓ Playing: %s - S%dE%d\n", showName, season, episode)
			
			if len(subtitleURLs) > 0 {
				statusMsg = fmt.Sprintf("%s\n✓ %d subtitle track(s) loaded", episodeInfo, len(subtitleURLs))
			} else {
				statusMsg = fmt.Sprintf("%s\n⚠ No subtitles found for this episode\n\nNote: You can still add external subtitles in your video player", episodeInfo)
			}
			
			if userSettings.AutoNext && !autoNextCancelled {
				statusMsg += "\n\n▶ Auto-next is enabled\nNext episode will play automatically when this one ends"
			} else if autoNextCancelled {
				statusMsg += "\n\n⏸ Auto-next is disabled for this session"
			}
			
			// Add queue info
			if !q.IsEmpty() {
				statusMsg += fmt.Sprintf("\n\n📋 %d item(s) in queue - will play after current content", q.Size())
			}
			
			dialog.ShowInformation("Playback Started", statusMsg, currentWindow)
		}
	}()
}

// playNextEpisode automatically plays the next episode
func playNextEpisode() {
	userSettings := settings.Get()
	if !userSettings.AutoNext || autoNextCancelled || autoNextInProgress {
		return
	}

	autoNextInProgress = true // Prevent multiple triggers
	fmt.Printf("🔄 Auto-next: Preparing next episode...\n")
	
	nextEpisode := currentEpisode + 1
	nextSeason := currentSeason

	// Check if we need to move to the next season
	if nextEpisode > totalEpisodes {
		nextSeason++
		nextEpisode = 1
		
		// Try to load next season details
		cacheKey := fmt.Sprintf("%d-%d", currentTVShowID, nextSeason)
		if seasonDetails, exists := seasonDetailsCache[cacheKey]; exists && len(seasonDetails.Episodes) > 0 {
			totalEpisodes = len(seasonDetails.Episodes)
		} else {
			// Try to fetch next season
			go func() {
				season, err := api.GetSeasonDetails(currentTVShowID, nextSeason)
				if err != nil || len(season.Episodes) == 0 {
					// No more episodes
					fmt.Printf("Auto-next: End of series\n")
					return
				}
				cacheKey := fmt.Sprintf("%d-%d", currentTVShowID, nextSeason)
				seasonDetailsCache[cacheKey] = season
				totalEpisodes = len(season.Episodes)
				currentSeason = nextSeason
				currentEpisode = nextEpisode
				
				// Play first episode of next season
				if len(season.Episodes) > 0 {
					ep := season.Episodes[0]
					fmt.Printf("▶ Auto-next: Will move to next season - S%dE%d - %s\n", nextSeason, ep.EpisodeNumber, ep.Name)
					
					// Show countdown with new season indication (includes season info in episode name)
					newSeasonEpisodeName := fmt.Sprintf("🎬 NEW SEASON! - %s", ep.Name)
					showAutoNextCountdown(currentTVShowID, nextSeason, ep.EpisodeNumber, currentShowName, newSeasonEpisodeName)
				} else {
					// No episodes in next season
					autoNextInProgress = false
					fyne.Do(func() {
						dialog.ShowInformation("Auto-Next", "End of series reached.", currentWindow)
					})
				}
			}()
			return
		}
	}

	// Update tracking
	currentSeason = nextSeason
	currentEpisode = nextEpisode

	// Get episode details from cache
	cacheKey := fmt.Sprintf("%d-%d", currentTVShowID, currentSeason)
	if seasonDetails, exists := seasonDetailsCache[cacheKey]; exists {
		for _, ep := range seasonDetails.Episodes {
			if ep.EpisodeNumber == nextEpisode {
				fmt.Printf("▶ Auto-next: Will play S%dE%d - %s in 5 seconds...\n", nextSeason, nextEpisode, ep.Name)
				
				// Show countdown notification with cancel option
				showAutoNextCountdown(currentTVShowID, nextSeason, nextEpisode, currentShowName, ep.Name)
				return
			}
		}
	}

	fmt.Printf("⚠ Auto-next: Could not find next episode\n")
	autoNextInProgress = false
	
	// Try to play from queue instead
	q := queue.Get()
	if !q.IsEmpty() {
		fmt.Println("📋 Auto-next ended, checking queue...")
		playNextInQueue()
	} else {
		// Show notification that auto-next couldn't continue
		fyne.Do(func() {
			dialog.ShowInformation("Auto-Next", "No more episodes available.", currentWindow)
		})
	}
}

// showAutoNextCountdown shows a countdown dialog with cancel option before auto-next
func showAutoNextCountdown(tvID, season, episode int, showName, episodeName string) {
	cancelled := false
	
	fyne.Do(func() {
		// Create countdown message
		countdownMsg := fmt.Sprintf("Next episode will start in 5 seconds...\n\nS%dE%d - %s", season, episode, episodeName)
		countdownLabel := widget.NewLabel(countdownMsg)
		countdownLabel.Alignment = fyne.TextAlignCenter
		
		// Create custom dialog with Cancel button
		var customDialog *dialog.CustomDialog
		
		cancelBtn := widget.NewButton("Cancel Auto-Next", func() {
			cancelled = true
			autoNextCancelled = true
			autoNextInProgress = false
			customDialog.Hide()
			dialog.ShowInformation("Auto-Next Cancelled", "Auto-next has been disabled for this session.\n\nYou can re-enable it in Settings.", currentWindow)
		})
		
		continueBtn := widget.NewButton("Play Now", func() {
			cancelled = true // Stop countdown
			customDialog.Hide()
			fmt.Printf("▶ Auto-next: Playing S%dE%d - %s (manual trigger)\n", season, episode, episodeName)
			watchEpisodeWithoutDialog(tvID, season, episode, showName, episodeName)
		})
		
		content := container.NewVBox(
			countdownLabel,
			widget.NewSeparator(),
			container.NewGridWithColumns(2, continueBtn, cancelBtn),
		)
		
		customDialog = dialog.NewCustom("Auto-Next Episode", "", content, currentWindow)
		customDialog.Show()
		
		// Start countdown
		go func() {
			for i := 5; i > 0; i-- {
				if cancelled {
					return
				}
				
				// Update countdown
				fyne.Do(func() {
					if !cancelled {
						countdownMsg := fmt.Sprintf("Next episode will start in %d second(s)...\n\nS%dE%d - %s", i, season, episode, episodeName)
						countdownLabel.SetText(countdownMsg)
					}
				})
				
				// Wait 1 second
				time.Sleep(1 * time.Second)
			}
			
			// If not cancelled, play next episode
			if !cancelled {
				fyne.Do(func() {
					customDialog.Hide()
				})
				
				fmt.Printf("▶ Auto-next: Playing S%dE%d - %s\n", season, episode, episodeName)
				watchEpisodeWithoutDialog(tvID, season, episode, showName, episodeName)
			} else {
				autoNextInProgress = false
			}
		}()
	})
}

// downloadEpisode downloads an episode
func downloadEpisode(tvID, season, episode int, showName, episodeName string) {
	progress := dialog.NewProgressInfinite("Loading", "Fetching stream URL...", currentWindow)
	progress.Show()

	go func() {
		streamInfo, err := api.GetStreamURL(tvID, "tv", season, episode)
		if err != nil {
			fyne.Do(func() {
				progress.Hide()
				dialog.ShowError(fmt.Errorf("failed to get stream: %v", err), currentWindow)
			})
			return
		}

		filename := fmt.Sprintf("%s_S%02dE%02d_%s.m3u8", showName, season, episode, episodeName)
		
		err = downloader.DownloadStream(streamInfo.StreamURL, filename, func(downloaded, total int64) {
			// Update progress
		})

		fyne.Do(func() {
			progress.Hide()

			if err != nil {
				dialog.ShowError(fmt.Errorf("download failed: %v", err), currentWindow)
				return
			}

			downloadPath := downloader.GetDownloadPath()
			dialog.ShowInformation("Stream URL Saved!", 
				fmt.Sprintf("Stream URL saved to:\n%s\n\nThe file contains the stream URL and download instructions.\n\nTo download the actual video, use ffmpeg:\nffmpeg -i \"[URL]\" -c copy output.mp4\n\nOr just click 'Watch' to play it in MPV!", downloadPath), 
				currentWindow)
		})
	}()
}

