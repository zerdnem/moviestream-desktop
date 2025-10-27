# Audio Track Feature Documentation

## Overview

The MovieStream application now supports **automatic audio track detection** and selection from the streaming service, as well as adding external audio tracks. This feature is particularly useful when watching videos in one language but wanting to hear audio in another language (e.g., watching a Korean movie with English audio).

## Features

### 1. **Automatic Audio Track Detection**
- **Automatically fetches** available audio tracks from the streaming API
- **Displays all available options** (e.g., English, Spanish, French, etc.)
- **One-click selection** to play with desired audio track
- Works seamlessly with the existing streaming infrastructure

### 2. **External Audio Track Integration**
- Add external audio tracks from URLs or local files (if API tracks aren't available)
- Support for multiple audio formats: AAC, MP3, OGG, Opus
- Seamless integration with video players (MPV, VLC, PotPlayer)

### 3. **Multi-Player Support**
- **MPV**: Full support with `--audio-file` parameter
- **VLC**: Support via `--input-slave` parameter  
- **PotPlayer**: Support via `/add` parameter
- **MPC-HC**: Limited command-line support (manual loading recommended)

### 4. **User-Friendly Interface**
- Dedicated "Watch with Audio Tracks" button on movie details
- **Automatic display of available audio tracks** from streaming service
- Dialog interface for selecting API tracks or adding external ones
- Visual list of tracks with easy selection

## Usage

### Watching a Movie with Different Audio Tracks

#### Method 1: Using API Audio Tracks (Recommended)

1. **Search for a movie** in the main interface
2. **Click on the movie** to view details
3. **Click "Watch with Audio Tracks"** button
4. In the Audio Track Dialog:
   - **View available audio tracks** (e.g., "English", "Spanish", "Original")
   - **Select your preferred audio track** from the list
   - Click "Play with Selected Track"
5. The video will start playing with the selected audio!

#### Method 2: Using External Audio Tracks (Fallback)

If no API tracks are available, you can still add external tracks:

1. Follow steps 1-3 above
2. In the Audio Track Dialog:
   - Scroll down to "Or add external audio tracks"
   - Enter a URL or local file path to an audio track
   - Select the language (optional)
   - Click "Add Track"
   - Repeat for multiple audio tracks if needed
3. **Click "Play with Audio Tracks"** to start playback

### Audio Track Sources

You can add audio tracks from:
- **Direct URLs**: `https://example.com/audio-english.aac`
- **Local Files**: `C:\path\to\audio\english.mp3`

### Supported Formats

- **AAC** (.aac) - Recommended for best quality
- **MP3** (.mp3) - Widely compatible
- **OGG** (.ogg) - Open format
- **Opus** (.opus) - Modern, efficient codec

## Player Controls

### MPV
- Press **#** to cycle through audio tracks
- Press **m** to mute/unmute
- OSD shows active audio track

### VLC
- Press **B** to cycle through audio tracks
- Audio > Audio Track menu for manual selection
- Shows all loaded tracks in menu

### PotPlayer
- Right-click > Sound > Audio Stream
- Select from loaded audio tracks
- Supports multiple simultaneous tracks

### MPC-HC
- Navigate > Audio menu
- Limited command-line support
- Manual loading recommended via GUI

## Technical Implementation

### Architecture

```
api/
├── stream.go        # Fetches audio tracks from streaming API
└── browser.go       # Browser automation for stream extraction

audiotracks/
├── types.go         # Audio track data structures
└── manager.go       # Audio track management and download

player/
├── launcher.go      # Updated to support audio tracks
└── mpv.go          # Audio track playback functions

gui/
├── audiotrackdialog.go  # UI for audio track selection
└── app.go              # Integration with main app
```

### How API Audio Tracks Work

The streaming API provides multiple audio tracks for content:

1. **Fetch Phase**: When getting a stream URL, `GetStreamURL()` also fetches available audio tracks
2. **Audio Track Data**: Each track includes:
   - `Name`: Track name (e.g., "English", "Spanish")
   - `Description`: Additional info (e.g., "English Dubbed")
   - `Data`: Encoded track identifier for the API
3. **Selection Phase**: User selects a track in the UI
4. **Stream Generation**: `GetStreamURLWithAudioTrack()` generates a new stream URL with the selected audio
5. **Playback**: Video plays with the chosen audio track

### Key Components

#### 1. **API Stream Info** (`api/stream.go`)
```go
type StreamInfo struct {
    StreamURL    string
    SubtitleURLs []SubtitleTrack
    AudioTracks  []AudioTrack  // Available audio tracks from API
}

type AudioTrack struct {
    Name        string  // e.g., "English", "Spanish"
    Description string  // e.g., "English Dubbed"
    Data        string  // API identifier for the track
}

// Get stream with all available audio tracks
func GetStreamURL(tmdbID int, contentType string, season, episode int) (*StreamInfo, error)

// Get stream URL for a specific audio track
func GetStreamURLWithAudioTrack(tmdbID int, contentType string, season, episode int, audioTrackData string) (string, error)
```

#### 2. **Audio Track Manager** (`audiotracks/manager.go`)
```go
type Manager struct {
    // Handles external audio track download and validation
}

// Download audio track from URL
func (m *Manager) DownloadAudioTrack(url string, filename string) (string, error)

// Validate audio file exists and is readable
func (m *Manager) ValidateAudioFile(path string) error
```

#### 3. **Player Integration** (`player/launcher.go`)
All player launch functions now accept `audioTrackURLs []string`:
```go
func launchMPV(exePath, streamURL, title string, 
               subtitleURLs []string, 
               audioTrackURLs []string, 
               onEnd OnPlaybackEndCallback) error
```

#### 4. **UI Dialog** (`gui/audiotrackdialog.go`)
```go
// Show dialog for selecting audio tracks (API or external)
func ShowAudioTrackDialog(
    title string, 
    tmdbID int,
    season, episode int, 
    streamURL string,
    subtitleURLs []string,
    availableAudioTracks []api.AudioTrack,  // From API
    onEnd player.OnPlaybackEndCallback)
```

## Examples

### Example 1: Selecting English Audio Track (API)

1. Search for "Parasite" (Korean movie)
2. Click movie card to view details
3. Click "Watch with Audio Tracks"
4. **See available audio tracks**:
   - Original - Korean
   - English - English Dubbed
   - Spanish - Spanish Dubbed
5. **Select "English - English Dubbed"**
6. Click "Play with Selected Track"
7. Video starts playing with English audio!

### Example 2: Switching Between Multiple API Audio Tracks

The streaming service often provides multiple audio tracks:
1. Click "Watch with Audio Tracks"
2. **See all available tracks**:
   - Original - Korean
   - English - English Dubbed
   - Spanish - Spanish Dubbed
   - French - French Dubbed
3. **Select any track** and play
4. Can restart and select a different track anytime

### Example 3: Local Audio File

If you have a local audio file:
1. Click "Watch with Audio Tracks"
2. Enter local path: `C:\Downloads\movie-audio-english.mp3`
3. Click "Add Track"
4. Play with the external audio

## Temporary File Management

Audio tracks downloaded from URLs are stored in:
```
%TEMP%\moviestream_audio\
```

Files are automatically cleaned up when:
- Player window closes
- Application exits
- Manual cleanup via `audiotracks.CleanupTempAudioFiles()`

## Limitations

1. **MPC-HC**: Limited command-line support for external audio tracks. Manual loading through the GUI is recommended.

2. **File Validation**: The app validates local files but cannot verify URLs until download. Invalid URLs will show errors during playback.

3. **Format Support**: Player-dependent. Most players support AAC, MP3, OGG, and Opus. Some exotic formats may not work.

4. **Sync Issues**: External audio tracks may have synchronization issues if:
   - Audio duration doesn't match video duration
   - Audio has different frame rate/timing
   - Use professional audio editing tools to sync before adding

## Future Enhancements

Potential future features:
- Automatic audio track search from repositories
- Audio track synchronization adjustment
- Built-in audio format conversion
- Audio track metadata display (codec, bitrate, channels)
- Audio track preview before playback
- Cloud storage integration for audio tracks

## Troubleshooting

### Audio Track Not Playing
- **Verify format**: Ensure the audio file is in a supported format
- **Check file size**: Empty or corrupted files won't play
- **Test URL**: Try opening the URL in a browser first
- **Player support**: Some players have limited format support

### Audio Out of Sync
- Use video editing software to adjust timing
- Consider using a different audio source
- Check if the audio track matches the video version

### Download Fails
- **Check URL**: Ensure the URL is accessible
- **Check permissions**: Verify write access to temp directory
- **Network issues**: Ensure stable internet connection
- **Server issues**: Try again later if server is down

### Player Won't Start
- **Check player installation**: Verify MPV/VLC is installed
- **Check file paths**: Ensure audio file paths are correct
- **View console output**: Check terminal for detailed error messages

## Settings

Audio language preference is configured in:
**Settings > Audio Language**

Available languages:
- English (en)
- Spanish (es)
- French (fr)
- German (de)
- Italian (it)
- Portuguese (pt)
- Japanese (ja)
- Korean (ko)
- Chinese (zh)
- Arabic (ar)
- Russian (ru)
- Hindi (hi)

This preference is used to prioritize audio tracks when multiple are available.

## API Reference

### `audiotracks.Manager`

#### `NewManager() *Manager`
Creates a new audio track manager instance.

#### `DownloadAudioTrack(url string, filename string) (string, error)`
Downloads an audio track from a URL and returns the local file path.

**Parameters:**
- `url`: URL of the audio track to download
- `filename`: Desired filename for the downloaded track

**Returns:**
- Local file path to the downloaded audio track
- Error if download fails

#### `ValidateAudioFile(path string) error`
Validates that an audio file exists and is readable.

**Parameters:**
- `path`: Path to the audio file

**Returns:**
- Error if validation fails, nil if successful

### `player.PlayWithMPVAndAudio`

```go
func PlayWithMPVAndAudio(
    streamURL string, 
    title string, 
    subtitleURLs []string, 
    audioTrackURLs []string, 
    onEnd OnPlaybackEndCallback) error
```

Plays a stream with external audio tracks and subtitles.

**Parameters:**
- `streamURL`: URL of the video stream
- `title`: Title to display in player
- `subtitleURLs`: Array of subtitle file URLs/paths
- `audioTrackURLs`: Array of audio track URLs/paths
- `onEnd`: Callback function when playback ends

**Returns:**
- Error if playback fails to start

## Conclusion

The audio track feature provides powerful flexibility for multilingual content consumption. Whether you're learning a new language, prefer dubbed content, or need accessibility features, external audio tracks enhance your viewing experience.

For issues or feature requests, please refer to the main project repository.

