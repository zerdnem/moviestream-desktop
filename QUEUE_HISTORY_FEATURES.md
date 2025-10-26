# Queue and Watch History Features

This document describes the new queuing and watch history features added to MovieStream.

## Overview

Two major features have been implemented:
1. **Playback Queue** - Queue movies and TV episodes to play automatically in sequence
2. **Watch History** - Track and review what you've watched, with easy access to continue watching shows

---

## 1. Playback Queue

### Features

- **Add to Queue**: Add movies or TV episodes to a playback queue
- **Auto-Play**: Items in the queue automatically play after the current content finishes
- **Queue Management**: View, reorder, and remove items from the queue
- **Persistent**: Queue is saved and restored between app sessions
- **Integration**: Works seamlessly with auto-next for TV shows

### How to Use

#### Adding Content to Queue

**Movies:**
1. Search for a movie
2. Click "View Details"
3. Click "➕ Add to Queue" button
4. Movie will be added to the queue

**TV Episodes:**
1. Search for a TV show
2. Click "View Episodes"
3. Select a season
4. Click "➕ Add to Queue" on any episode
5. Episode will be added to the queue

#### Managing Queue

1. Click "📋 Queue" button on the main screen
2. View all queued items in order
3. Options for each item:
   - **▶ Play Now** - Remove from queue and play immediately
   - **Remove** - Remove from queue without playing
4. **Clear Queue** - Remove all items at once

### How Queue Works

1. When content finishes playing, the app checks if there are items in the queue
2. If queue has items, the next item automatically starts playing
3. For TV shows with auto-next enabled:
   - Auto-next takes priority (plays next episode in series)
   - When the series ends, queue items play next
4. Queue items are processed in FIFO (First In, First Out) order

---

## 2. Watch History

### Features

- **Automatic Tracking**: All watched movies and episodes are automatically recorded
- **Timestamps**: Shows when each item was watched (e.g., "2 hours ago", "Yesterday")
- **Continue Watching**: Quick access to resume the last watched TV show
- **Persistent**: History is saved between app sessions
- **Privacy**: Easily clear history or remove individual items

### How to Use

#### Viewing History

1. Click "🕒 History" button on the main screen
2. View chronological list of watched content (most recent first)
3. Each item shows:
   - Title and episode details (for TV shows)
   - When it was watched
   - Icon: 🎬 for movies, 📺 for TV shows

#### Continue Watching

1. Open watch history
2. Click "▶ Continue Watching" button at the top
3. Instantly resumes playback of the last watched TV episode

#### Managing History

- **Watch Again** - Play any item from history again
- **Remove** - Remove individual items from history
- **Clear History** - Remove all history at once

### How History Works

1. When you start watching content, it's automatically added to history
2. For TV shows, each episode is tracked separately
3. Duplicate entries are automatically removed (rewatching updates the timestamp)
4. History keeps the last 100 watched items

---

## Technical Implementation

### New Packages

#### `queue/queue.go`
- Queue data structure with thread-safe operations
- Persistence using Fyne preferences
- Support for both movies and TV episodes

#### `history/history.go`
- History tracking with timestamps
- Duplicate detection and removal
- Chronological sorting and display helpers

### UI Components

#### `gui/queueview.go`
- Queue display and management interface
- Play and remove controls

#### `gui/historyview.go`
- History display interface
- Continue watching functionality
- Time-based display formatting

### Integration Points

#### `main.go`
- Initialize queue and history systems on startup

#### `gui/app.go`
- Added Queue and History buttons to main UI
- Added "Add to Queue" buttons to movie details
- Integrated queue playback with movie player callbacks
- Added watch history recording for movies

#### `gui/tvdetails.go`
- Added "Add to Queue" buttons to episode lists
- Integrated queue playback with episode callbacks
- Added watch history recording for episodes
- Auto-next and queue work together seamlessly

#### `player/launcher.go`
- Existing callback mechanism used for auto-play features
- No changes needed (already supported callbacks)

---

## Data Storage

All data is stored using Fyne's preferences system:
- **Queue**: `playback_queue` preference (JSON)
- **History**: `watch_history` preference (JSON)

Data persists between app sessions and is automatically loaded on startup.

---

## Usage Examples

### Example 1: Movie Marathon
1. Search for "Lord of the Rings"
2. Add all movies to queue in order
3. Click watch on first movie
4. Sit back and enjoy - movies play automatically in sequence

### Example 2: Binge Watching with Queue
1. Add favorite episodes from different shows to queue
2. Watch them in the order you added them
3. Queue items play after each episode finishes

### Example 3: Continue Watching
1. Watch a few episodes of a show
2. Close the app
3. Reopen app later
4. Click History → Continue Watching
5. Resume from where you left off

---

## User Interface

### Main Screen
- New buttons added:
  - **📋 Queue** - Access playback queue
  - **🕒 History** - Access watch history

### Movie Details
- New button: **➕ Add to Queue**

### Episode List
- New button for each episode: **➕ Add to Queue**

### Playback Status Messages
- Shows queue status when playing content
- Example: "📋 5 item(s) in queue - will play next automatically"

---

## Benefits

1. **Seamless Viewing**: No need to manually select next content
2. **Custom Playlists**: Create your own viewing order across movies and shows
3. **Memory**: Never forget where you left off in a series
4. **Convenience**: Quick access to recently watched content
5. **Flexibility**: Works alongside existing auto-next feature

---

## Future Enhancements (Ideas)

- Reorder queue items by drag-and-drop
- Search/filter history
- Export/import queue and history
- Statistics (most watched shows, total watch time)
- Multiple named queues (e.g., "Weekend", "Favorites")
- Resume playback at exact timestamp (requires player integration)

---

## Troubleshooting

**Queue items not playing:**
- Ensure video player is properly installed
- Check that content is available/streamable
- Queue plays after current content finishes (be patient)

**History not recording:**
- History is recorded when playback starts
- Ensure app has write permissions for preferences

**Continue watching goes to wrong episode:**
- History tracks exact episode you watched
- If you manually select different episodes, history updates accordingly

---

## Code Quality

- ✅ Thread-safe operations with mutexes
- ✅ No linter errors
- ✅ Proper error handling
- ✅ Clean separation of concerns
- ✅ Comprehensive documentation
- ✅ Works with existing features (auto-next)

---

**Enjoy your enhanced MovieStream experience!** 🎬📺

