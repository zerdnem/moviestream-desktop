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

var (
	currentWindow    fyne.Window
	mainContainer    *fyne.Container
	searchEntry      *widget.Entry
	contentTypeRadio *widget.RadioGroup
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

	// MPV status check
	mpvStatus := widget.NewLabel("")
	if player.CheckMPVInstalled() {
		mpvStatus.SetText("✓ MPV Player detected")
		mpvStatus.Importance = widget.SuccessImportance
	} else {
		mpvStatus.SetText("⚠ MPV Player not found - Install from https://mpv.io/")
		mpvStatus.Importance = widget.WarningImportance
	}

	// Search container
	searchContainer := container.NewVBox(
		container.NewCenter(title),
		widget.NewSeparator(),
		widget.NewLabel("Select content type:"),
		contentTypeRadio,
		widget.NewLabel("Enter search query:"),
		searchEntry,
		searchBtn,
		widget.NewSeparator(),
		mpvStatus,
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
				resultsWidget = createMovieResults(movies)
			}
		} else {
			tvShows, searchErr := api.SearchTVShows(query)
			err = searchErr
			if err == nil {
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

	// Back button
	backBtn := widget.NewButton("← Back to Search", func() {
		CreateMainUI(currentWindow)
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
	)

	scrollContent := container.NewVScroll(content)
	currentWindow.SetContent(scrollContent)
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

		if err := player.PlayWithMPV(streamInfo.StreamURL, movie.Title, subtitleURLs); err != nil {
			dialog.ShowError(err, currentWindow)
			return
		}

		subMsg := ""
		if len(subtitleURLs) > 0 {
			subMsg = fmt.Sprintf(" with %d subtitle track(s)", len(subtitleURLs))
		}
		dialog.ShowInformation("Success", fmt.Sprintf("Playing '%s' in MPV%s", movie.Title, subMsg), currentWindow)
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

