# MovieStream Icon & Logo Design

## Application Icon (512x512)

### Design Elements
The application icon features a professional, cinema-themed design:

**Background**
- Deep purple/blue gradient (RGB: 20, 20, 30)
- Creates a premium, modern look
- Provides good contrast for icon elements

**Film Strip**
- Top and bottom horizontal strips (gray-purple: RGB 60, 60, 80)
- Authentic film perforation holes along the strips
- 8 evenly-spaced circular perforations on each strip
- Represents classic cinema and movie production

**Play Button**
- Large centered circle (purple: RGB 120, 100, 255)
- Inner glow effect for depth (RGB 140, 120, 255)
- White play triangle pointing right
- Symbolizes streaming and playback functionality

**Design Philosophy**
- Combines traditional cinema (film strip) with modern streaming (play button)
- Memorable and recognizable at all sizes
- Professional appearance suitable for desktop applications
- Color scheme matches the app's Geist dark theme

## UI Logo (180x50)

### Design Elements
A simplified, horizontal version for the application header:

**Structure**
- Horizontal film strip with perforations
- 5 rectangular perforation holes (top and bottom)
- Centered play triangle
- Transparent background

**Colors**
- Uses app's theme colors (GeistGray10 for main elements)
- Matches existing UI aesthetic
- Provides brand consistency

**Placement**
- Replaces "MovieStream" text in header
- 180x50px size optimized for header space
- Left-aligned with navigation buttons on right

## Technical Implementation

### Icon Generation
- Programmatically generated using Go's image library
- No external image files required
- Converted to PNG bytes for Fyne resource
- Set as window icon via `window.SetIcon()`

### Logo Integration
- Created via `CreateAppLogo()` function
- Returns `canvas.Image` widget
- Async loading (immediate display with theme colors)
- Used in both main UI and navigation views

### Code Structure
```go
// Create app icon for window
iconImg := gui.CreateAppIcon()
myWindow.SetIcon(fyne.NewStaticResource("icon.png", iconToBytes(iconImg)))

// Create logo for UI header
logo := CreateAppLogo(180, 50)
```

## Files
- `gui/logo.go` - Icon and logo generation functions
- `main.go` - Icon setup and helper functions
- `gui/app.go` - Logo integration in main UI
- `gui/navigation.go` - Logo integration in navigation

## Benefits

1. **Professional Branding**: Unique, recognizable visual identity
2. **No External Dependencies**: All graphics generated programmatically
3. **Theme Consistency**: Colors match the Geist dark theme
4. **Scalable**: Vector-like generation ensures quality at any size
5. **Cross-Platform**: Works on Windows, macOS, and Linux

## Future Enhancements

Potential improvements:
- Export icon as .ico file for Windows executables
- Create .icns file for macOS app bundles
- Add animated logo for splash screen
- Create icon variants for different platforms
- Add dark/light mode variations

## Design Credits

Icon Design: Cinema-inspired with film strip and play button motifs
Theme Integration: Matches Geist Design System monochrome palette

