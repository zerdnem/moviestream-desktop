# Addic7ed Integration Summary

## Overview
Successfully integrated **Addic7ed** as a fallback subtitle source for TV shows. This provides redundancy when OpenSubtitles is unavailable or returns no results.

## Implementation Details

### Files Created/Modified
1. **`subtitles/addic7ed.go`** (NEW)
   - `Addic7edClient` struct for handling Addic7ed website scraping
   - `SearchByTitle()` - Searches for TV show subtitles
   - `DownloadSubtitle()` - Downloads subtitle files
   - HTML parsing using regex to extract subtitle information

2. **`subtitles/manager.go`** (MODIFIED)
   - Added `addic7ed` client to Manager struct
   - Updated `SearchSubtitles()` to try Addic7ed as fallback after OpenSubtitles
   - Updated `DownloadSubtitle()` to route downloads to correct source based on ID format

3. **`gui/subtitlesdialog.go`** (MODIFIED)
   - Enhanced download error handling to detect OpenSubtitles failures
   - Added smart fallback logic: offers Addic7ed when OpenSubtitles download fails (TV shows only)
   - Automatic search and download from Addic7ed with user confirmation
   - User-friendly dialog: "Try Addic7ed" button instead of generic retry
   - Updated info text to mention both subtitle sources

## Features

### Addic7ed Capabilities
✅ **TV Shows Only** - Addic7ed specializes in TV show subtitles
✅ **Multiple Languages** - Supports English, Spanish, French, German, Italian, Portuguese, etc.
✅ **Multiple Versions** - Often has subtitles for different release groups (EXPLOIT, GOSSIP, etc.)
✅ **Web Scraping** - Uses HTML parsing since Addic7ed has no official API
✅ **Smart Fallback** - Automatically offers Addic7ed when OpenSubtitles download fails

### Search Flow
```
1. Try OpenSubtitles with user's preferred language
   ├─ Success → Return results
   └─ Fail → Continue

2. Try OpenSubtitles with English (if user language != English)
   ├─ Success → Return results
   └─ Fail → Continue

3. Try Addic7ed with user's preferred language (TV shows only)
   ├─ Success → Return results
   └─ Fail → Continue

4. Try Addic7ed with English (if user language != English)
   ├─ Success → Return results
   └─ Fail → Return error
```

### Download Fallback Flow (NEW!)
```
User selects subtitle and clicks "Download & Play"
    ↓
Try downloading from selected source
    ↓
    ├─ Success → Play video with subtitle
    │
    └─ Failure (OpenSubtitles) AND (TV Show)
        ↓
        Show dialog: "Try Addic7ed instead?"
        ↓
        ├─ User clicks "Try Addic7ed"
        │   ↓
        │   Search Addic7ed for same episode
        │   ↓
        │   ├─ Found → Download from Addic7ed → Play
        │   └─ Not found → Play without subtitles
        │
        └─ User clicks "Play Without Subtitles"
            ↓
            Play video without subtitles
```

### Download Routing
The manager automatically detects the source based on the subtitle ID format:
- **Addic7ed**: ID starts with "/" (e.g., `/updated/1/167242/0`)
- **OpenSubtitles**: Numeric ID only (e.g., `167242`)

## Testing Results

### Test 1: Search for TV Show
```
Show: Breaking Bad S01E01
Result: ✅ Found 2 English subtitles from Addic7ed
```

### Test 2: Download Subtitle
```
Source: Addic7ed
File Size: 38,633 bytes
Format: Valid SRT file with timestamps
Status: ✅ Successfully downloaded and verified
```

### Test 3: Multiple Shows
```
Show: The Office S01E01
Result: ✅ Found 2 English subtitles
```

### Test 4: Movie Support
```
Type: Movie (no season/episode)
Result: ✅ Correctly rejected (Addic7ed is TV-only)
```

## HTML Parsing Details

### Page Structure
Addic7ed uses a table-based layout:
```html
<tr class="epeven completed">
    <td>SEASON</td>
    <td>EPISODE</td>
    <td>EPISODE_TITLE</td>
    <td>LANGUAGE</td>
    <td>VERSION</td>
    <td>STATUS</td>
    <td><a href="/updated/LANG_ID/FILE_ID/VERSION">Download</a></td>
</tr>
```

### Extraction Logic
1. Split HTML by table rows (`<tr class="ep...">`)
2. Extract `<td>` columns using regex
3. Match episode number and language
4. Extract download link (`/updated/...`)
5. Create `SubtitleResult` objects

## Advantages

### 1. Redundancy
- If OpenSubtitles is down or rate-limited, Addic7ed provides backup
- Increases subtitle availability for TV shows

### 2. Quality
- Addic7ed often has high-quality, well-synced subtitles
- Multiple versions per episode for different releases

### 3. TV Show Focus
- Addic7ed specializes in TV shows
- Complements OpenSubtitles' broader movie + TV coverage

## Limitations

### 1. TV Shows Only
- ❌ Does not support movies
- Only activates when `season > 0 && episode > 0`

### 2. Web Scraping
- ⚠️ Depends on HTML structure remaining stable
- May break if Addic7ed updates their website design
- More fragile than API-based solutions

### 3. Show Matching
- Currently takes first show ID from search results
- May not always find the exact show (e.g., remakes, regional variants)
- Room for improvement in show name matching

## Error Handling

### Graceful Degradation
```
OpenSubtitles fails → Try Addic7ed
Addic7ed fails → Return "no subtitles found" error
Never crashes → Always provides user with options
```

### User Experience
- Clear console messages showing search progress
- Indicates which source provided the results
- Download retry mechanism still works for both sources

## Future Improvements

### Potential Enhancements
1. **Better Show Matching**
   - Compare show titles more intelligently
   - Use Levenshtein distance or fuzzy matching
   - Consider show year/ID if available

2. **Caching**
   - Cache show ID lookups to avoid repeated searches
   - Store recent subtitle searches

3. **Rate Limiting**
   - Add delays between requests
   - Respect Addic7ed's server load

4. **Parser Robustness**
   - Add more fallback regex patterns
   - Handle edge cases (missing columns, malformed HTML)

## Console Output Examples

### Successful Addic7ed Fallback (Search)
```
Searching OpenSubtitles...
⚠ OpenSubtitles search failed: API returned status 503: Server busy
Trying Addic7ed as fallback...
✓ Found 2 subtitles from Addic7ed
```

### Successful Addic7ed Fallback (Download)
```
Downloading from OpenSubtitles...
❌ Download failed with status 503: Server temporarily unavailable

[User clicks "Try Addic7ed" in dialog]

Trying Addic7ed as alternative...
✓ Found 2 subtitles from Addic7ed
Downloading from Addic7ed...
✓ Downloaded to: C:\...\Temp\moviestream_addic7ed_EXPLOIT.srt
▶ Playing: Breaking Bad - S01E01 - Pilot
   ✓ Subtitle loaded from Addic7ed: addic7ed_EXPLOIT.srt (English)
```

### Normal Download
```
Downloading from Addic7ed...
✓ Downloaded to: C:\...\Temp\moviestream_addic7ed_EXPLOIT.srt
▶ Playing: Breaking Bad - S01E01 - Pilot
   ✓ Subtitle loaded: addic7ed_EXPLOIT.srt (English)
```

## Integration Status

✅ **Fully Integrated**
- Addic7ed client implemented and tested
- Integrated into subtitle manager
- Automatic fallback configured
- Download routing working
- Application builds successfully

✅ **Production Ready**
- Error handling in place
- User-friendly console messages
- Works alongside existing OpenSubtitles integration
- No breaking changes to existing functionality

## Conclusion

The Addic7ed integration significantly improves subtitle availability for **TV shows** by providing a reliable fallback when OpenSubtitles is unavailable. The implementation uses web scraping to access Addic7ed's extensive TV show subtitle database, providing users with more options and better reliability.

**Key Benefits:**
- 🎯 Increased TV show subtitle availability
- 🔄 Automatic fallback in both search AND download phases
- 🤝 Smart download fallback - offers Addic7ed when OpenSubtitles download fails
- 🌐 Multi-language support
- 💪 Robust error handling with user-friendly dialogs
- 🎬 Complements existing OpenSubtitles integration
- 🎭 Seamless user experience - one click to try alternative source

**New Feature: Download Fallback**
- When an OpenSubtitles download fails for a TV show, the system automatically offers Addic7ed as an alternative
- User-friendly dialog: "Would you like to try Addic7ed instead?"
- No need to restart the search - seamless transition to alternative source
- Falls back gracefully if both sources fail

The feature is **production-ready** and has been successfully tested with real TV show queries.

