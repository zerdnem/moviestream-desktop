package gui

import (
	"fmt"
	"image/color"
	"moviestream-gui/api"
	"moviestream-gui/downloader"
	"moviestream-gui/history"
	"moviestream-gui/player"
	"moviestream-gui/queue"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	viewModeRadio    *widget.RadioGroup
	// Search state tracking
	lastSearchQuery       string
	lastSearchMovies      []api.Movie
	lastSearchTVShows     []api.TVShow
	lastSearchContentType string
	lastViewMode          string
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

	// Initialize defaults if not set
	if lastSearchContentType == "" {
		lastSearchContentType = "Movies"
	}
	if lastViewMode == "" {
		lastViewMode = "Grid View"
	}

	// Modern content type selector
	contentTypeRadio = widget.NewRadioGroup([]string{"Movies", "TV Shows"}, nil)
	contentTypeRadio.SetSelected(lastSearchContentType)
	contentTypeRadio.Horizontal = true

	// View mode selector
	viewModeRadio = widget.NewRadioGroup([]string{"Grid View", "List View"}, func(selected string) {
		// Re-render results with new view mode
		if lastSearchQuery != "" && (len(lastSearchMovies) > 0 || len(lastSearchTVShows) > 0) {
			lastViewMode = selected
			refreshSearchResults()
		}
	})
	viewModeRadio.SetSelected(lastViewMode)
	viewModeRadio.Horizontal = true

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
		viewModeRadio,
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

// refreshSearchResults re-renders the search results with the current view mode
func refreshSearchResults() {
	var resultsWidget fyne.CanvasObject
	
	if lastSearchContentType == "Movies" && len(lastSearchMovies) > 0 {
		resultsWidget = createMovieResults(lastSearchMovies)
	} else if lastSearchContentType == "TV Shows" && len(lastSearchTVShows) > 0 {
		resultsWidget = createTVResults(lastSearchTVShows)
	}
	
	if resultsWidget != nil {
		resultsScroll := container.NewVScroll(resultsWidget)
		mainContainer.Objects[0] = resultsScroll
		mainContainer.Refresh()
	}
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

	// Check view mode (default to Grid View if not set)
	if lastViewMode == "" || lastViewMode == "Grid View" {
		return createMovieGridView(movies)
	}
	return createMovieListView(movies)
}

// createMovieGridView creates a grid view with cover images
func createMovieGridView(movies []api.Movie) fyne.CanvasObject {
	var items []fyne.CanvasObject
	
	// HBO Max style section header
	sectionTitle := CreateTitle(fmt.Sprintf("Movies • %d Results", len(movies)))
	items = append(items, 
		widget.NewLabel(""),
		sectionTitle,
		widget.NewLabel(""),
	)

	// Create grid container (3 columns)
	var gridItems []fyne.CanvasObject
	for _, movie := range movies {
		m := movie // Capture for closure
		
		// Load poster image
		posterURL := api.GetPosterURL(m.PosterPath)
		posterImg := LoadImageFromURL(posterURL, 150, 225)
		posterImg.FillMode = canvas.ImageFillContain
		
		// Title
		titleLabel := widget.NewLabelWithStyle(m.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		titleLabel.Wrapping = fyne.TextWrapWord
		
		// Rating
		rating := fmt.Sprintf("⭐ %.1f", m.VoteAverage)
		if m.VoteAverage == 0 {
			rating = "⭐ N/A"
		}
		ratingLabel := widget.NewLabel(rating)
		ratingLabel.Alignment = fyne.TextAlignCenter
		
		// Create card content
		cardContent := container.NewVBox(
			posterImg,
			titleLabel,
			ratingLabel,
		)
		
		card := CreateMovieCard(cardContent)
		
		// Wrap in tappable container (no visible button styling)
		itemCard := NewTappableCard(card, func() {
			showMovieDetails(m)
		})
		
		gridItems = append(gridItems, itemCard)
	}
	
	grid := container.NewGridWithColumns(3, gridItems...)
	items = append(items, grid)
	
	return container.NewVBox(items...)
}

// createMovieListView creates a list view with cover images
func createMovieListView(movies []api.Movie) fyne.CanvasObject {
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
		
		// Load poster image
		posterURL := api.GetPosterURL(m.PosterPath)
		posterImg := LoadImageFromURL(posterURL, 100, 150)
		posterImg.FillMode = canvas.ImageFillContain
		
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

		// Text content
		textContent := container.NewVBox(
			titleLabel,
			infoLabel,
			overviewLabel,
			watchBtn,
		)
		
		// Combine poster and text
		cardContent := container.NewBorder(nil, nil, posterImg, nil, textContent)
		
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

	// Check view mode (default to Grid View if not set)
	if lastViewMode == "" || lastViewMode == "Grid View" {
		return createTVGridView(tvShows)
	}
	return createTVListView(tvShows)
}

// createTVGridView creates a grid view with cover images
func createTVGridView(tvShows []api.TVShow) fyne.CanvasObject {
	var items []fyne.CanvasObject
	
	// HBO Max style section header
	sectionTitle := CreateTitle(fmt.Sprintf("TV Shows • %d Results", len(tvShows)))
	items = append(items, 
		widget.NewLabel(""),
		sectionTitle,
		widget.NewLabel(""),
	)

	// Create grid container (3 columns)
	var gridItems []fyne.CanvasObject
	for _, show := range tvShows {
		s := show // Capture for closure
		
		// Load poster image
		posterURL := api.GetPosterURL(s.PosterPath)
		posterImg := LoadImageFromURL(posterURL, 150, 225)
		posterImg.FillMode = canvas.ImageFillContain
		
		// Title
		titleLabel := widget.NewLabelWithStyle(s.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		titleLabel.Wrapping = fyne.TextWrapWord
		
		// Rating
		rating := fmt.Sprintf("⭐ %.1f", s.VoteAverage)
		if s.VoteAverage == 0 {
			rating = "⭐ N/A"
		}
		ratingLabel := widget.NewLabel(rating)
		ratingLabel.Alignment = fyne.TextAlignCenter
		
		// Create card content
		cardContent := container.NewVBox(
			posterImg,
			titleLabel,
			ratingLabel,
		)
		
		card := CreateMovieCard(cardContent)
		
		// Wrap in tappable container (no visible button styling)
		itemCard := NewTappableCard(card, func() {
			showTVDetails(s)
		})
		
		gridItems = append(gridItems, itemCard)
	}
	
	grid := container.NewGridWithColumns(3, gridItems...)
	items = append(items, grid)
	
	return container.NewVBox(items...)
}

// createTVListView creates a list view with cover images
func createTVListView(tvShows []api.TVShow) fyne.CanvasObject {
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
		
		// Load poster image
		posterURL := api.GetPosterURL(s.PosterPath)
		posterImg := LoadImageFromURL(posterURL, 100, 150)
		posterImg.FillMode = canvas.ImageFillContain
		
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

		// Text content
		textContent := container.NewVBox(
			titleLabel,
			infoLabel,
			overviewLabel,
			watchBtn,
		)
		
		// Combine poster and text
		cardContent := container.NewBorder(nil, nil, posterImg, nil, textContent)
		
		card := CreateMovieCard(cardContent)

		items = append(items, card)
	}

	return container.NewVBox(items...)
}

// showMovieDetails shows detailed view for a movie with watch/download options
func showMovieDetails(movie api.Movie) {
	// Sleek back button - FIX: correct parameter order (label, icon, func)
	backBtn := CreateIconButton("Back", IconBack, func() {
		GoBackToSearch()
	})
	backBtn.Importance = widget.LowImportance

	// Hero title section with enhanced visibility on backdrop
	titleText := canvas.NewText(movie.Title, GetTextColor())
	titleText.TextSize = 28
	titleText.Alignment = fyne.TextAlignCenter
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	
	// Add a subtle glow effect with a semi-transparent background
	cardColor := GetCardColor()
	r, g, b, _ := cardColor.RGBA()
	titleBg := canvas.NewRectangle(&color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: 180,
	})
	titleContainer := container.NewStack(
		titleBg,
		container.NewPadded(titleText),
	)

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
	infoLabel.Alignment = fyne.TextAlignCenter

	// Overview section with background for readability
	overviewLabel := widget.NewLabel(movie.Overview)
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
		container.NewCenter(titleContainer),
		container.NewCenter(infoLabel),
		watchBtn,
		widget.NewSeparator(),
		CreateHeader("Overview"),
		overviewContainer,
		widget.NewSeparator(),
		container.NewGridWithColumns(2, downloadBtn, addToQueueBtn),
	)

	// Use parallax background with backdrop image
	backdropURL := api.GetBackdropURL(movie.BackdropPath)
	parallaxView := CreateParallaxView(backdropURL, content)
	currentWindow.SetContent(parallaxView)
}

// addMovieToQueue adds a movie to the playback queue
func addMovieToQueue(movie api.Movie) {
	q := queue.Get()
	q.AddMovie(movie.ID, movie.Title)
	fmt.Printf("✓ Added to queue: %s (Queue size: %d)\n", movie.Title, q.Size())
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

		// Print status to console instead of showing dialog
		fmt.Printf("\n▶ Playing: %s\n", movie.Title)
		if len(subtitleURLs) > 0 {
			fmt.Printf("   ✓ %d subtitle track(s) loaded\n", len(subtitleURLs))
		} else {
			fmt.Printf("   ⚠ No subtitles found\n")
		}
		
		q := queue.Get()
		if !q.IsEmpty() {
			fmt.Printf("   📋 Queue: %d item(s) waiting\n", q.Size())
		}
		fmt.Println()
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

