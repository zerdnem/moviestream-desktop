package gui

import (
	"fmt"
	"moviestream-gui/history"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowHistoryView displays the watch history with a modern design
func ShowHistoryView() {
	h := history.Get()
	items := h.GetAll()

	// Navigation buttons
	backBtn := CreateIconButton("", IconBack, func() {
		GoBackToSearch()
	})
	backBtn.Importance = widget.LowImportance

	// Hero section with title and stats
	titleText := CreateLargeTitle(fmt.Sprintf("Watch History"))
	statsText := widget.NewLabelWithStyle(
		fmt.Sprintf("%d items watched", len(items)),
		fyne.TextAlignLeading,
		fyne.TextStyle{},
	)
	
	heroSection := container.NewVBox(
		container.NewHBox(backBtn, CreateIconLabel(IconHistory, "History")),
		widget.NewSeparator(),
		container.NewCenter(titleText),
		container.NewCenter(statsText),
	)

	// Continue watching section (prominent placement)
	var continueSection *fyne.Container
	if lastShow, ok := h.GetLastWatchedShow(); ok {
		continueBtn := CreateIconButtonWithImportance("Continue Watching", IconPlay, widget.HighImportance, func() {
			watchEpisodeWithAutoNext(lastShow.TMDBID, lastShow.Season, lastShow.Episode, lastShow.Title, lastShow.EpisodeName, 0)
		})
		
		continueLabel := widget.NewLabelWithStyle(
			"Pick up where you left off",
			fyne.TextAlignCenter,
			fyne.TextStyle{Italic: true},
		)
		
		continueSection = container.NewVBox(
			widget.NewSeparator(),
			container.NewCenter(continueLabel),
			container.NewCenter(container.NewHBox(continueBtn)),
			widget.NewSeparator(),
		)
	}

	// History items with modern card design
	var itemWidgets []fyne.CanvasObject

	if len(items) == 0 {
		// Empty state with icon
		emptyIcon := canvas.NewImageFromResource(GetIconResource(IconHistory))
		emptyIcon.FillMode = canvas.ImageFillContain
		emptyIcon.SetMinSize(fyne.NewSize(64, 64))
		
		emptyMsg := container.NewCenter(
			container.NewVBox(
				widget.NewLabel(""),
				widget.NewLabel(""),
				container.NewCenter(emptyIcon),
				widget.NewLabel(""),
				CreateHeader("No Watch History"),
				widget.NewLabel("Your watched movies and shows will appear here"),
				widget.NewLabel(""),
			),
		)
		itemWidgets = append(itemWidgets, emptyMsg)
	} else {
		// Section header
		recentHeader := widget.NewLabelWithStyle(
			"Recently Watched",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		)
		itemWidgets = append(itemWidgets, widget.NewSeparator(), recentHeader)

		// Create modern cards for each history item
		for i, item := range items {
			idx := i
			hItem := item
			
			// Type icon
			var typeIcon fyne.CanvasObject
			if hItem.Type == "movie" {
				icon := canvas.NewImageFromResource(GetIconResource(IconMovie))
				icon.FillMode = canvas.ImageFillContain
				icon.SetMinSize(fyne.NewSize(20, 20))
				typeIcon = icon
			} else {
				icon := canvas.NewImageFromResource(GetIconResource(IconTV))
				icon.FillMode = canvas.ImageFillContain
				icon.SetMinSize(fyne.NewSize(20, 20))
				typeIcon = icon
			}

			// Title with icon
			titleLabel := widget.NewLabelWithStyle(
				hItem.GetDisplayTitle(),
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			)
			titleRow := container.NewHBox(typeIcon, titleLabel)

			// Time watched with clock icon
			timeIcon := canvas.NewImageFromResource(GetIconResource(IconClock))
			timeIcon.FillMode = canvas.ImageFillContain
			timeIcon.SetMinSize(fyne.NewSize(16, 16))
			
			timeLabel := widget.NewLabelWithStyle(
				hItem.GetDisplayTime(),
				fyne.TextAlignLeading,
				fyne.TextStyle{},
			)
			timeRow := container.NewHBox(timeIcon, timeLabel)

			// Action buttons with modern styling
			watchAgainBtn := CreateIconButtonWithImportance("Watch Again", IconPlay, widget.HighImportance, func() {
				if hItem.Type == "movie" {
					watchMovieByID(hItem.TMDBID, hItem.Title)
				} else {
					watchEpisodeWithAutoNext(hItem.TMDBID, hItem.Season, hItem.Episode, hItem.Title, hItem.EpisodeName, 0)
				}
			})

			removeBtn := CreateIconButton("", IconDelete, func() {
				dialog.ShowConfirm("Remove from History", 
					fmt.Sprintf("Remove '%s' from your watch history?", hItem.GetDisplayTitle()),
					func(confirmed bool) {
						if confirmed {
							h.Remove(idx)
							ShowHistoryView() // Refresh view
						}
					}, currentWindow)
			})
			removeBtn.Importance = widget.DangerImportance

			buttonContainer := container.NewHBox(watchAgainBtn, removeBtn)

			// Create modern card with background
			cardContent := container.NewVBox(
				titleRow,
				timeRow,
				widget.NewSeparator(),
				buttonContainer,
			)

			// Add card with background
			card := CreateCard(cardContent)
			itemWidgets = append(itemWidgets, card)
		}

		// Clear all button at bottom
		itemWidgets = append(itemWidgets, widget.NewSeparator())
		clearBtn := CreateIconButton("Clear All History", IconDelete, func() {
			dialog.ShowConfirm("Clear History", 
				"Are you sure you want to clear your entire watch history? This cannot be undone.",
				func(confirmed bool) {
					if confirmed {
						h.Clear()
						ShowHistoryView() // Refresh view
						dialog.ShowInformation("History Cleared", "All watch history has been removed.", currentWindow)
					}
				}, currentWindow)
		})
		clearBtn.Importance = widget.DangerImportance
		itemWidgets = append(itemWidgets, container.NewCenter(container.NewHBox(clearBtn)))
	}

	// Assemble content
	contentItems := []fyne.CanvasObject{heroSection}
	if continueSection != nil {
		contentItems = append(contentItems, continueSection)
	}
	contentItems = append(contentItems, itemWidgets...)

	content := container.NewVBox(contentItems...)
	scrollContent := container.NewVScroll(content)

	// Main layout
	mainLayout := container.NewPadded(scrollContent)

	currentWindow.SetContent(mainLayout)
}

