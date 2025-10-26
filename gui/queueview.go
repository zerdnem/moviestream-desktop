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

	// Back button
	backBtn := widget.NewButton("← Back to Search", func() {
		GoBackToSearch()
	})

	// Clear queue button
	clearBtn := widget.NewButton("Clear Queue", func() {
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

	// Header
	header := container.NewVBox(
		backBtn,
		widget.NewSeparator(),
		widget.NewLabelWithStyle(fmt.Sprintf("📋 Playback Queue (%d items)", len(items)), 
			fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: false}),
		widget.NewSeparator(),
	)

	// Queue items
	var itemWidgets []fyne.CanvasObject

	if len(items) == 0 {
		itemWidgets = append(itemWidgets, 
			widget.NewLabel("Queue is empty"),
			widget.NewLabel("Add movies or episodes to the queue to watch them in sequence."),
		)
	} else {
		for i, item := range items {
			idx := i // Capture for closure
			qItem := item // Capture for closure
			
			// Item display
			titleLabel := widget.NewLabelWithStyle(
				fmt.Sprintf("%d. %s", idx+1, qItem.GetDisplayTitle()),
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			)

			// Remove button
			removeBtn := widget.NewButton("Remove", func() {
				q.RemoveAt(idx)
				ShowQueueView() // Refresh view
			})

			// Play now button
			playNowBtn := widget.NewButton("▶ Play Now", func() {
				q.RemoveAt(idx)
				playQueueItem(&qItem)
			})

			buttonContainer := container.NewGridWithColumns(2, playNowBtn, removeBtn)

			itemCard := container.NewVBox(
				titleLabel,
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

