# Subtitle Download Feature - Testing Guide

## Quick Test Scenarios

### Scenario 1: Movie Without Subtitles
**Steps:**
1. Launch MovieStream
2. Search for a movie (e.g., "The Matrix")
3. Click "Watch" on any result
4. Wait for stream to load

**Expected Behavior:**
- If no subtitles are embedded in the stream:
  - A dialog appears: "No subtitles found in stream"
  - Dialog shows: "Would you like to search for subtitles from OpenSubtitles?"
  - Progress indicator shows "Searching for subtitles..."
  - After search completes, available subtitles are listed with language info
  - User can select a subtitle and click "Download & Play"
  - Video plays with the downloaded subtitle

### Scenario 2: Movie With Subtitles
**Steps:**
1. Launch MovieStream
2. Search for a movie
3. Click "Watch" on a result that has embedded subtitles

**Expected Behavior:**
- No dialog appears
- Video plays immediately with embedded subtitles
- Console shows: "✓ X subtitle track(s) loaded"

### Scenario 3: TV Show Episode Without Subtitles
**Steps:**
1. Launch MovieStream
2. Search for a TV show (e.g., "Breaking Bad")
3. Select a season and episode
4. Click "Watch"

**Expected Behavior:**
- If no subtitles are embedded:
  - Subtitle download dialog appears
  - Search includes season and episode parameters
  - Shows episode-specific subtitle results
  - Can download and play with subtitles

### Scenario 4: Play Without Subtitles
**Steps:**
1. Trigger subtitle download dialog (by playing content without subtitles)
2. Click "Play Without Subtitles" button

**Expected Behavior:**
- Dialog closes
- Video starts playing immediately without subtitles
- Console shows: "⚠ No subtitles loaded"

### Scenario 5: Cancel Playback
**Steps:**
1. Trigger subtitle download dialog
2. Click "Cancel" button

**Expected Behavior:**
- Dialog closes
- No video playback starts
- Returns to main interface

### Scenario 6: Queue with Mixed Subtitle Availability
**Steps:**
1. Add multiple items to queue (some with subtitles, some without)
2. Start playing from queue
3. Observe behavior as queue progresses

**Expected Behavior:**
- Items with embedded subtitles play automatically
- Items without subtitles trigger the download dialog
- Queue progression maintains after subtitle selection/download

### Scenario 7: Language Preference
**Steps:**
1. Open Settings
2. Change subtitle language to Spanish
3. Play content without embedded subtitles

**Expected Behavior:**
- Subtitle search prioritizes Spanish subtitles
- Results show Spanish subtitles first if available
- Falls back to English if Spanish not available

### Scenario 8: No Results Found
**Steps:**
1. Play very obscure content that likely has no subtitles available

**Expected Behavior:**
- Dialog appears and searches
- Shows message: "No subtitles found. You can try playing without subtitles or search manually on OpenSubtitles.org"
- "Download & Play" button remains disabled
- Can still click "Play Without Subtitles"

## Testing Checklist

### Functional Tests
- [ ] Dialog appears when no subtitles found
- [ ] Dialog doesn't appear when subtitles are embedded
- [ ] Subtitle search completes successfully
- [ ] Search results display correctly with language info
- [ ] Subtitle selection works
- [ ] Download & Play button downloads and plays
- [ ] Play Without Subtitles button works
- [ ] Cancel button closes dialog without playing
- [ ] Subtitle loads in video player correctly
- [ ] Works with MPV
- [ ] Works with VLC
- [ ] Works with MPC-HC
- [ ] Works with PotPlayer

### Integration Tests
- [ ] Queue continues after subtitle download
- [ ] Auto-next works after subtitle download (TV shows)
- [ ] History records correctly after subtitle selection
- [ ] Language preference from settings is respected
- [ ] Multiple subtitle formats work (SRT, VTT)

### Edge Cases
- [ ] Network error during search (shows error message)
- [ ] Network error during download (shows error message)
- [ ] Very long title names (display correctly in dialog)
- [ ] Special characters in titles (handles correctly)
- [ ] Rapid clicking of buttons (no crashes)
- [ ] Dialog opened multiple times (no memory leaks)

### UI/UX Tests
- [ ] Dialog is properly sized and scrollable
- [ ] Loading indicator is visible during search
- [ ] Results are readable and well-formatted
- [ ] Buttons are clearly labeled
- [ ] Dialog is modal (blocks other actions)
- [ ] Dialog centers on screen

## Manual Testing Commands

### Test Build
```bash
go build
```

### Run Application
```bash
./moviestream-gui
```

### Check Console Output
Watch for these messages:
- `⚠ No subtitles available for this content` - Subtitle dialog should appear
- `✓ X subtitle track(s) loaded` - Embedded subtitles found
- `✓ Subtitle loaded: filename.srt (Language)` - Downloaded subtitle loaded

## Debugging Tips

### If Dialog Doesn't Appear
1. Check console for subtitle count: should be 0
2. Verify `ShowSubtitleDownloadDialog` is called
3. Check for any error messages

### If Search Fails
1. Check network connectivity
2. Verify OpenSubtitles API is accessible
3. Check console for API error messages
4. Try with a different title

### If Download Fails
1. Check temp directory permissions
2. Verify download URL is valid
3. Check network connectivity during download
4. Verify sufficient disk space

### If Video Doesn't Play with Subtitle
1. Check that subtitle file was created in temp directory
2. Verify video player supports the subtitle format
3. Check player launch arguments in console
4. Try loading subtitle manually in player

## Performance Notes

### Expected Timing
- Subtitle search: 1-3 seconds
- Subtitle download: 0.5-2 seconds
- Total delay from "no subs" to "playing": 2-5 seconds (user interaction included)

### API Rate Limits
- OpenSubtitles allows reasonable request rates for free users
- If testing extensively, add small delays between tests
- Consider implementing request caching for production use

## Known Limitations

1. **API Dependency**: Requires internet connection and OpenSubtitles API availability
2. **Search Accuracy**: Results depend on title matching accuracy
3. **Language Coverage**: Not all content has subtitles in all languages
4. **No Authentication**: Using free API tier without authentication
5. **Temporary Files**: Subtitles are downloaded to temp directory and may be cleaned up by OS

## Future Testing Considerations

- Load testing with multiple simultaneous downloads
- Long-term testing for memory leaks
- Testing with slow network connections
- Testing with different subtitle encodings
- Testing subtitle sync accuracy

