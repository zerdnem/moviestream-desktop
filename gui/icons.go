package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Icon names for the application using Fyne's built-in icons
type IconName int

const (
	IconSearch IconName = iota
	IconSettings
	IconMovie
	IconTV
	IconPlay
	IconDownload
	IconBack
	IconStar
	IconCalendar
	IconCheck
	IconCancel
	IconInfo
	IconWarning
	IconFolder
	IconMenu
	IconHome
)

// GetIcon returns the appropriate Fyne icon resource
func GetIcon(name IconName) fyne.Resource {
	switch name {
	case IconSearch:
		return theme.SearchIcon()
	case IconSettings:
		return theme.SettingsIcon()
	case IconMovie:
		return theme.MediaVideoIcon()
	case IconTV:
		return theme.MediaVideoIcon()
	case IconPlay:
		return theme.MediaPlayIcon()
	case IconDownload:
		return theme.DownloadIcon()
	case IconBack:
		return theme.NavigateBackIcon()
	case IconStar:
		return theme.ConfirmIcon()
	case IconCalendar:
		return theme.InfoIcon()
	case IconCheck:
		return theme.ConfirmIcon()
	case IconCancel:
		return theme.CancelIcon()
	case IconInfo:
		return theme.InfoIcon()
	case IconWarning:
		return theme.WarningIcon()
	case IconFolder:
		return theme.FolderIcon()
	case IconMenu:
		return theme.MenuIcon()
	case IconHome:
		return theme.HomeIcon()
	default:
		return theme.InfoIcon()
	}
}

// ButtonWithIcon creates a button with an icon
func ButtonWithIcon(label string, icon IconName, tapped func()) *fyne.Container {
	btn := &iconButton{
		text: label,
		icon: GetIcon(icon),
		onTapped: tapped,
	}
	return btn.Create()
}

type iconButton struct {
	text     string
	icon     fyne.Resource
	onTapped func()
}

func (ib *iconButton) Create() *fyne.Container {
	// For now, we'll just use regular buttons with text
	// Fyne's Button widget already supports icons via widget.NewButtonWithIcon
	return nil
}

