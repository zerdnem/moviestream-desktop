package gui

import (
	"moviestream-gui/player"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// GoBackToSearch returns to the search results if they exist, otherwise returns to main UI
func GoBackToSearch() {
	if lastSearchQuery != "" && (len(lastSearchMovies) > 0 || len(lastSearchTVShows) > 0) {
		// Recreate the main UI with search results
		createMainUIWithResults()
	} else {
		// No previous search, go to main UI
		CreateMainUI(currentWindow)
	}
}

// createMainUIWithResults recreates the main UI and populates it with the last search results
func createMainUIWithResults() {
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

	// Search input - restore previous search
	if searchEntry == nil {
		searchEntry = widget.NewEntry()
	}
	searchEntry.SetText(lastSearchQuery)
	searchEntry.SetPlaceHolder("Search thousands of movies and TV shows...")
	searchEntry.OnSubmitted = func(query string) {
		performSearch(query)
	}

	// Content type selector - restore previous selection
	if contentTypeRadio == nil {
		contentTypeRadio = widget.NewRadioGroup([]string{"Movies", "TV Shows"}, nil)
		contentTypeRadio.Horizontal = true
	}
	// Default to Movies if not set
	if lastSearchContentType == "" {
		lastSearchContentType = "Movies"
	}
	contentTypeRadio.SetSelected(lastSearchContentType)

	// View mode selector - restore previous selection
	if viewModeRadio == nil {
		viewModeRadio = widget.NewRadioGroup([]string{"Grid View", "List View"}, func(selected string) {
			// Re-render results with new view mode
			if lastSearchQuery != "" && (len(lastSearchMovies) > 0 || len(lastSearchTVShows) > 0) {
				lastViewMode = selected
				refreshSearchResults()
			}
		})
		viewModeRadio.Horizontal = true
	}
	if lastViewMode == "" {
		lastViewMode = "Grid View"
	}
	viewModeRadio.SetSelected(lastViewMode)

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

	// Recreate results widget
	var resultsWidget fyne.CanvasObject
	if lastSearchContentType == "Movies" && len(lastSearchMovies) > 0 {
		resultsWidget = createMovieResults(lastSearchMovies)
	} else if lastSearchContentType == "TV Shows" && len(lastSearchTVShows) > 0 {
		resultsWidget = createTVResults(lastSearchTVShows)
	} else if len(lastSearchMovies) > 0 {
		// Handle old format
		resultsWidget = createMovieResults(lastSearchMovies)
	} else if len(lastSearchTVShows) > 0 {
		resultsWidget = createTVResults(lastSearchTVShows)
	}

	// Results scroll container
	resultsScroll := container.NewVScroll(resultsWidget)

	// Main layout with generous padding
	mainContainer = container.NewBorder(
		container.NewPadded(topSection),
		nil,
		nil,
		nil,
		resultsScroll,
	)

	currentWindow.SetContent(mainContainer)
}

