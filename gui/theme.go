package gui

import (
	"image/color"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Color palette from https://www.color-hex.com/color-palette/1050053
var (
	// Primary colors
	PrimaryBlue   = color.RGBA{R: 0, G: 175, B: 240, A: 255}     // #00aff0
	AccentBlue    = color.RGBA{R: 1, G: 140, B: 241, A: 255}     // #018cf1
	DarkBg        = color.RGBA{R: 39, G: 39, B: 43, A: 255}      // #27272b
	White         = color.RGBA{R: 255, G: 255, B: 255, A: 255}   // #ffffff
	Gray          = color.RGBA{R: 138, G: 150, B: 163, A: 255}   // #8a96a3
	
	// Additional colors for better UI
	DarkCard      = color.RGBA{R: 50, G: 50, B: 55, A: 255}
	LightBg       = color.RGBA{R: 248, G: 249, B: 250, A: 255}
	LightCard     = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	LightText     = color.RGBA{R: 33, G: 37, B: 41, A: 255}
	DarkHover     = color.RGBA{R: 60, G: 60, B: 65, A: 255}
	LightHover    = color.RGBA{R: 235, G: 235, B: 240, A: 255}
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
		if t.isDark {
			return DarkBg
		}
		return LightBg
	case theme.ColorNameButton:
		return PrimaryBlue
	case theme.ColorNameDisabledButton:
		return Gray
	case theme.ColorNameForeground:
		if t.isDark {
			return White
		}
		return LightText
	case theme.ColorNameHover:
		if t.isDark {
			return DarkHover
		}
		return LightHover
	case theme.ColorNameInputBackground:
		if t.isDark {
			return DarkCard
		}
		return LightCard
	case theme.ColorNamePrimary:
		return PrimaryBlue
	case theme.ColorNameFocus:
		return AccentBlue
	case theme.ColorNameSelection:
		return AccentBlue
	case theme.ColorNameShadow:
		if t.isDark {
			return color.RGBA{R: 0, G: 0, B: 0, A: 100}
		}
		return color.RGBA{R: 0, G: 0, B: 0, A: 50}
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
		return 6
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 10
	case theme.SizeNameScrollBarSmall:
		return 6
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// Helper functions to get colors based on current theme
func GetBackgroundColor() color.Color {
	if currentTheme.isDark {
		return DarkBg
	}
	return LightBg
}

func GetCardColor() color.Color {
	if currentTheme.isDark {
		return DarkCard
	}
	return LightCard
}

func GetTextColor() color.Color {
	if currentTheme.isDark {
		return White
	}
	return LightText
}

func GetSecondaryTextColor() color.Color {
	return Gray
}

func GetPrimaryColor() color.Color {
	return PrimaryBlue
}

func GetAccentColor() color.Color {
	return AccentBlue
}

