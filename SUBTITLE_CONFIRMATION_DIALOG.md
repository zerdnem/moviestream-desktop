# Subtitle Confirmation Dialog Feature

## Overview
Added **confirmation dialogs** that ask users if they want to play without subtitles when subtitle fetching or downloading fails. This gives users full control instead of automatically playing videos without subtitles.

## What Changed?

### Before
```
Subtitle fetch/download fails
    ↓
Video automatically plays without subtitles
    ↓
User: "Wait, I wanted subtitles!"
```

### After (NEW!)
```
Subtitle fetch/download fails
    ↓
Dialog: "Would you like to play without subtitles?"
    ↓
User chooses:
  → "Play Without Subtitles" - Video plays
  → "Cancel" - Nothing happens, user stays in app
```

## Dialog Scenarios

### 1. No Subtitles Found (Initial Search)

**When:** No subtitles found from any source during initial search

**Dialog:**
```
Title: "No Subtitles Found"
Message: "No subtitles were found for this content.

Would you like to play without subtitles?"

[Cancel] [Play Without Subtitles]
```

**Console Output:**
```
🔍 Searching for subtitles...
⚠ No subtitles found.
```

---

### 2. No Valid Subtitles

**When:** Search returned results but none were valid/usable

**Dialog:**
```
Title: "No Valid Subtitles"
Message: "No valid subtitles were found for this content.

Would you like to play without subtitles?"

[Cancel] [Play Without Subtitles]
```

**Console Output:**
```
✓ Auto-selected: [filename] from [source]
⚠ No valid subtitles found.
```

---

### 3. Subtitle Download Failed (Single Source)

**When:** Download from primary source failed, no alternative available (movies)

**Dialog:**
```
Title: "Subtitle Download Failed"
Message: "Failed to download subtitles.

Would you like to play without subtitles?"

[Cancel] [Play Without Subtitles]
```

**Console Output:**
```
✓ Auto-selected: movie.srt from OpenSubtitles
⬇ Downloading subtitle from OpenSubtitles...
⚠ Download failed: [error]
⚠ Subtitle download failed.
```

---

### 4. Both Sources Failed to Provide (TV Shows)

**When:** OpenSubtitles failed, Addic7ed search also failed

**Dialog:**
```
Title: "Subtitle Download Failed"
Message: "Both OpenSubtitles and Addic7ed failed to provide subtitles.

Would you like to play without subtitles?"

[Cancel] [Play Without Subtitles]
```

**Console Output:**
```
✓ Auto-selected: show.s01e01.srt from OpenSubtitles
⬇ Downloading subtitle from OpenSubtitles...
⚠ Download failed: [error]
🔄 Trying Addic7ed as alternative...
⚠ Addic7ed also failed.
```

---

### 5. All Subtitle Downloads Failed (TV Shows)

**When:** OpenSubtitles download failed, Addic7ed download also failed

**Dialog:**
```
Title: "All Subtitle Downloads Failed"
Message: "Failed to download subtitles from both OpenSubtitles and Addic7ed.

Would you like to play without subtitles?"

[Cancel] [Play Without Subtitles]
```

**Console Output:**
```
✓ Auto-selected: show.s01e01.srt from OpenSubtitles
⬇ Downloading subtitle from OpenSubtitles...
⚠ Download failed: [error]
🔄 Trying Addic7ed as alternative...
⬇ Downloading subtitle from Addic7ed...
⚠ Both sources failed.
```

## User Actions

### User Clicks "Play Without Subtitles"
- ✅ Video starts playing immediately
- ✅ No subtitles loaded
- ✅ Console shows: `▶ Playing: [title]` and `⚠ No subtitles loaded`
- ✅ Queue/auto-next features continue normally

### User Clicks "Cancel"
- ✅ Dialog closes
- ✅ Video does NOT start playing
- ✅ User returns to main app interface
- ✅ Console shows: `✗ Playback cancelled by user`
- ✅ User can try a different movie/show

## User Experience Flow

### Example 1: User Wants to Watch Without Subtitles
```
User clicks "Watch" on movie
    ↓
"Searching for subtitles..." (loader)
    ↓
No subtitles found
    ↓
Dialog: "No subtitles found. Play without?"
    ↓
User clicks "Play Without Subtitles"
    ↓
Video plays ✓
```

### Example 2: User Needs Subtitles
```
User clicks "Watch" on movie
    ↓
"Searching for subtitles..." (loader)
    ↓
No subtitles found
    ↓
Dialog: "No subtitles found. Play without?"
    ↓
User clicks "Cancel"
    ↓
Returns to app, can try different content ✓
```

### Example 3: Fallback Fails, User Still Wants to Watch
```
User clicks "Watch" on TV show
    ↓
"Searching..." → "Downloading from OpenSubtitles..."
    ↓
OpenSubtitles fails
    ↓
"Searching Addic7ed..." → Both failed
    ↓
Dialog: "Both sources failed. Play without?"
    ↓
User clicks "Play Without Subtitles"
    ↓
Video plays ✓
```

## Benefits

### 1. User Control ✓
- No surprises - users choose whether to continue
- Clear options at every failure point
- Can back out if subtitles are critical

### 2. Better UX ✓
- Explicit confirmation before playing
- User knows exactly what's happening
- Prevents unwanted playback

### 3. Clear Communication ✓
- Dialogs explain what failed
- User understands why no subtitles
- Console output shows full sequence

### 4. Flexibility ✓
- Users who need subtitles can cancel
- Users who don't care can continue
- No forced behavior

### 5. No Breaking Changes ✓
- Successful subtitle downloads work as before
- Only affects failure scenarios
- Queue and auto-next still work

## Technical Implementation

### Dialog Structure
All dialogs follow the same pattern:

```go
confirmDialog := dialog.NewConfirm(
    "Dialog Title",
    "Message explaining what happened.\n\nWould you like to play without subtitles?",
    func(play bool) {
        if play {
            // User chose to play without subtitles
            if err := player.PlayWithMPVAndCallback(streamURL, title, []string{}, onEnd); err != nil {
                dialog.ShowError(err, currentWindow)
            } else {
                fmt.Printf("\n▶ Playing: %s\n", title)
                fmt.Printf("   ⚠ No subtitles loaded\n\n")
            }
        } else {
            // User cancelled
            fmt.Println("✗ Playback cancelled by user")
        }
    },
    currentWindow,
)
confirmDialog.SetDismissText("Cancel")
confirmDialog.SetConfirmText("Play Without Subtitles")
confirmDialog.Show()
```

### Code Changes in `gui/autosubtitles.go`

**Line 33-62:** No subtitles found dialog
**Line 91-120:** No valid subtitles dialog
**Line 165-191:** Both sources failed to provide dialog
**Line 208-233:** All downloads failed dialog
**Line 247-273:** Single source failed dialog (movies)

Each failure point now shows a confirmation dialog instead of auto-playing.

## Testing Checklist

### Test 1: No Subtitles Found ✓
1. Play obscure movie with no subtitles
2. **Verify:** See "No Subtitles Found" dialog
3. Click "Play Without Subtitles"
4. **Verify:** Video plays

### Test 2: User Cancels ✓
1. Play movie with no subtitles
2. **Verify:** See confirmation dialog
3. Click "Cancel"
4. **Verify:** Video does NOT play
5. **Verify:** Can use app normally

### Test 3: Download Fails ✓
1. Play movie when OpenSubtitles is busy (503)
2. **Verify:** See download failure dialog
3. **Verify:** Dialog explains what happened
4. Choose action
5. **Verify:** Behavior matches choice

### Test 4: Both Sources Fail (TV) ✓
1. Play TV show when both sources unavailable
2. **Verify:** See multiple loaders
3. **Verify:** See "Both sources failed" dialog
4. Choose action
5. **Verify:** Behavior matches choice

### Test 5: Queue Integration ✓
1. Add multiple items to queue
2. First item has no subtitles
3. Click "Play Without Subtitles"
4. **Verify:** Queue continues to next item

### Test 6: Auto-Next Integration ✓
1. Play TV episode with auto-next
2. No subtitles for current episode
3. Click "Play Without Subtitles"
4. **Verify:** Episode plays, auto-next triggers

## Comparison

| Aspect | Before | After |
|--------|--------|-------|
| **Subtitle failure** | Auto-plays without | Asks user first |
| **User control** | None | Full control |
| **User awareness** | Surprised | Informed |
| **Can cancel** | No | Yes |
| **Flexibility** | Forced playback | User choice |

## Console Output Examples

### Scenario 1: User Confirms Playback
```
🔍 Searching for subtitles...
⚠ No subtitles found.

[User clicks "Play Without Subtitles"]

▶ Playing: Movie Title
   ⚠ No subtitles loaded
```

### Scenario 2: User Cancels
```
🔍 Searching for subtitles...
⚠ No subtitles found.

[User clicks "Cancel"]

✗ Playback cancelled by user
```

### Scenario 3: Fallback Failure
```
✓ Auto-selected: show.s01e01.srt from OpenSubtitles
⬇ Downloading subtitle from OpenSubtitles...
⚠ Download failed: API returned status 503
🔄 Trying Addic7ed as alternative...
⚠ Addic7ed also failed.

[User clicks "Play Without Subtitles"]

▶ Playing: Show Name - S01E01 - Episode Title
   ⚠ No subtitles loaded
```

## Build Status

✅ **No linter errors**
✅ **Builds successfully**
✅ **Production ready**

## User Feedback Impact

**Before:**
- "Why did it start playing without subtitles?"
- "I wanted subtitles, this is annoying!"
- "Can I go back and try something else?"

**After:**
- "Oh, no subtitles available. I'll watch anyway." ✓
- "No subtitles? I'll cancel and watch something else." ✓
- "I have full control over what happens." ✓

## Conclusion

The confirmation dialog feature provides **user control and transparency** when subtitle downloads fail. Users can now make informed decisions about whether to continue without subtitles or cancel and try different content.

**Key Benefits:**
- 🎯 User control at every failure point
- 💬 Clear communication about what failed
- ❌ Option to cancel if subtitles are critical
- ✅ Option to continue if subtitles aren't needed
- 🔄 Works with all existing features (queue, auto-next)

**Result:** Better user experience with no forced behavior - users decide what happens.

