package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// TappableCard is a container that can be tapped without visible button styling
type TappableCard struct {
	widget.BaseWidget
	content fyne.CanvasObject
	onTap   func()
}

// NewTappableCard creates a new tappable card with the given content and tap callback
func NewTappableCard(content fyne.CanvasObject, onTap func()) *TappableCard {
	tc := &TappableCard{
		content: content,
		onTap:   onTap,
	}
	tc.ExtendBaseWidget(tc)
	return tc
}

// CreateRenderer creates the renderer for the tappable card
func (tc *TappableCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(tc.content)
}

// Tapped handles tap events
func (tc *TappableCard) Tapped(_ *fyne.PointEvent) {
	if tc.onTap != nil {
		tc.onTap()
	}
}

// TappedSecondary handles secondary tap events (right-click)
func (tc *TappableCard) TappedSecondary(_ *fyne.PointEvent) {
	// Not used
}

// Cursor returns the cursor type for the tappable card
func (tc *TappableCard) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

