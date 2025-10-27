# Download Fallback Feature

## Overview
Enhanced the subtitle download system to **intelligently suggest Addic7ed** as an alternative when OpenSubtitles download fails for TV shows.

## What Changed?

### Before
```
OpenSubtitles download fails
    ↓
Dialog: "Retry Download" or "Play Without Subtitles"
    ↓
User clicks "Retry Download"
    ↓
Same OpenSubtitles download is attempted again
    ↓
Often fails again (503 server busy)
    ↓
Play without subtitles
```

### After (NEW!)
```
OpenSubtitles download fails (TV show)
    ↓
System detects: "This is OpenSubtitles + This is a TV show"
    ↓
Dialog: "Try Addic7ed instead?" or "Play Without Subtitles"
    ↓
User clicks "Try Addic7ed"
    ↓
System searches Addic7ed for the same episode
    ↓
Found subtitle → Download from Addic7ed → Play with subtitle! ✓
```

## User Experience

### Dialog Message (OpenSubtitles Failure)
```
OpenSubtitles download failed:
API returned status 503: Server temporarily unavailable

Would you like to try Addic7ed instead?
(Addic7ed is a reliable alternative for TV shows)

[Play Without Subtitles]  [Try Addic7ed]
```

### Dialog Message (Other Failures)
```
Failed to download subtitle:
Network error

Options:
• Click OK to retry
• Or play without subtitles

[Play Without Subtitles]  [Retry Download]
```

## Smart Detection

The system automatically detects when to offer Addic7ed:

### Offers Addic7ed When:
✅ Download source was **OpenSubtitles** (detected by ID format)
✅ Content is a **TV show** (season > 0 AND episode > 0)
✅ Initial download **failed**

### Uses Standard Retry When:
- Download source was **Addic7ed** (already a fallback)
- Content is a **movie** (Addic7ed doesn't support movies)
- Initial download was successful

## Technical Implementation

### Detection Logic
```go
isOpenSubtitles := len(selectedSub.ID) > 0 && selectedSub.ID[0] != '/'
isTVShow := season > 0 && episode > 0

if isOpenSubtitles && isTVShow {
    // Offer Addic7ed
} else {
    // Standard retry
}
```

### ID Format Detection
- **OpenSubtitles**: Numeric ID like `"167242"`
- **Addic7ed**: Path like `"/updated/1/167242/0"`

The system uses this to determine which source failed and route accordingly.

## Fallback Sequence

### Complete Fallback Chain
```
1. Try OpenSubtitles download
   ├─ Success → Play ✓
   └─ Fail (TV show only) → Continue

2. User offered "Try Addic7ed"
   ├─ User declines → Play without subtitles
   └─ User accepts → Continue

3. Search Addic7ed for same episode
   ├─ Not found → Play without subtitles
   └─ Found → Continue

4. Download from Addic7ed
   ├─ Success → Play with Addic7ed subtitle ✓
   └─ Fail → Play without subtitles
```

## Benefits

### 1. Higher Success Rate
- When OpenSubtitles is down/busy (503 errors), Addic7ed provides backup
- TV show users get **second chance** at finding subtitles

### 2. User-Friendly
- Clear explanation of what went wrong
- Suggests viable alternative automatically
- One-click solution (no manual searching)

### 3. Intelligent
- Only suggests Addic7ed when it makes sense (TV shows)
- Doesn't waste time on movies or other content types
- Detects source automatically

### 4. Non-Intrusive
- Doesn't force the alternative
- User can still choose "Play Without Subtitles"
- Graceful degradation at each step

## Code Changes

### Modified File: `gui/subtitlesdialog.go`

**Lines 93-211**: Enhanced download error handling

**Key Changes:**
1. Added source detection: `isOpenSubtitles` check
2. Added content type detection: `isTVShow` check
3. Created Addic7ed-specific error dialog
4. Implemented automatic Addic7ed search on user confirmation
5. Added progress indicators for each step
6. Preserved original retry logic for non-OpenSubtitles failures

**New Dependencies:**
- `subtitles.NewAddic7edClient()` - Direct client access for fallback search

## User Scenarios

### Scenario 1: OpenSubtitles Busy (Common)
```
1. User plays Breaking Bad S01E01
2. No subtitles in stream → Dialog appears
3. User selects OpenSubtitles subtitle
4. OpenSubtitles returns 503 (server busy)
5. Dialog: "Would you like to try Addic7ed instead?"
6. User clicks "Try Addic7ed"
7. Addic7ed search finds 2 subtitles
8. First subtitle downloads successfully
9. Video plays with Addic7ed subtitle ✓
```

### Scenario 2: Both Sources Fail (Rare)
```
1. User plays obscure show S05E12
2. OpenSubtitles download fails (503)
3. Dialog offers Addic7ed
4. User clicks "Try Addic7ed"
5. Addic7ed search finds no results
6. Error: "Addic7ed search also failed"
7. Video plays without subtitles
```

### Scenario 3: Movie Download Fails
```
1. User plays The Matrix (movie)
2. OpenSubtitles download fails
3. Dialog: "Retry Download" (standard retry)
4. User clicks retry
5. Retries same source
(Addic7ed not offered - movies not supported)
```

## Console Output

### Successful Fallback
```
Downloading from OpenSubtitles...
❌ Download request failed with status 503: Server temporarily unavailable

[User clicks "Try Addic7ed"]

Trying Addic7ed as alternative...
✓ Found 2 subtitles from Addic7ed
Downloading from Addic7ed...
✓ Downloaded to: C:\...\Temp\moviestream_addic7ed_EXPLOIT.srt

▶ Playing: Breaking Bad - S01E01 - Pilot
   ✓ Subtitle loaded from Addic7ed: addic7ed_EXPLOIT.srt (English)
```

## Testing

### To Test This Feature:

1. **Simulate OpenSubtitles Failure:**
   - Play a TV show episode
   - Wait for subtitle dialog
   - Select an OpenSubtitles subtitle
   - When OpenSubtitles is busy/rate-limited, download will fail
   - Dialog should offer "Try Addic7ed"

2. **Verify Addic7ed Fallback:**
   - Click "Try Addic7ed"
   - Should see progress indicator
   - Should see "Found X subtitles from Addic7ed"
   - Should download and play successfully

3. **Verify Movie Behavior:**
   - Play a movie
   - If download fails, should see standard "Retry Download"
   - Should NOT see "Try Addic7ed" option

## Future Enhancements

### Potential Improvements:
1. **Parallel Search**: Search both sources simultaneously
2. **User Preference**: Remember if user prefers Addic7ed
3. **Quality Indicators**: Show download success rates per source
4. **Manual Source Selection**: Let user choose source upfront

## Status

✅ **Fully Implemented**
✅ **Tested with Addic7ed**
✅ **Builds Successfully**
✅ **Production Ready**

This feature significantly improves the reliability of subtitle downloads for TV shows by providing an intelligent, automatic fallback mechanism.

