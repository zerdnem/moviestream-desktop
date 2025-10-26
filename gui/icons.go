package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Icon system using Fyne theme icons and clean symbols
const (
	// Navigation
	IconSearch      = "icon-search"
	IconHome        = "icon-home"
	IconBack        = "icon-back"
	IconNext        = "icon-next"
	IconMenu        = "icon-menu"
	
	// Media
	IconPlay        = "icon-play"
	IconPause       = "icon-pause"
	IconStop        = "icon-stop"
	IconMovie       = "icon-movie"
	IconTV          = "icon-tv"
	IconVideo       = "icon-video"
	
	// Actions
	IconDownload    = "icon-download"
	IconUpload      = "icon-upload"
	IconAdd         = "icon-add"
	IconRemove      = "icon-remove"
	IconDelete      = "icon-delete"
	IconEdit        = "icon-edit"
	IconSave        = "icon-save"
	
	// Status
	IconCheck       = "icon-check"
	IconCross       = "icon-cross"
	IconInfo        = "icon-info"
	IconWarning     = "icon-warning"
	IconError       = "icon-error"
	IconStar        = "icon-star"
	
	// UI
	IconSettings    = "icon-settings"
	IconUser        = "icon-user"
	IconFolder      = "icon-folder"
	IconFile        = "icon-file"
	IconClock       = "icon-clock"
	IconCalendar    = "icon-calendar"
	IconQueue       = "icon-queue"
	IconHistory     = "icon-history"
)

// CreateIconButton creates a button with Fyne icon
func CreateIconButton(label string, iconName string, tapped func()) *widget.Button {
	// Map icon names to Fyne theme icons
	var icon fyne.Resource
	switch iconName {
	case IconSearch:
		icon = theme.SearchIcon()
	case IconSettings:
		icon = theme.SettingsIcon()
	case IconPlay:
		icon = theme.MediaPlayIcon()
	case IconDownload:
		icon = theme.DownloadIcon()
	case IconBack:
		icon = theme.NavigateBackIcon()
	case IconHome:
		icon = theme.HomeIcon()
	case IconMenu:
		icon = theme.MenuIcon()
	case IconDelete:
		icon = theme.DeleteIcon()
	case IconAdd:
		icon = theme.ContentAddIcon()
	case IconRemove:
		icon = theme.ContentRemoveIcon()
	case IconCheck:
		icon = theme.ConfirmIcon()
	case IconCross:
		icon = theme.CancelIcon()
	case IconInfo:
		icon = theme.InfoIcon()
	case IconWarning:
		icon = theme.WarningIcon()
	case IconFolder:
		icon = theme.FolderIcon()
	case IconFile:
		icon = theme.FileIcon()
	case IconQueue:
		icon = theme.ListIcon()
	case IconHistory:
		icon = theme.HistoryIcon()
	}
	
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
	
	// For high importance buttons in monochrome, use white background
	if importance == widget.HighImportance {
		// The theme will handle this
	}
	
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
	border := canvas.NewRectangle(WarpBorder)
	
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

// GetFyneIcon returns Fyne's built-in theme icons as fallback
func GetFyneIcon(name string) fyne.Resource {
	switch name {
	case "search":
		return theme.SearchIcon()
	case "settings":
		return theme.SettingsIcon()
	case "play":
		return theme.MediaPlayIcon()
	case "download":
		return theme.DownloadIcon()
	case "back":
		return theme.NavigateBackIcon()
	case "home":
		return theme.HomeIcon()
	case "menu":
		return theme.MenuIcon()
	case "delete":
		return theme.DeleteIcon()
	default:
		return theme.InfoIcon()
	}
}

