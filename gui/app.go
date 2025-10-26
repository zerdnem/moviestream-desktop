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

	// Title
	title := canvas.NewText("MovieStream", color.RGBA{R: 255, G: 100, B: 100, A: 255})
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Search input
	searchEntry = widget.NewEntry()
	searchEntry.SetPlaceHolder("Search for movies or TV shows...")
	searchEntry.OnSubmitted = func(query string) {
		performSearch(query)
	}

	// Content type selector
	contentTypeRadio = widget.NewRadioGroup([]string{"Movies", "TV Shows"}, nil)
	contentTypeRadio.SetSelected("Movies")
	contentTypeRadio.Horizontal = true

	// Search button
	searchBtn := widget.NewButton("Search", func() {
		query := searchEntry.Text
		if query != "" {
			performSearch(query)
		}
	})

	// Settings button
	settingsBtn := widget.NewButton("⚙ Settings", func() {
		ShowSettingsDialog()
	})

	// Queue button
	queueBtn := widget.NewButton("📋 Queue", func() {
		ShowQueueView()
	})

	// History button
	historyBtn := widget.NewButton("🕒 History", func() {
		ShowHistoryView()
	})

	// Video player status check
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
		playerStatus.SetText(fmt.Sprintf("✓ Video Players: %s", playerNames))
		playerStatus.Importance = widget.SuccessImportance
	} else {
		playerStatus.SetText("⚠ No video player detected - Install MPV, VLC, or another supported player")
		playerStatus.Importance = widget.WarningImportance
	}

	// Search container
	searchContainer := container.NewVBox(
		container.NewCenter(title),
		widget.NewSeparator(),
		widget.NewLabel("Select content type:"),
		contentTypeRadio,
		widget.NewLabel("Enter search query:"),
		searchEntry,
		container.NewGridWithColumns(2, searchBtn, settingsBtn),
		container.NewGridWithColumns(2, queueBtn, historyBtn),
		widget.NewSeparator(),
		playerStatus,
	)

	// Results container (initially empty)
	resultsContainer := container.NewVBox()

	// Scroll container for results
	scrollResults := container.NewVScroll(resultsContainer)
	scrollResults.SetMinSize(fyne.NewSize(400, 500))

	// Main layout
	mainContainer = container.NewBorder(
		searchContainer,
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

		if contentTypeRadio.Selected == "Movies" {
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
			resultsScroll.SetMinSize(fyne.NewSize(400, 500))

			mainContainer.Objects[0] = resultsScroll
			mainContainer.Refresh()
		})
	}()
}

// createMovieResults creates the results list for movies
func createMovieResults(movies []api.Movie) fyne.CanvasObject {
	if len(movies) == 0 {
		return widget.NewLabel("No movies found")
	}

	var items []fyne.CanvasObject
	
	items = append(items, widget.NewLabelWithStyle(
		fmt.Sprintf("Found %d movie(s)", len(movies)),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	))
	items = append(items, widget.NewSeparator())

	for _, movie := range movies {
		m := movie // Capture for closure
		
		// Movie card
		titleLabel := widget.NewLabelWithStyle(m.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		
		infoText := fmt.Sprintf("Release: %s | Rating: %.1f/10", m.ReleaseDate, m.VoteAverage)
		infoLabel := widget.NewLabel(infoText)
		
		overviewLabel := widget.NewLabel(m.Overview)
		overviewLabel.Wrapping = fyne.TextWrapWord

		detailsBtn := widget.NewButton("View Details", func() {
			showMovieDetails(m)
		})

		card := container.NewVBox(
			titleLabel,
			infoLabel,
			overviewLabel,
			detailsBtn,
			widget.NewSeparator(),
		)

		items = append(items, card)
	}

	return container.NewVBox(items...)
}

// createTVResults creates the results list for TV shows
func createTVResults(tvShows []api.TVShow) fyne.CanvasObject {
	if len(tvShows) == 0 {
		return widget.NewLabel("No TV shows found")
	}

	var items []fyne.CanvasObject
	
	items = append(items, widget.NewLabelWithStyle(
		fmt.Sprintf("Found %d TV show(s)", len(tvShows)),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	))
	items = append(items, widget.NewSeparator())

	for _, show := range tvShows {
		s := show // Capture for closure
		
		// TV show card
		titleLabel := widget.NewLabelWithStyle(s.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		
		infoText := fmt.Sprintf("First Air: %s | Rating: %.1f/10", s.FirstAirDate, s.VoteAverage)
		infoLabel := widget.NewLabel(infoText)
		
		overviewLabel := widget.NewLabel(s.Overview)
		overviewLabel.Wrapping = fyne.TextWrapWord

		detailsBtn := widget.NewButton("View Episodes", func() {
			showTVDetails(s)
		})

		card := container.NewVBox(
			titleLabel,
			infoLabel,
			overviewLabel,
			detailsBtn,
			widget.NewSeparator(),
		)

		items = append(items, card)
	}

	return container.NewVBox(items...)
}

// showMovieDetails shows detailed view for a movie with watch/download options
func showMovieDetails(movie api.Movie) {
	// Create details content
	titleText := canvas.NewText(movie.Title, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	titleText.TextSize = 20
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.Alignment = fyne.TextAlignCenter

	infoText := fmt.Sprintf("Release Date: %s\nRating: %.1f/10", movie.ReleaseDate, movie.VoteAverage)
	infoLabel := widget.NewLabel(infoText)

	overviewLabel := widget.NewLabel(movie.Overview)
	overviewLabel.Wrapping = fyne.TextWrapWord

	// Watch button
	watchBtn := widget.NewButton("▶ Watch", func() {
		watchMovie(movie)
	})

	// Download button
	downloadBtn := widget.NewButton("⬇ Download", func() {
		downloadMovie(movie)
	})

	// Add to Queue button
	addToQueueBtn := widget.NewButton("➕ Add to Queue", func() {
		addMovieToQueue(movie)
	})

	// Back button
	backBtn := widget.NewButton("← Back to Search", func() {
		GoBackToSearch()
	})

	content := container.NewVBox(
		backBtn,
		widget.NewSeparator(),
		container.NewCenter(titleText),
		infoLabel,
		widget.NewSeparator(),
		overviewLabel,
		widget.NewSeparator(),
		container.NewGridWithColumns(2, watchBtn, downloadBtn),
		addToQueueBtn,
	)

	scrollContent := container.NewVScroll(content)
	currentWindow.SetContent(scrollContent)
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

