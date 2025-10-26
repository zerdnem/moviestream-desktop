package gui

import (
	"fmt"
	"moviestream-gui/queue"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowQueueView displays the playback queue
func ShowQueueView() {
	q := queue.Get()
	items := q.GetAll()

	// Modern back button
	backBtn := CreateIconButton("Back", IconBack, func() {
		GoBackToSearch()
	})
	backBtn.Importance = widget.LowImportance

	// Modern clear queue button
	clearBtn := CreateIconButton("Clear Queue", IconDelete, func() {
		if q.IsEmpty() {
			dialog.ShowInformation("Queue Empty", "The queue is already empty.", currentWindow)
			return
		}

		dialog.ShowConfirm("Clear Queue", "Are you sure you want to clear the entire queue?", func(confirmed bool) {
			if confirmed {
				q.Clear()
				ShowQueueView() // Refresh view
				dialog.ShowInformation("Queue Cleared", "All items have been removed from the queue.", currentWindow)
			}
		}, currentWindow)
	})
	clearBtn.Importance = widget.DangerImportance

	// Modern header with accent color
	titleContainer := CreateTitleWithIcon(IconQueue, fmt.Sprintf("Queue (%d)", len(items)))

	// Compact header layout
	header := container.NewVBox(
		backBtn,
		widget.NewSeparator(),
		container.NewCenter(titleContainer),
		widget.NewSeparator(),
	)

	// Queue items with modern styling
	var itemWidgets []fyne.CanvasObject

	if len(items) == 0 {
		emptyMsg := container.NewCenter(
			container.NewVBox(
				widget.NewLabel(""),
				CreateHeader("Queue is empty"),
				widget.NewLabel("Add movies or episodes to watch them in sequence"),
			),
		)
		itemWidgets = append(itemWidgets, emptyMsg)
	} else {
		for i, item := range items {
			idx := i // Capture for closure
			qItem := item // Capture for closure
			
			// Compact item display
			titleLabel := widget.NewLabelWithStyle(
				fmt.Sprintf("%d. %s", idx+1, qItem.GetDisplayTitle()),
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			)

			// Modern action buttons
			playNowBtn := CreateIconButtonWithImportance("Play Now", IconPlay, widget.HighImportance, func() {
				q.RemoveAt(idx)
				playQueueItem(&qItem)
			})

			removeBtn := CreateIconButton("Remove", IconRemove, func() {
				q.RemoveAt(idx)
				ShowQueueView() // Refresh view
			})
			removeBtn.Importance = widget.DangerImportance

			buttonContainer := container.NewGridWithColumns(2, playNowBtn, removeBtn)

			// Modern card
			itemCard := container.NewVBox(
				titleLabel,
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

// playQueueItem plays a queue item
func playQueueItem(item *queue.QueueItem) {
	if item.Type == "movie" {
		// Import api package types (we'll need to fetch movie details)
		watchMovieByID(item.TMDBID, item.Title)
	} else {
		// TV episode
		watchEpisodeWithAutoNext(item.TMDBID, item.Season, item.Episode, item.Title, item.EpisodeName, 0)
	}
}

