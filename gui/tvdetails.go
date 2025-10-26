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
	// Modern back button
	backBtn := CreateIconButton("Back", IconBack, func() {
		GoBackToSearch()
	})
	backBtn.Importance = widget.LowImportance

	// Modern title with accent color - make it more prominent on backdrop
	titleText := canvas.NewText(tvDetails.Name, GetTextColor())
	titleText.TextSize = 28
	titleText.Alignment = fyne.TextAlignCenter
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	
	// Add a subtle glow effect with a semi-transparent background
	cardColor := GetCardColor()
	r, g, b, _ := cardColor.RGBA()
	glowBg := canvas.NewRectangle(&color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: 180,
	})
	glowContainer := container.NewStack(
		glowBg,
		container.NewPadded(titleText),
	)

	// Compact info with icons
	rating := fmt.Sprintf("★ %.1f/10", tvDetails.VoteAverage)
	if tvDetails.VoteAverage == 0 {
		rating = "★ N/A"
	}
	infoText := fmt.Sprintf("%s | %s | %d Seasons",
		rating, tvDetails.FirstAirDate, len(tvDetails.Seasons))
	infoLabel := widget.NewLabel(infoText)
	infoLabel.Alignment = fyne.TextAlignCenter

	// Overview with background for readability
	overviewLabel := widget.NewLabel(tvDetails.Overview)
	overviewLabel.Wrapping = fyne.TextWrapWord
	
	cardColor2 := GetCardColor()
	r2, g2, b2, _ := cardColor2.RGBA()
	overviewBg := canvas.NewRectangle(&color.NRGBA{
		R: uint8(r2 >> 8),
		G: uint8(g2 >> 8),
		B: uint8(b2 >> 8),
		A: 200,
	})
	overviewContainer := container.NewStack(
		overviewBg,
		container.NewPadded(overviewLabel),
	)

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
			container.NewCenter(glowContainer),
			container.NewCenter(infoLabel),
			widget.NewSeparator(),
			CreateHeader("Overview"),
			overviewContainer,
			widget.NewSeparator(),
			widget.NewLabel("No seasons available"),
		)
		
		// Use parallax background
		backdropURL := api.GetBackdropURL(tvDetails.BackdropPath)
		parallaxView := CreateParallaxView(backdropURL, content)
		currentWindow.SetContent(parallaxView)
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

	// Compact content layout
	content := container.NewVBox(
		backBtn,
		widget.NewSeparator(),
		container.NewCenter(glowContainer),
		container.NewCenter(infoLabel),
		widget.NewSeparator(),
		CreateHeader("Overview"),
		overviewContainer,
		widget.NewSeparator(),
		CreateHeader("Episodes"),
		seasonSelect,
		widget.NewSeparator(),
		episodesList,
	)

	// Use parallax background with backdrop image
	backdropURL := api.GetBackdropURL(tvDetails.BackdropPath)
	parallaxView := CreateParallaxView(backdropURL, content)
	currentWindow.SetContent(parallaxView)
	
	// Load first season after UI is set up
	if len(seasonOptions) > 0 {
		seasonSelect.SetSelected(seasonOptions[0])
	}
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
				
				// Modern compact episode card
				episodeTitle := fmt.Sprintf("E%d: %s", ep.EpisodeNumber, ep.Name)
				titleLabel := widget.NewLabelWithStyle(episodeTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
				
				// Truncate overview for compact display
				overview := ep.Overview
				if len(overview) > 120 {
					overview = overview[:117] + "..."
				}
				overviewLabel := widget.NewLabel(overview)
				overviewLabel.Wrapping = fyne.TextWrapWord

				// Modern action buttons
				watchBtn := CreateIconButtonWithImportance("Watch", IconPlay, widget.HighImportance, func() {
					watchEpisodeWithAutoNext(tvID, seasonNum, ep.EpisodeNumber, showName, ep.Name, len(season.Episodes))
				})

				downloadBtn := CreateIconButton("Download", IconDownload, func() {
					downloadEpisode(tvID, seasonNum, ep.EpisodeNumber, showName, ep.Name)
				})

				addToQueueBtn := CreateIconButton("Queue", IconAdd, func() {
					addEpisodeToQueue(tvID, showName, seasonNum, ep.EpisodeNumber, ep.Name)
				})
				addToQueueBtn.Importance = widget.LowImportance

				// Compact button layout
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
	fmt.Printf("✓ Added to queue: %s - S%dE%d: %s (Queue size: %d)\n", showName, season, episode, episodeName, q.Size())
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

		// Print status to console instead of showing dialog
		fmt.Printf("\n▶ Playing: %s - S%dE%d - %s\n", showName, season, episode, episodeName)
		if len(subtitleURLs) > 0 {
			fmt.Printf("   ✓ %d subtitle track(s) loaded\n", len(subtitleURLs))
		} else {
			fmt.Printf("   ⚠ No subtitles found\n")
		}
		
		if userSettings.AutoNext && !autoNextCancelled {
			fmt.Printf("   ▶ Auto-next enabled - Next episode will play automatically\n")
		} else if autoNextCancelled {
			fmt.Printf("   ⏸ Auto-next disabled for this session\n")
		}
		
		if !q.IsEmpty() {
			fmt.Printf("   📋 Queue: %d item(s) waiting\n", q.Size())
		}
		fmt.Println()
	}()
}

// playNextEpisode automatically plays the next episode
func playNextEpisode() {
	userSettings := settings.Get()
	if !userSettings.AutoNext || autoNextCancelled || autoNextInProgress {
		return
	}

	autoNextInProgress = true // Prevent multiple triggers
	fmt.Printf("🔄 Auto-next: Preparing next episode (Current: S%dE%d)...\n", currentSeason, currentEpisode)
	
	// Get current season details from cache
	cacheKey := fmt.Sprintf("%d-%d", currentTVShowID, currentSeason)
	seasonDetails, exists := seasonDetailsCache[cacheKey]
	
	if !exists || len(seasonDetails.Episodes) == 0 {
		fmt.Printf("⚠ Auto-next: Season details not found in cache\n")
		autoNextInProgress = false
		return
	}

	// Find the current episode in the list
	currentEpisodeIndex := -1
	for i, ep := range seasonDetails.Episodes {
		if ep.EpisodeNumber == currentEpisode {
			currentEpisodeIndex = i
			break
		}
	}

	if currentEpisodeIndex == -1 {
		fmt.Printf("⚠ Auto-next: Could not find current episode S%dE%d in episode list\n", currentSeason, currentEpisode)
		autoNextInProgress = false
		return
	}

	// Check if there's a next episode in the current season
	if currentEpisodeIndex + 1 < len(seasonDetails.Episodes) {
		// Play next episode in current season
		nextEp := seasonDetails.Episodes[currentEpisodeIndex + 1]
		currentEpisode = nextEp.EpisodeNumber
		
		fmt.Printf("▶ Auto-next: Next episode in season - S%dE%d - %s\n", currentSeason, nextEp.EpisodeNumber, nextEp.Name)
		showAutoNextCountdown(currentTVShowID, currentSeason, nextEp.EpisodeNumber, currentShowName, nextEp.Name)
		return
	}

	// No more episodes in current season, try next season
	fmt.Printf("🔄 Auto-next: End of season %d, checking for season %d...\n", currentSeason, currentSeason+1)
	nextSeason := currentSeason + 1
	
	// Try to load next season details
	nextSeasonKey := fmt.Sprintf("%d-%d", currentTVShowID, nextSeason)
	if nextSeasonDetails, exists := seasonDetailsCache[nextSeasonKey]; exists && len(nextSeasonDetails.Episodes) > 0 {
		// Next season is cached, play first episode
		firstEp := nextSeasonDetails.Episodes[0]
		currentSeason = nextSeason
		currentEpisode = firstEp.EpisodeNumber
		totalEpisodes = len(nextSeasonDetails.Episodes)
		
		fmt.Printf("▶ Auto-next: Moving to next season - S%dE%d - %s\n", nextSeason, firstEp.EpisodeNumber, firstEp.Name)
		newSeasonEpisodeName := fmt.Sprintf("🎬 NEW SEASON! - %s", firstEp.Name)
		showAutoNextCountdown(currentTVShowID, nextSeason, firstEp.EpisodeNumber, currentShowName, newSeasonEpisodeName)
		return
	}

	// Next season not cached, try to fetch it
	go func() {
		season, err := api.GetSeasonDetails(currentTVShowID, nextSeason)
		if err != nil || len(season.Episodes) == 0 {
			// No more seasons/episodes
			fmt.Printf("⚠ Auto-next: No more seasons available (tried season %d)\n", nextSeason)
			autoNextInProgress = false
			
			// Try to play from queue instead
			q := queue.Get()
			if !q.IsEmpty() {
				fmt.Println("📋 Auto-next ended, checking queue...")
				playNextInQueue()
			} else {
				fyne.Do(func() {
					dialog.ShowInformation("Auto-Next", "End of series reached.", currentWindow)
				})
			}
			return
		}
		
		// Cache the next season
		seasonDetailsCache[nextSeasonKey] = season
		
		// Play first episode of next season
		firstEp := season.Episodes[0]
		currentSeason = nextSeason
		currentEpisode = firstEp.EpisodeNumber
		totalEpisodes = len(season.Episodes)
		
		fmt.Printf("▶ Auto-next: Fetched next season - S%dE%d - %s\n", nextSeason, firstEp.EpisodeNumber, firstEp.Name)
		newSeasonEpisodeName := fmt.Sprintf("🎬 NEW SEASON! - %s", firstEp.Name)
		showAutoNextCountdown(currentTVShowID, nextSeason, firstEp.EpisodeNumber, currentShowName, newSeasonEpisodeName)
	}()
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
			dialog.ShowInformation("Auto-Next Cancelled", "Auto-next has been disabled for this session.\n\nIt will automatically re-enable when you start watching a new show.", currentWindow)
		})
		
		continueBtn := widget.NewButton("Play Now", func() {
			cancelled = true // Stop countdown
			customDialog.Hide()
			fmt.Printf("▶ Auto-next: Playing S%dE%d - %s (manual trigger)\n", season, episode, episodeName)
			// Update tracking variables before playing
			currentTVShowID = tvID
			currentShowName = showName
			currentSeason = season
			currentEpisode = episode
			autoNextInProgress = false // Reset for next auto-next
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
				// Update tracking variables before playing
				currentTVShowID = tvID
				currentShowName = showName
				currentSeason = season
				currentEpisode = episode
				autoNextInProgress = false // Reset for next auto-next
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

