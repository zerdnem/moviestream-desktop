package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	
	"github.com/gosthome/icons"
	"github.com/gosthome/icons/fynico"
	_ "github.com/gosthome/icons/fynico/templarian/mdi" // Import MDI icons
)

// Icon system using Material Design Icons
const (
	// Navigation
	IconSearch      = "mdi:magnify"
	IconHome        = "mdi:home"
	IconBack        = "mdi:arrow-left"
	IconNext        = "mdi:arrow-right"
	IconMenu        = "mdi:menu"
	
	// Media
	IconPlay        = "mdi:play"
	IconPause       = "mdi:pause"
	IconStop        = "mdi:stop"
	IconMovie       = "mdi:movie"
	IconTV          = "mdi:television"
	IconVideo       = "mdi:video"
	
	// Actions
	IconDownload    = "mdi:download"
	IconUpload      = "mdi:upload"
	IconAdd         = "mdi:plus"
	IconRemove      = "mdi:minus"
	IconDelete      = "mdi:delete"
	IconEdit        = "mdi:pencil"
	IconSave        = "mdi:content-save"
	
	// Status
	IconCheck       = "mdi:check"
	IconCross       = "mdi:close"
	IconInfo        = "mdi:information"
	IconWarning     = "mdi:alert"
	IconError       = "mdi:alert-circle"
	IconStar        = "mdi:star"
	
	// UI
	IconSettings    = "mdi:cog"
	IconUser        = "mdi:account"
	IconFolder      = "mdi:folder"
	IconFile        = "mdi:file"
	IconClock       = "mdi:clock"
	IconCalendar    = "mdi:calendar"
	IconQueue       = "mdi:playlist-play"
	IconHistory     = "mdi:history"
)

// GetIconResource returns a themed icon resource from the icon name
func GetIconResource(iconName string) fyne.Resource {
	p, err := icons.Parse(iconName)
	if err != nil {
		// Fallback to theme icon
		return theme.InfoIcon()
	}
	
	r := fynico.Collections.Lookup(p.Collection, p.Icon)
	if r == nil {
		// Fallback to theme icon
		return theme.InfoIcon()
	}
	
	// Return the icon with theme foreground color
	return theme.NewThemedResource(r)
}

// CreateIconButton creates a button with Material Design Icon
func CreateIconButton(label string, iconName string, tapped func()) *widget.Button {
	icon := GetIconResource(iconName)
	
	if icon != nil && label != "" {
		return widget.NewButtonWithIcon(label, icon, tapped)
	} else if icon != nil {
		return widget.NewButtonWithIcon("", icon, tapped)
	}
	return widget.NewButton(label, tapped)
}

// CreateIconButtonWithImportance creates a button with icon and importance styling
func CreateIconButtonWithImportance(label string, iconName string, importance widget.Importance, tapped func()) *widget.Button {
	btn := CreateIconButton(label, iconName, tapped)
	btn.Importance = importance
	return btn
}

// CreateCompactButton creates a compact button for the modern UI
func CreateCompactButton(label string, iconName string, tapped func()) *widget.Button {
	return CreateIconButton(label, iconName, tapped)
}

// CreateCard creates a modern card container with background
func CreateCard(content ...fyne.CanvasObject) *fyne.Container {
	bg := canvas.NewRectangle(GetCardColor())
	
	cardContent := container.NewVBox(content...)
	
	return container.NewStack(
		bg,
		container.NewPadded(cardContent),
	)
}

// CreateMovieCard creates a streaming-service style card for movies/shows
func CreateMovieCard(content ...fyne.CanvasObject) *fyne.Container {
	bg := canvas.NewRectangle(GetCardColor())
	
	// Compact padding
	cardContent := container.NewVBox(content...)
	padded := container.NewPadded(cardContent)
	
	return container.NewStack(bg, padded)
}

// CreateCardWithBorder creates a card with a border
func CreateCardWithBorder(content ...fyne.CanvasObject) *fyne.Container {
	bg := canvas.NewRectangle(GetCardColor())
	border := canvas.NewRectangle(GeistGray4) // Default border color
	
	cardContent := container.NewVBox(content...)
	padded := container.NewPadded(cardContent)
	
	return container.NewStack(
		border,
		container.NewPadded(
			container.NewStack(bg, padded),
		),
	)
}

// CreateHeader creates a styled header text
func CreateHeader(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.TextStyle = fyne.TextStyle{Bold: true}
	return label
}

// CreateAccentText creates text with accent color
func CreateAccentText(text string, size float32) *canvas.Text {
	t := canvas.NewText(text, GetAccentColor())
	t.TextSize = size
	t.TextStyle = fyne.TextStyle{Bold: true}
	return t
}

// CreateTitle creates a large title with accent color
func CreateTitle(text string) *canvas.Text {
	return CreateAccentText(text, 18)
}

// CreateLargeTitle creates an extra large title for hero sections
func CreateLargeTitle(text string) *canvas.Text {
	return CreateAccentText(text, 22)
}

// CreateSubtitle creates a subtitle text
func CreateSubtitle(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.TextStyle = fyne.TextStyle{Bold: false}
	return label
}

// CreateIconLabel creates a label with an icon
func CreateIconLabel(iconName string, text string) *fyne.Container {
	icon := canvas.NewImageFromResource(GetIconResource(iconName))
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSize(20, 20))
	
	label := widget.NewLabel(text)
	
	return container.NewHBox(icon, label)
}

// CreateTitleWithIcon creates a title with an icon
func CreateTitleWithIcon(iconName string, text string) *fyne.Container {
	icon := canvas.NewImageFromResource(GetIconResource(iconName))
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSize(24, 24))
	
	title := CreateTitle(text)
	
	return container.NewHBox(icon, title)
}

