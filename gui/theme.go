package gui

import (
	"image/color"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Geist Design System monochrome palette[](https://vercel.com/geist/colors)
var (
	// Backgrounds
	GeistBackground1 = color.RGBA{R: 0, G: 0, B: 0, A: 255}      // #000000 - Default element background
	GeistBackground2 = color.RGBA{R: 10, G: 10, B: 10, A: 255}   // #0a0a0a - Secondary background
	
	// Gray scale (Dark mode values)
	GeistGray1  = color.RGBA{R: 17, G: 17, B: 17, A: 255}   // Component default background
	GeistGray2  = color.RGBA{R: 23, G: 23, B: 23, A: 255}   // Component hover background
	GeistGray3  = color.RGBA{R: 31, G: 31, B: 31, A: 255}   // Component active background
	GeistGray4  = color.RGBA{R: 46, G: 46, B: 46, A: 255}   // Default border
	GeistGray5  = color.RGBA{R: 62, G: 62, B: 62, A: 255}   // Hover border
	GeistGray6  = color.RGBA{R: 77, G: 77, B: 77, A: 255}   // Active border
	GeistGray7  = color.RGBA{R: 99, G: 99, B: 99, A: 255}   // High contrast background
	GeistGray8  = color.RGBA{R: 117, G: 117, B: 117, A: 255}  // Hover high contrast background
	GeistGray9  = color.RGBA{R: 150, G: 150, B: 150, A: 255}  // Secondary text and icons
	GeistGray10 = color.RGBA{R: 238, G: 238, B: 238, A: 255}  // Primary text and icons
	
	// Alpha variants for overlays
	GeistGrayAlpha1 = color.RGBA{R: 255, G: 255, B: 255, A: 6}
	GeistGrayAlpha2 = color.RGBA{R: 255, G: 255, B: 255, A: 13}
	GeistGrayAlpha3 = color.RGBA{R: 255, G: 255, B: 255, A: 20}
	GeistGrayAlpha4 = color.RGBA{R: 255, G: 255, B: 255, A: 31}
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
	// Note: Theme preference is managed by Fyne's built-in theme system
}

func (t *MovieStreamTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return GeistBackground1
	case theme.ColorNameButton:
		return GeistGray1 // Default component background for buttons
	case theme.ColorNameDisabledButton:
		return GeistBackground1 // Darker background for disabled buttons
	case theme.ColorNameForeground:
		return GeistGray10 // Primary text and icons
	case theme.ColorNameDisabled:
		return GeistGray9 // Secondary text for disabled state
	case theme.ColorNameHover:
		return GeistGray2 // Hover background for buttons and interactive elements
	case theme.ColorNamePressed:
		return GeistGray3 // Active (pressed) background for buttons
	case theme.ColorNameInputBackground:
		return GeistGray1 // Component default background
	case theme.ColorNamePrimary:
		return GeistGray10 // Primary text for selected items
	case theme.ColorNameFocus:
		return GeistGray10 // High contrast white for focus
	case theme.ColorNameSelection:
		return GeistGray3 // Active background for selected items
	case theme.ColorNameShadow:
		return color.RGBA{R: 0, G: 0, B: 0, A: 150}
	case theme.ColorNameError:
		return GeistGray9 // Secondary text (monochrome)
	case theme.ColorNameSuccess:
		return GeistGray10 // Primary text (monochrome)
	case theme.ColorNameWarning:
		return GeistGray9 // Secondary text (monochrome)
	case theme.ColorNameHeaderBackground:
		return GeistBackground2
	case theme.ColorNameInputBorder:
		return GeistGray4 // Default border
	case theme.ColorNameSeparator:
		return GeistGray4 // Default border
	case theme.ColorNameMenuBackground:
		return GeistGray1
	case theme.ColorNameOverlayBackground:
		return GeistGray1
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t *MovieStreamTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *MovieStreamTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	// Use default theme icons - they will automatically use our monochrome foreground color
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

// Helper functions to get Geist monochrome colors
func GetBackgroundColor() color.Color {
	return GeistBackground1
}

func GetCardColor() color.Color {
	return GeistGray1
}

func GetTextColor() color.Color {
	return GeistGray10
}

func GetSecondaryTextColor() color.Color {
	return GeistGray9
}

func GetPrimaryColor() color.Color {
	return GeistGray4 // Updated to match theme primary
}

func GetAccentColor() color.Color {
	return GeistGray4
}

func GetSuccessColor() color.Color {
	return GeistGray10
}

func GetWarningColor() color.Color {
	return GeistGray9
}

func GetErrorColor() color.Color {
	return GeistGray9
}