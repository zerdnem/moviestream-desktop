package gui

import (
	"image/color"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Monochrome color palette (https://www.color-hex.com/color-palette/8528)
var (
	// Primary monochrome colors
	WarpAccent     = color.RGBA{R: 192, G: 192, B: 192, A: 255}   // #c0c0c0 Silver
	WarpBackground = color.RGBA{R: 0, G: 0, B: 0, A: 255}         // #000000 Black
	WarpForeground = color.RGBA{R: 255, G: 255, B: 255, A: 255}   // #ffffff White
	
	// Terminal colors - Normal (monochrome)
	WarpBlack      = color.RGBA{R: 0, G: 0, B: 0, A: 255}         // #000000 Black
	WarpBlue       = color.RGBA{R: 192, G: 192, B: 192, A: 255}   // #c0c0c0 Silver
	WarpCyan       = color.RGBA{R: 255, G: 255, B: 255, A: 255}   // #ffffff White
	WarpGreen      = color.RGBA{R: 192, G: 192, B: 192, A: 255}   // #c0c0c0 Silver
	WarpRed        = color.RGBA{R: 128, G: 128, B: 128, A: 255}   // #808080 Gray
	WarpMagenta    = color.RGBA{R: 192, G: 192, B: 192, A: 255}   // #c0c0c0 Silver
	WarpWhite      = color.RGBA{R: 255, G: 255, B: 255, A: 255}   // #ffffff White
	WarpYellow     = color.RGBA{R: 192, G: 192, B: 192, A: 255}   // #c0c0c0 Silver
	
	// Terminal colors - Bright (monochrome)
	WarpBrightBlack   = color.RGBA{R: 64, G: 64, B: 64, A: 255}     // #404040 Dark gray
	WarpBrightBlue    = color.RGBA{R: 192, G: 192, B: 192, A: 255}  // #c0c0c0 Silver
	WarpBrightCyan    = color.RGBA{R: 255, G: 255, B: 255, A: 255}  // #ffffff White
	WarpBrightGreen   = color.RGBA{R: 192, G: 192, B: 192, A: 255}  // #c0c0c0 Silver
	WarpBrightRed     = color.RGBA{R: 128, G: 128, B: 128, A: 255}  // #808080 Gray
	WarpBrightMagenta = color.RGBA{R: 192, G: 192, B: 192, A: 255}  // #c0c0c0 Silver
	WarpBrightWhite   = color.RGBA{R: 255, G: 255, B: 255, A: 255}  // #ffffff White
	WarpBrightYellow  = color.RGBA{R: 192, G: 192, B: 192, A: 255}  // #c0c0c0 Silver
	
	// UI Component colors (monochrome)
	WarpCard       = color.RGBA{R: 32, G: 32, B: 32, A: 255}      // Very dark gray
	WarpCardHover  = color.RGBA{R: 80, G: 80, B: 80, A: 255}      // #505050 Medium dark gray
	WarpBorder     = color.RGBA{R: 64, G: 64, B: 64, A: 255}      // #404040 Dark gray
	WarpDisabled   = color.RGBA{R: 128, G: 128, B: 128, A: 255}   // #808080 Gray
	WarpOverlay    = color.RGBA{R: 0, G: 0, B: 0, A: 180}         // Black overlay
	
	// Special button colors
	WarpButtonPrimary   = color.RGBA{R: 255, G: 255, B: 255, A: 255}  // White for primary buttons
	WarpButtonSecondary = color.RGBA{R: 64, G: 64, B: 64, A: 255}     // Dark gray for secondary
)

type MovieStreamTheme struct {
	isDark bool
}

var currentTheme *MovieStreamTheme

func init() {
	currentTheme = &MovieStreamTheme{isDark: true}
}

// GetCurrentTheme returns the current theme instance
func GetCurrentTheme() *MovieStreamTheme {
	return currentTheme
}

// IsDark returns whether the theme is in dark mode
func (t *MovieStreamTheme) IsDark() bool {
	return t.isDark
}

// SetDark sets the theme to dark or light mode
func (t *MovieStreamTheme) SetDark(dark bool) {
	t.isDark = dark
}

func (t *MovieStreamTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return WarpBackground
	case theme.ColorNameButton:
		return WarpBrightBlack // #404040 Dark gray buttons
	case theme.ColorNameDisabledButton:
		return WarpCard // Very dark for disabled state
	case theme.ColorNameForeground:
		return WarpForeground
	case theme.ColorNameHover:
		return WarpRed // #808080 Gray for hover
	case theme.ColorNameInputBackground:
		return WarpBrightBlack // Match button background
	case theme.ColorNamePrimary:
		return WarpAccent
	case theme.ColorNameFocus:
		return WarpAccent // Silver focus
	case theme.ColorNameSelection:
		return WarpAccent // Silver selection
	case theme.ColorNameShadow:
		return color.RGBA{R: 0, G: 0, B: 0, A: 150}
	case theme.ColorNameError:
		return WarpRed
	case theme.ColorNameSuccess:
		return WarpAccent // Use silver for success
	case theme.ColorNameWarning:
		return WarpRed // Use gray for warning
	case theme.ColorNameHeaderBackground:
		return WarpCard
	case theme.ColorNameInputBorder:
		return WarpBorder
	case theme.ColorNameSeparator:
		return WarpBorder
	case theme.ColorNameMenuBackground:
		return WarpCard
	case theme.ColorNameOverlayBackground:
		return WarpCard
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t *MovieStreamTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *MovieStreamTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *MovieStreamTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 4 // Compact padding
	case theme.SizeNameInlineIcon:
		return 16
	case theme.SizeNameScrollBar:
		return 6
	case theme.SizeNameScrollBarSmall:
		return 3
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameText:
		return 13 // Compact text
	case theme.SizeNameHeadingText:
		return 18 // Compact headers
	case theme.SizeNameSubHeadingText:
		return 15
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameInputBorder:
		return 1
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// Helper functions to get Warp Dark colors
func GetBackgroundColor() color.Color {
	return WarpBackground
}

func GetCardColor() color.Color {
	return WarpCard
}

func GetTextColor() color.Color {
	return WarpForeground
}

func GetSecondaryTextColor() color.Color {
	return WarpBrightBlack
}

func GetPrimaryColor() color.Color {
	return WarpAccent
}

func GetAccentColor() color.Color {
	return WarpBrightBlue
}

func GetSuccessColor() color.Color {
	return WarpGreen
}

func GetWarningColor() color.Color {
	return WarpYellow
}

func GetErrorColor() color.Color {
	return WarpRed
}

