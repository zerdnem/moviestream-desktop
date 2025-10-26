package gui

import (
	"fmt"
	"moviestream-gui/history"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowHistoryView displays the watch history
func ShowHistoryView() {
	h := history.Get()
	items := h.GetAll()

	// Back button
	backBtn := widget.NewButton("← Back to Search", func() {
		GoBackToSearch()
	})

	// Clear history button
	clearBtn := widget.NewButton("Clear History", func() {
		if len(items) == 0 {
			dialog.ShowInformation("History Empty", "The watch history is already empty.", currentWindow)
			return
		}

		dialog.ShowConfirm("Clear History", "Are you sure you want to clear your entire watch history?", func(confirmed bool) {
			if confirmed {
				h.Clear()
				ShowHistoryView() // Refresh view
				dialog.ShowInformation("History Cleared", "All watch history has been removed.", currentWindow)
			}
		}, currentWindow)
	})

	// Continue watching button (for last watched show)
	var continueWatchingBtn *widget.Button
	if lastShow, ok := h.GetLastWatchedShow(); ok {
		continueWatchingBtn = widget.NewButton("▶ Continue Watching", func() {
			// Navigate to the episode or show
			watchEpisodeWithAutoNext(lastShow.TMDBID, lastShow.Season, lastShow.Episode, lastShow.Title, lastShow.EpisodeName, 0)
		})
	}

	// Header
	headerWidgets := []fyne.CanvasObject{
		backBtn,
		widget.NewSeparator(),
		widget.NewLabelWithStyle(fmt.Sprintf("🕒 Watch History (%d items)", len(items)), 
			fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	}

	if continueWatchingBtn != nil {
		headerWidgets = append(headerWidgets, continueWatchingBtn)
	}

	headerWidgets = append(headerWidgets, widget.NewSeparator())

	header := container.NewVBox(headerWidgets...)

	// History items
	var itemWidgets []fyne.CanvasObject

	if len(items) == 0 {
		itemWidgets = append(itemWidgets, 
			widget.NewLabel("No watch history"),
			widget.NewLabel("Start watching movies or TV shows to see them here."),
		)
	} else {
		for i, item := range items {
			idx := i // Capture for closure
			hItem := item // Capture for closure
			
			// Item display
			titleLabel := widget.NewLabelWithStyle(
				hItem.GetDisplayTitle(),
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			)

			timeLabel := widget.NewLabel(fmt.Sprintf("Watched: %s", hItem.GetDisplayTime()))

			// Watch again button
			watchAgainBtn := widget.NewButton("▶ Watch Again", func() {
				if hItem.Type == "movie" {
					watchMovieByID(hItem.TMDBID, hItem.Title)
				} else {
					watchEpisodeWithAutoNext(hItem.TMDBID, hItem.Season, hItem.Episode, hItem.Title, hItem.EpisodeName, 0)
				}
			})

			// Remove button
			removeBtn := widget.NewButton("Remove", func() {
				h.Remove(idx)
				ShowHistoryView() // Refresh view
			})

			buttonContainer := container.NewGridWithColumns(2, watchAgainBtn, removeBtn)

			itemCard := container.NewVBox(
				titleLabel,
				timeLabel,
				buttonContainer,
				widget.NewSeparator(),
			)

			itemWidgets = append(itemWidgets, itemCard)
		}

		// Add clear button at the bottom
		itemWidgets = append(itemWidgets, clearBtn)
	}

	// Content
	content := container.NewVBox(itemWidgets...)
	scrollContent := container.NewVScroll(content)
	scrollContent.SetMinSize(fyne.NewSize(400, 500))

	// Main layout
	mainLayout := container.NewBorder(
		header,
		nil,
		nil,
		nil,
		scrollContent,
	)

	currentWindow.SetContent(mainLayout)
}

