package gui

import (
	"fmt"
	"image/color"
	"moviestream-gui/player"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
	// Title
	title := canvas.NewText("MovieStream", color.RGBA{R: 255, G: 100, B: 100, A: 255})
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Search input - restore previous search
	if searchEntry == nil {
		searchEntry = widget.NewEntry()
	}
	searchEntry.SetText(lastSearchQuery)
	searchEntry.SetPlaceHolder("Search for movies or TV shows...")
	searchEntry.OnSubmitted = func(query string) {
		performSearch(query)
	}

	// Content type selector - restore previous selection
	if contentTypeRadio == nil {
		contentTypeRadio = widget.NewRadioGroup([]string{"Movies", "TV Shows"}, nil)
		contentTypeRadio.Horizontal = true
	}
	contentTypeRadio.SetSelected(lastSearchContentType)

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

	// Video player status (simplified for back navigation)
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

	// Recreate results widget
	var resultsWidget fyne.CanvasObject
	if lastSearchContentType == "Movies" && len(lastSearchMovies) > 0 {
		resultsWidget = createMovieResults(lastSearchMovies)
	} else if lastSearchContentType == "TV Shows" && len(lastSearchTVShows) > 0 {
		resultsWidget = createTVResults(lastSearchTVShows)
	}

	// Results scroll container
	resultsScroll := container.NewVScroll(resultsWidget)
	resultsScroll.SetMinSize(fyne.NewSize(400, 500))

	// Main layout
	mainContainer = container.NewBorder(
		searchContainer,
		nil,
		nil,
		nil,
		resultsScroll,
	)

	currentWindow.SetContent(mainContainer)
}

