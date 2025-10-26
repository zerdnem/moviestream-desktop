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

	// Modern back button
	backBtn := CreateIconButton("Back", IconBack, func() {
		GoBackToSearch()
	})
	backBtn.Importance = widget.LowImportance

	// Modern clear history button
	clearBtn := CreateIconButton("Clear History", IconDelete, func() {
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
	clearBtn.Importance = widget.DangerImportance

	// Modern continue watching button (for last watched show)
	var continueWatchingBtn *widget.Button
	if lastShow, ok := h.GetLastWatchedShow(); ok {
		continueWatchingBtn = CreateIconButtonWithImportance("Continue Watching", IconPlay, widget.HighImportance, func() {
			// Navigate to the episode or show
			watchEpisodeWithAutoNext(lastShow.TMDBID, lastShow.Season, lastShow.Episode, lastShow.Title, lastShow.EpisodeName, 0)
		})
	}

	// Modern header with accent color
	titleContainer := CreateTitleWithIcon(IconHistory, fmt.Sprintf("History (%d)", len(items)))

	// Compact header layout
	headerWidgets := []fyne.CanvasObject{
		backBtn,
		widget.NewSeparator(),
		container.NewCenter(titleContainer),
	}

	if continueWatchingBtn != nil {
		headerWidgets = append(headerWidgets, continueWatchingBtn)
	}

	headerWidgets = append(headerWidgets, widget.NewSeparator())

	header := container.NewVBox(headerWidgets...)

	// History items with modern styling
	var itemWidgets []fyne.CanvasObject

	if len(items) == 0 {
		emptyMsg := container.NewCenter(
			container.NewVBox(
				widget.NewLabel(""),
				CreateHeader("No watch history"),
				widget.NewLabel("Start watching movies or TV shows to see them here"),
			),
		)
		itemWidgets = append(itemWidgets, emptyMsg)
	} else {
		for i, item := range items {
			idx := i // Capture for closure
			hItem := item // Capture for closure
			
			// Compact item display
			titleLabel := widget.NewLabelWithStyle(
				hItem.GetDisplayTitle(),
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			)

			// Compact time display
			timeLabel := widget.NewLabel(hItem.GetDisplayTime())

			// Modern action buttons
			watchAgainBtn := CreateIconButtonWithImportance("Watch Again", IconPlay, widget.HighImportance, func() {
				if hItem.Type == "movie" {
					watchMovieByID(hItem.TMDBID, hItem.Title)
				} else {
					watchEpisodeWithAutoNext(hItem.TMDBID, hItem.Season, hItem.Episode, hItem.Title, hItem.EpisodeName, 0)
				}
			})

			removeBtn := CreateIconButton("Remove", IconRemove, func() {
				h.Remove(idx)
				ShowHistoryView() // Refresh view
			})
			removeBtn.Importance = widget.DangerImportance

			buttonContainer := container.NewGridWithColumns(2, watchAgainBtn, removeBtn)

			// Modern card
			itemCard := container.NewVBox(
				titleLabel,
				timeLabel,
				buttonContainer,
				widget.NewSeparator(),
			)

			itemWidgets = append(itemWidgets, itemCard)
		}

		// Add clear button at the bottom with spacing
		itemWidgets = append(itemWidgets, widget.NewSeparator(), clearBtn)
	}

	// Content with modern layout
	content := container.NewVBox(itemWidgets...)
	scrollContent := container.NewVScroll(content)

	// Main layout with padding
	mainLayout := container.NewBorder(
		container.NewPadded(header),
		nil,
		nil,
		nil,
		scrollContent,
	)

	currentWindow.SetContent(mainLayout)
}

