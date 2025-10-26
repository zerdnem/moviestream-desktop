package gui

import (
	"fmt"
	"moviestream-gui/api"
	"moviestream-gui/downloader"
	"moviestream-gui/history"
	"moviestream-gui/player"
	"moviestream-gui/queue"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	currentWindow    fyne.Window
	mainContainer    *fyne.Container
	searchEntry      *widget.Entry
	contentTypeRadio *widget.RadioGroup
	// Search state tracking
	lastSearchQuery      string
	lastSearchMovies     []api.Movie
	lastSearchTVShows    []api.TVShow
	lastSearchContentType string
)

// CreateMainUI creates the main user interface
func CreateMainUI(window fyne.Window) {
	currentWindow = window

	// HBO Max style large title
	title := CreateLargeTitle("MovieStream")
	title.Alignment = fyne.TextAlignLeading
	
	// Sleek navigation buttons
	settingsBtn := CreateIconButton("", IconSettings, func() {
		ShowSettingsDialog()
	})
	settingsBtn.Importance = widget.LowImportance
	
	queueBtn := CreateIconButton("", IconQueue, func() {
		ShowQueueView()
	})
	queueBtn.Importance = widget.LowImportance
	
	historyBtn := CreateIconButton("", IconHistory, func() {
		ShowHistoryView()
	})
	historyBtn.Importance = widget.LowImportance
	
	navButtons := container.NewHBox(
		queueBtn,
		historyBtn,
		settingsBtn,
	)
	
	// Spacious header bar
	headerBar := container.NewBorder(
		nil, nil,
		container.NewPadded(title),
		container.NewPadded(navButtons),
	)

	// Modern content type selector
	contentTypeRadio = widget.NewRadioGroup([]string{"Movies", "TV Shows"}, nil)
	contentTypeRadio.SetSelected("Movies")
	contentTypeRadio.Horizontal = true

	// Large search input
	searchEntry = widget.NewEntry()
	searchEntry.SetPlaceHolder("Search thousands of movies and TV shows...")
	searchEntry.OnSubmitted = func(query string) {
		performSearch(query)
	}

	// Compact search button (icon only)
	searchBtn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		query := searchEntry.Text
		if query != "" {
			performSearch(query)
		}
	})
	searchBtn.Importance = widget.HighImportance

	// Compact search bar with inline button
	searchBar := container.NewBorder(nil, nil, nil, searchBtn, searchEntry)

	// Modern search section with styling
	searchSection := container.NewVBox(
		contentTypeRadio,
		searchBar,
	)

	// Minimal player status
	playerStatus := widget.NewLabel("")
	installedPlayers := player.GetInstalledPlayers()
	if len(installedPlayers) > 0 {
		playerNames := ""
		for i, p := range installedPlayers {
			if i > 0 {
				playerNames += ", "
			}
			playerNames += p.Name
		}
		playerStatus.SetText("✓ " + playerNames)
		playerStatus.Importance = widget.SuccessImportance
	} else {
		playerStatus.SetText("⚠ No video player detected")
		playerStatus.Importance = widget.WarningImportance
	}

	// Spacious top section
	topSection := container.NewVBox(
		headerBar,
		widget.NewSeparator(),
		container.NewPadded(searchSection),
		widget.NewSeparator(),
		playerStatus,
	)

	// Compact welcome section
	heroTitle := CreateTitle("Discover Movies & TV Shows")
	heroTitle.Alignment = fyne.TextAlignCenter
	
	resultsContainer := container.NewVBox(
		widget.NewLabel(""),
		container.NewCenter(heroTitle),
	)

	// Scroll container for results
	scrollResults := container.NewVScroll(resultsContainer)

	// Main layout with generous padding
	mainContainer = container.NewBorder(
		container.NewPadded(topSection),
		nil,
		nil,
		nil,
		scrollResults,
	)

	window.SetContent(mainContainer)
}

// performSearch executes the search based on selected content type
func performSearch(query string) {
	if query == "" {
		return
	}

	// Show loading dialog
	progress := dialog.NewProgressInfinite("Searching", fmt.Sprintf("Searching for '%s'...", query), currentWindow)
	progress.Show()

	go func() {
		var err error
		var resultsWidget fyne.CanvasObject

		// Check if searching for movies or TV shows
		isMovies := contentTypeRadio.Selected == "Movies"

		if isMovies {
			movies, searchErr := api.SearchMovies(query)
			err = searchErr
			if err == nil {
				// Store search state
				lastSearchQuery = query
				lastSearchMovies = movies
				lastSearchTVShows = nil
				lastSearchContentType = "Movies"
				
				resultsWidget = createMovieResults(movies)
			}
		} else {
			tvShows, searchErr := api.SearchTVShows(query)
			err = searchErr
			if err == nil {
				// Store search state
				lastSearchQuery = query
				lastSearchMovies = nil
				lastSearchTVShows = tvShows
				lastSearchContentType = "TV Shows"
				
				resultsWidget = createTVResults(tvShows)
			}
		}

		// Always hide progress dialog
		fyne.Do(func() {
			progress.Hide()
		})

		if err != nil {
			dialog.ShowError(fmt.Errorf("search failed: %v", err), currentWindow)
			return
		}

		// Update the main container with results (must be done on UI thread)
		fyne.Do(func() {
			resultsScroll := container.NewVScroll(resultsWidget)
			mainContainer.Objects[0] = resultsScroll
			mainContainer.Refresh()
		})
	}()
}

// createMovieResults creates the results list for movies
func createMovieResults(movies []api.Movie) fyne.CanvasObject {
	if len(movies) == 0 {
		return container.NewCenter(
			container.NewVBox(
				widget.NewLabel(""),
				CreateHeader("No movies found"),
				widget.NewLabel("Try a different search term"),
			),
		)
	}

	var items []fyne.CanvasObject
	
	// HBO Max style section header
	sectionTitle := CreateTitle(fmt.Sprintf("Movies • %d Results", len(movies)))
	items = append(items, 
		widget.NewLabel(""),
		sectionTitle,
		widget.NewLabel(""),
	)

	for _, movie := range movies {
		m := movie // Capture for closure
		
		// Streaming service style card
		titleLabel := widget.NewLabelWithStyle(m.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		
		// Rating and year
		rating := fmt.Sprintf("⭐ %.1f", m.VoteAverage)
		if m.VoteAverage == 0 {
			rating = "⭐ N/A"
		}
		year := ""
		if len(m.ReleaseDate) >= 4 {
			year = m.ReleaseDate[:4]
		}
		infoText := fmt.Sprintf("%s  •  %s", rating, year)
		infoLabel := widget.NewLabel(infoText)
		
		// Compact overview
		overview := m.Overview
		if len(overview) > 120 {
			overview = overview[:117] + "..."
		}
		overviewLabel := widget.NewLabel(overview)
		overviewLabel.Wrapping = fyne.TextWrapWord

		// Compact action buttons
		watchBtn := CreateIconButtonWithImportance("Watch", IconPlay, widget.HighImportance, func() {
			showMovieDetails(m)
		})

		// Compact card with minimal spacing
		cardContent := container.NewVBox(
			titleLabel,
			infoLabel,
			overviewLabel,
			watchBtn,
		)
		
		card := CreateMovieCard(cardContent)

		items = append(items, card)
	}

	return container.NewVBox(items...)
}

// createTVResults creates the results list for TV shows
func createTVResults(tvShows []api.TVShow) fyne.CanvasObject {
	if len(tvShows) == 0 {
		return container.NewCenter(
			container.NewVBox(
				widget.NewLabel(""),
				CreateHeader("No TV shows found"),
				widget.NewLabel("Try a different search term"),
			),
		)
	}

	var items []fyne.CanvasObject
	
	// HBO Max style section header
	sectionTitle := CreateTitle(fmt.Sprintf("TV Shows • %d Results", len(tvShows)))
	items = append(items, 
		widget.NewLabel(""),
		sectionTitle,
		widget.NewLabel(""),
	)

	for _, show := range tvShows {
		s := show // Capture for closure
		
		// Streaming service style card
		titleLabel := widget.NewLabelWithStyle(s.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		
		// Rating and year
		rating := fmt.Sprintf("⭐ %.1f", s.VoteAverage)
		if s.VoteAverage == 0 {
			rating = "⭐ N/A"
		}
		year := ""
		if len(s.FirstAirDate) >= 4 {
			year = s.FirstAirDate[:4]
		}
		infoText := fmt.Sprintf("%s  •  %s", rating, year)
		infoLabel := widget.NewLabel(infoText)
		
		// Compact overview
		overview := s.Overview
		if len(overview) > 120 {
			overview = overview[:117] + "..."
		}
		overviewLabel := widget.NewLabel(overview)
		overviewLabel.Wrapping = fyne.TextWrapWord

		// Compact action buttons
		watchBtn := CreateIconButtonWithImportance("Episodes", IconPlay, widget.HighImportance, func() {
			showTVDetails(s)
		})

		// Compact card with minimal spacing
		cardContent := container.NewVBox(
			titleLabel,
			infoLabel,
			overviewLabel,
			watchBtn,
		)
		
		card := CreateMovieCard(cardContent)

		items = append(items, card)
	}

	return container.NewVBox(items...)
}

// showMovieDetails shows detailed view for a movie with watch/download options
func showMovieDetails(movie api.Movie) {
	// Sleek back button
	backBtn := CreateIconButton(IconBack, "", func() {
		GoBackToSearch()
	})
	backBtn.Importance = widget.LowImportance

	// Hero title section
	titleText := CreateLargeTitle(movie.Title)
	titleText.Alignment = fyne.TextAlignLeading

	// Premium info display
	rating := fmt.Sprintf("⭐ %.1f", movie.VoteAverage)
	if movie.VoteAverage == 0 {
		rating = "⭐ N/A"
	}
	year := ""
	if len(movie.ReleaseDate) >= 4 {
		year = movie.ReleaseDate[:4]
	}
	infoText := fmt.Sprintf("%s  •  %s  •  Movie", rating, year)
	infoLabel := widget.NewLabel(infoText)

	// Overview section
	overviewTitle := CreateHeader("Overview")
	overviewLabel := widget.NewLabel(movie.Overview)
	overviewLabel.Wrapping = fyne.TextWrapWord

	// Premium action buttons - larger and more prominent
	watchBtn := CreateIconButtonWithImportance("Watch Now", IconPlay, widget.HighImportance, func() {
		watchMovie(movie)
	})

	downloadBtn := CreateIconButton("Download", IconDownload, func() {
		downloadMovie(movie)
	})

	addToQueueBtn := CreateIconButton("Add to Queue", IconAdd, func() {
		addMovieToQueue(movie)
	})
	addToQueueBtn.Importance = widget.LowImportance

	// Compact layout
	content := container.NewVBox(
		backBtn,
		widget.NewSeparator(),
		titleText,
		infoLabel,
		watchBtn,
		widget.NewSeparator(),
		overviewTitle,
		overviewLabel,
		widget.NewSeparator(),
		container.NewGridWithColumns(2, downloadBtn, addToQueueBtn),
	)

	scrollContent := container.NewVScroll(content)
	currentWindow.SetContent(container.NewPadded(scrollContent))
}

// addMovieToQueue adds a movie to the playback queue
func addMovieToQueue(movie api.Movie) {
	q := queue.Get()
	q.AddMovie(movie.ID, movie.Title)
	dialog.ShowInformation("Added to Queue", 
		fmt.Sprintf("'%s' has been added to the playback queue.", movie.Title), 
		currentWindow)
}

// watchMovieByID starts playing a movie by ID and title
func watchMovieByID(movieID int, title string) {
	movie := api.Movie{
		ID:    movieID,
		Title: title,
	}
	watchMovie(movie)
}

// watchMovie starts playing a movie in MPV
func watchMovie(movie api.Movie) {
	progress := dialog.NewProgressInfinite("Loading Stream", 
		"Fetching stream URL...\nThis may take 10-15 seconds for browser automation", 
		currentWindow)
	progress.Show()

	go func() {
		streamInfo, err := api.GetStreamURL(movie.ID, "movie", 0, 0)
		
		// Always hide progress
		fyne.Do(func() {
			progress.Hide()
		})

		if err != nil {
			errorMsg := fmt.Sprintf("Failed to get stream:\n%v\n\nTips:\n• This movie might not be available on this platform\n• Try a different, more popular movie\n• Some older or less common titles may not work", err)
			dialog.ShowError(fmt.Errorf(errorMsg), currentWindow)
			return
		}

		// Extract subtitle URLs
		var subtitleURLs []string
		for _, sub := range streamInfo.SubtitleURLs {
			subtitleURLs = append(subtitleURLs, sub.URL)
		}

		// Create callback for queue auto-play
		onEndCallback := func() {
			playNextInQueue()
		}

		if err := player.PlayWithMPVAndCallback(streamInfo.StreamURL, movie.Title, subtitleURLs, onEndCallback); err != nil {
			dialog.ShowError(err, currentWindow)
			return
		}

		// Record in watch history
		h := history.Get()
		h.AddMovie(movie.ID, movie.Title)

		// Show subtitle status
		var statusMsg string
		if len(subtitleURLs) > 0 {
			statusMsg = fmt.Sprintf("✓ Playing '%s'\n\n%d subtitle track(s) loaded successfully", movie.Title, len(subtitleURLs))
		} else {
			statusMsg = fmt.Sprintf("✓ Playing '%s'\n\n⚠ No subtitles found for this content\n\nNote: You can still add external subtitles in your video player", movie.Title)
		}
		
		// Add queue info
		q := queue.Get()
		if !q.IsEmpty() {
			statusMsg += fmt.Sprintf("\n\n📋 %d item(s) in queue - will play next automatically", q.Size())
		}
		
		dialog.ShowInformation("Playback Started", statusMsg, currentWindow)
	}()
}

// downloadMovie downloads a movie
func downloadMovie(movie api.Movie) {
	progress := dialog.NewProgressInfinite("Loading", "Fetching stream URL...", currentWindow)
	progress.Show()

	go func() {
		streamInfo, err := api.GetStreamURL(movie.ID, "movie", 0, 0)
		if err != nil {
			fyne.Do(func() {
				progress.Hide()
				dialog.ShowError(fmt.Errorf("failed to get stream: %v", err), currentWindow)
			})
			return
		}

		filename := fmt.Sprintf("%s.m3u8", movie.Title)
		
		err = downloader.DownloadStream(streamInfo.StreamURL, filename, func(downloaded, total int64) {
			// Update progress (simplified for now)
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

// playNextInQueue plays the next item in the queue
func playNextInQueue() {
	q := queue.Get()
	
	if q.IsEmpty() {
		fmt.Println("Queue is empty, no more items to play")
		return
	}

	nextItem, ok := q.GetNext()
	if !ok {
		return
	}

	fmt.Printf("📋 Playing next item from queue: %s\n", nextItem.GetDisplayTitle())
	
	// Play the next item
	if nextItem.Type == "movie" {
		watchMovieByID(nextItem.TMDBID, nextItem.Title)
	} else {
		watchEpisodeWithAutoNext(nextItem.TMDBID, nextItem.Season, nextItem.Episode, nextItem.Title, nextItem.EpisodeName, 0)
	}
}

