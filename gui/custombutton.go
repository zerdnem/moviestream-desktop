package gui

import (
	"image/color"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreatePrimaryButton creates a white button with black text for primary actions
func CreatePrimaryButton(label string, icon fyne.Resource, tapped func()) *fyne.Container {
	// White background
	bg := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	
	// Black text
	var btn *widget.Button
	if icon != nil {
		btn = widget.NewButtonWithIcon(label, icon, tapped)
	} else {
		btn = widget.NewButton(label, tapped)
	}
	
	// Override button appearance
	btn.Importance = widget.HighImportance
	
	content := container.NewStack(bg, btn)
	return content
}

