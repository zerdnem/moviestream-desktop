# Subtitle Download Feature

## Overview
Added a comprehensive subtitle download feature that automatically detects when no subtitles are available in the stream and provides users with the option to download subtitles from OpenSubtitles.org.

## Features

### 1. Automatic Detection
- The application now automatically detects when a movie or TV show has no embedded subtitles
- When no subtitles are found, a modal dialog appears instead of just playing without subtitles

### 2. OpenSubtitles Integration
- Integrated with OpenSubtitles API to search for subtitles
- Searches by title, season, and episode (for TV shows)
- Respects user's subtitle language preference from settings
- Falls back to English subtitles if preferred language is not available

### 3. User-Friendly Dialog
The subtitle download dialog provides:
- **Automatic Search**: Searches for subtitles as soon as the dialog opens
- **Multiple Results**: Displays all available subtitle matches
- **Language Info**: Shows language name and code for each subtitle
- **File Names**: Displays the subtitle file name for transparency
- **Selection Interface**: Radio buttons to select preferred subtitle
- **Multiple Options**:
  - Download & Play: Download selected subtitle and play with it
  - Play Without Subtitles: Skip subtitle download and play anyway
  - Cancel: Cancel playback entirely

### 4. Smart Subtitle Handling
- Downloads subtitles to temporary directory
- Automatically loads subtitle into the video player
- Supports all video players (MPV, VLC, MPC-HC, PotPlayer)
- Cleans up temporary files after playback

## Implementation Details

### New Packages

#### `subtitles/opensubtitles.go`
- Implements OpenSubtitles REST API client
- Handles subtitle search by title or IMDB ID
- Supports season/episode parameters for TV shows
- Downloads and saves subtitle files
- Maps language codes between internal format and OpenSubtitles format

#### `subtitles/manager.go`
- High-level subtitle management interface
- Coordinates subtitle search across sources
- Handles language preferences
- Provides simple API for GUI integration

#### `gui/subtitlesdialog.go`
- Creates and manages the subtitle download modal dialog
- Implements asynchronous subtitle search
- Provides user selection interface
- Handles download and playback coordination

### Modified Files

#### `gui/app.go` (watchMovie function)
- Added check for empty subtitle list
- Shows subtitle download dialog when no subtitles found
- Maintains queue auto-play functionality

#### `gui/tvdetails.go` (watchEpisodeInternal function)
- Added check for empty subtitle list
- Shows subtitle download dialog when no subtitles found
- Maintains auto-next and queue functionality

## User Experience

### Before
When no subtitles were found:
```
⚠ No subtitles available for this content
[Video plays without subtitles]
```

### After
When no subtitles are found:
1. A modal dialog appears: "No subtitles found in stream"
2. Dialog shows: "Would you like to search for subtitles from OpenSubtitles?"
3. Application automatically searches for subtitles
4. User sees available subtitle options with languages
5. User can:
   - Select a subtitle and download it
   - Play without subtitles
   - Cancel playback

### With Subtitles Downloaded
```
▶ Playing: Movie Title
   ✓ Subtitle loaded: subtitle_file.srt (English)
```

## Technical Notes

### API Usage
- Uses OpenSubtitles REST API (no authentication required for basic searches)
- Respects API rate limits
- User-Agent: "MovieStream v1.0"

### File Handling
- Subtitles are downloaded to system temp directory
- Files are prefixed with "moviestream_" for easy identification
- Temporary files are cleaned up by video players after use

### Language Support
Supports subtitle search in multiple languages:
- English (en/eng)
- Spanish (es/spa)
- French (fr/fre)
- German (de/ger)
- Italian (it/ita)
- Portuguese (pt/por)
- Japanese (ja/jpn)
- Korean (ko/kor)
- Chinese (zh/chi)
- Arabic (ar/ara)
- Russian (ru/rus)
- Hindi (hi/hin)

### Player Compatibility
The feature works with all supported video players:
- **MPV**: Loads subtitle via --sub-file parameter
- **VLC**: Downloads and loads via --sub-file parameter
- **MPC-HC**: Downloads and loads via /sub parameter
- **PotPlayer**: Downloads and loads via /sub parameter

## Future Enhancements

Potential improvements for future versions:
1. Cache subtitle search results to avoid repeated API calls
2. Add manual subtitle file upload option
3. Support additional subtitle sources (subscene.com, etc.)
4. Subtitle preview/rating display
5. Automatic subtitle sync adjustment
6. Save downloaded subtitles permanently (optional)
7. OpenSubtitles API authentication for enhanced features

## Dependencies

No new external dependencies required. Uses only Go standard library packages:
- `encoding/json` - API response parsing
- `net/http` - HTTP requests
- `net/url` - URL parameter encoding
- `io` - File operations
- `os` - File system operations
- `path/filepath` - Path handling
- `strings` - String manipulation

## Testing

To test the feature:
1. Build and run the application
2. Search for a movie or TV show
3. Select one known to have no embedded subtitles
4. Click "Watch"
5. The subtitle download dialog should appear
6. Select a subtitle and click "Download & Play"
7. Video should play with downloaded subtitle

Alternatively:
- Click "Play Without Subtitles" to skip subtitle download
- Click "Cancel" to abort playback

