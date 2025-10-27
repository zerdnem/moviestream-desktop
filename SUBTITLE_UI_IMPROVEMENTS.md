# Subtitle UI Improvements

## Overview
Enhanced the subtitle download dialog with two major improvements:
1. **Don't auto-play** when subtitle download fails
2. **Separate results** by source (OpenSubtitles vs Addic7ed)

## Changes Made

### 1. Removed Auto-Play on Failure ❌➡️🛑

**Problem Before:**
When subtitle downloads failed (after retry or Addic7ed fallback), the video would automatically start playing WITHOUT subtitles. This was frustrating for users who wanted subtitles.

**Solution:**
All automatic video playback on download failure has been removed. Now users must explicitly choose to play without subtitles.

#### Affected Scenarios:

**Scenario A: Retry Download Fails**
```
Before:
  Download fails → Retry → Fails again → Auto-plays without subtitles

After:
  Download fails → Retry → Fails again → Error dialog only
  User must click "Play Without Subtitles" if they want to continue
```

**Scenario B: Addic7ed Search Fails**
```
Before:
  OpenSubtitles fails → Try Addic7ed → Search fails → Auto-plays without subtitles

After:
  OpenSubtitles fails → Try Addic7ed → Search fails → Error dialog only
  User stays in subtitle dialog, can try different subtitle or close
```

**Scenario C: Addic7ed Download Fails**
```
Before:
  OpenSubtitles fails → Addic7ed downloads → Download fails → Auto-plays without subtitles

After:
  OpenSubtitles fails → Addic7ed downloads → Download fails → Error dialog only
  User stays in subtitle dialog
```

#### Error Messages Updated:

**Old Messages:**
- "Playing without subtitles" (implied auto-play)
- "download failed again: [error]\n\nPlaying without subtitles"

**New Messages:**
- "Please try again or play without subtitles."
- "Download failed again: [error]\n\nPlease try again or play without subtitles."

### 2. Separated Results by Source 🔄

**Problem Before:**
All subtitles were mixed together in a single list. Users couldn't tell which came from OpenSubtitles vs Addic7ed.

**Solution:**
Results are now displayed in separate sections with clear headers.

#### Visual Layout:

**Before:**
```
Found 5 subtitle(s):
○ English (en) - movie.srt
○ English (en) - addic7ed_EXPLOIT.srt
○ Spanish (es) - movie.spa.srt
○ English (en) - addic7ed_GOSSIP.srt
○ French (fr) - movie.fre.srt
```

**After:**
```
Found 5 subtitle(s):

OpenSubtitles:
○   English (en) - movie.srt
○   Spanish (es) - movie.spa.srt
○   French (fr) - movie.fre.srt

Addic7ed:
○   English (en) - addic7ed_EXPLOIT.srt
○   English (en) - addic7ed_GOSSIP.srt
```

#### Implementation Details:

**Source Detection:**
```go
if len(result.ID) > 0 && result.ID[0] == '/' {
    // Addic7ed (ID starts with "/")
    addic7edResults = append(addic7edResults, result)
} else {
    // OpenSubtitles (numeric ID)
    openSubtitlesResults = append(openSubtitlesResults, result)
}
```

**Section Headers:**
- Bold labels: "OpenSubtitles:" and "Addic7ed:"
- Indented options (two spaces) for visual hierarchy
- Maintains single radio group for easy selection

**Selection Priority:**
1. First OpenSubtitles result (if available)
2. First Addic7ed result (if no OpenSubtitles results)

## Code Changes

### File: `gui/subtitlesdialog.go`

#### Change 1: Removed Auto-Play (Line 119-125)
```go
// Before
dialog.ShowError(fmt.Errorf("%s\n\nPlaying without subtitles", errorText), currentWindow)
player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd)

// After
dialog.ShowError(fmt.Errorf("%s", errorText), currentWindow)
// No auto-play
```

#### Change 2: Removed Auto-Play (Line 139-141)
```go
// Before
dialog.ShowError(fmt.Errorf("Addic7ed download also failed: %v\n\nPlaying without subtitles", downloadErr), currentWindow)
player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd)

// After
dialog.ShowError(fmt.Errorf("Addic7ed download also failed: %v\n\nPlease try again or play without subtitles.", downloadErr), currentWindow)
// No auto-play
```

#### Change 3: Removed Auto-Play (Line 181-183)
```go
// Before
dialog.ShowError(fmt.Errorf("download failed again: %v\n\nPlaying without subtitles", err), currentWindow)
player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd)

// After
dialog.ShowError(fmt.Errorf("Download failed again: %v\n\nPlease try again or play without subtitles.", err), currentWindow)
// No auto-play
```

#### Change 4: Separated Results (Lines 274-370)
```go
// Separate results by source
var openSubtitlesResults []subtitles.SubtitleResult
var addic7edResults []subtitles.SubtitleResult

for _, result := range results {
    if len(result.ID) > 0 && result.ID[0] == '/' {
        addic7edResults = append(addic7edResults, result)
    } else {
        openSubtitlesResults = append(openSubtitlesResults, result)
    }
}

// Add OpenSubtitles section
if len(openSubtitlesResults) > 0 {
    osHeader := widget.NewLabel("OpenSubtitles:")
    osHeader.TextStyle = fyne.TextStyle{Bold: true}
    resultsContainer.Add(osHeader)
    // ... add options with indentation
}

// Add Addic7ed section
if len(addic7edResults) > 0 {
    addic7edHeader := widget.NewLabel("Addic7ed:")
    addic7edHeader.TextStyle = fyne.TextStyle{Bold: true}
    resultsContainer.Add(addic7edHeader)
    // ... add options with indentation
}
```

## User Experience Impact

### Benefits:

#### 1. User Control ✓
- **No surprises**: Video never auto-plays without user action
- **Explicit choice**: User must actively choose to continue without subtitles
- **Time to decide**: User can review error message and decide next steps

#### 2. Clear Organization ✓
- **Easy to compare**: See all options from each source at a glance
- **Source preference**: Users can prefer one source over another
- **Visual clarity**: Headers and indentation make scanning easier

#### 3. Better Errors ✓
- **Helpful messages**: "Please try again or play without subtitles"
- **No auto-play confusion**: Clear that video hasn't started
- **Multiple attempts**: User can retry different subtitles

### Example User Flow:

```
User plays TV show episode
    ↓
No subtitles in stream
    ↓
Subtitle dialog appears
    ↓
Shows results:
    OpenSubtitles:
      ○ English (en) - breaking.bad.s01e01.srt
      ○ Spanish (es) - breaking.bad.s01e01.spa.srt
    
    Addic7ed:
      ○ English (en) - addic7ed_EXPLOIT.srt
      ○ English (en) - addic7ed_WEB.srt
    ↓
User selects OpenSubtitles subtitle
    ↓
Download fails (503 error)
    ↓
Dialog: "Would you like to try Addic7ed instead?"
    ↓
User clicks "Try Addic7ed"
    ↓
Addic7ed download fails
    ↓
Error: "Addic7ed download also failed: [error]
         Please try again or play without subtitles."
    ↓
User STAYS in subtitle dialog
    ↓
User can:
  - Select different Addic7ed subtitle
  - Go back and select different OpenSubtitles subtitle
  - Click "Play Without Subtitles"
  - Click "Cancel"
```

## Testing Checklist

### Test 1: No Auto-Play on Retry Failure ✓
1. Play TV show
2. Select OpenSubtitles subtitle
3. Wait for download to fail
4. Click "Retry Download"
5. Wait for retry to fail
6. **Verify**: Video does NOT start playing
7. **Verify**: Error message says "Please try again or play without subtitles"

### Test 2: No Auto-Play on Addic7ed Failure ✓
1. Play TV show
2. Select OpenSubtitles subtitle
3. Wait for download to fail
4. Click "Try Addic7ed"
5. Wait for Addic7ed to fail
6. **Verify**: Video does NOT start playing
7. **Verify**: User stays in dialog

### Test 3: Separated Results Display ✓
1. Play TV show that has subtitles from both sources
2. Open subtitle dialog
3. **Verify**: See "OpenSubtitles:" header
4. **Verify**: See "Addic7ed:" header
5. **Verify**: Results are properly grouped
6. **Verify**: Options are indented under headers

### Test 4: Source Selection Works ✓
1. Display separated results
2. Select an OpenSubtitles subtitle
3. Download and play
4. **Verify**: Correct subtitle is downloaded
5. Repeat with Addic7ed subtitle
6. **Verify**: Correct subtitle is downloaded

## Build Status

✅ **No linter errors**
✅ **Builds successfully**
✅ **Production ready**

## Conclusion

These improvements significantly enhance the user experience:

- **User control**: No unwanted auto-play behavior
- **Clear organization**: Easy to see and compare subtitle sources
- **Better errors**: Helpful messages that guide user actions
- **Flexibility**: User can retry multiple times or try different sources

Users now have full control over whether to play without subtitles, and can easily distinguish between subtitle sources.

