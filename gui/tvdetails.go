package gui

import (
	"fmt"
	"image/color"
	"moviestream-gui/api"
	"moviestream-gui/downloader"
	"moviestream-gui/player"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// showTVDetails shows detailed view for a TV show with episode list
func showTVDetails(tvShow api.TVShow) {
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
		CreateMainUI(currentWindow)
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
					watchEpisode(tvID, seasonNum, ep.EpisodeNumber, showName, ep.Name)
				})

				// Download button
				downloadBtn := widget.NewButton("⬇ Download", func() {
					downloadEpisode(tvID, seasonNum, ep.EpisodeNumber, showName, ep.Name)
				})

				buttonContainer := container.NewGridWithColumns(2, watchBtn, downloadBtn)

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

// watchEpisode starts playing an episode in MPV
func watchEpisode(tvID, season, episode int, showName, episodeName string) {
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
		if err := player.PlayWithMPV(streamInfo.StreamURL, title, subtitleURLs); err != nil {
			dialog.ShowError(err, currentWindow)
			return
		}

		subMsg := ""
		if len(subtitleURLs) > 0 {
			subMsg = fmt.Sprintf(" with %d subtitle track(s)", len(subtitleURLs))
		}
		dialog.ShowInformation("Success", fmt.Sprintf("Playing episode in MPV%s", subMsg), currentWindow)
	}()
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

