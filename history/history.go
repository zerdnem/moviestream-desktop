package history

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

// HistoryItem represents a watched item
type HistoryItem struct {
	Type        string    `json:"type"`         // "movie" or "tv"
	TMDBID      int       `json:"tmdb_id"`      // TMDB ID
	Title       string    `json:"title"`        // Display title
	Season      int       `json:"season"`       // For TV shows
	Episode     int       `json:"episode"`      // For TV shows
	EpisodeName string    `json:"episode_name"` // Episode name for TV shows
	WatchedAt   time.Time `json:"watched_at"`   // When it was watched
}

// History manages watch history
type History struct {
	items []HistoryItem
	mutex sync.RWMutex
	app   fyne.App
}

var (
	globalHistory *History
)

// Initialize sets up the history system
func Initialize(app fyne.App) {
	globalHistory = &History{
		items: make([]HistoryItem, 0),
		app:   app,
	}
	globalHistory.Load()
}

// Get returns the global history instance
func Get() *History {
	if globalHistory == nil {
		return &History{
			items: make([]HistoryItem, 0),
		}
	}
	return globalHistory
}

// AddMovie adds a movie to watch history
func (h *History) AddMovie(tmdbID int, title string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	item := HistoryItem{
		Type:      "movie",
		TMDBID:    tmdbID,
		Title:     title,
		WatchedAt: time.Now(),
	}

	// Remove duplicates (same movie)
	h.items = h.removeDuplicates(h.items, item)
	
	// Add to beginning of list (most recent first)
	h.items = append([]HistoryItem{item}, h.items...)
	
	// Keep only last 100 items
	if len(h.items) > 100 {
		h.items = h.items[:100]
	}

	h.save()
}

// AddEpisode adds a TV episode to watch history
func (h *History) AddEpisode(tmdbID int, showName string, season, episode int, episodeName string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	item := HistoryItem{
		Type:        "tv",
		TMDBID:      tmdbID,
		Title:       showName,
		Season:      season,
		Episode:     episode,
		EpisodeName: episodeName,
		WatchedAt:   time.Now(),
	}

	// Remove duplicates (same episode)
	h.items = h.removeDuplicates(h.items, item)
	
	// Add to beginning of list (most recent first)
	h.items = append([]HistoryItem{item}, h.items...)
	
	// Keep only last 100 items
	if len(h.items) > 100 {
		h.items = h.items[:100]
	}

	h.save()
}

// GetLastWatched returns the most recently watched item
func (h *History) GetLastWatched() (*HistoryItem, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if len(h.items) == 0 {
		return nil, false
	}

	return &h.items[0], true
}

// GetLastWatchedShow returns the most recently watched TV show with its last episode
func (h *History) GetLastWatchedShow() (*HistoryItem, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for i := range h.items {
		if h.items[i].Type == "tv" {
			return &h.items[i], true
		}
	}

	return nil, false
}

// GetShowHistory returns all episodes watched for a specific show
func (h *History) GetShowHistory(tmdbID int) []HistoryItem {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var showItems []HistoryItem
	for _, item := range h.items {
		if item.Type == "tv" && item.TMDBID == tmdbID {
			showItems = append(showItems, item)
		}
	}

	return showItems
}

// GetAll returns all history items
func (h *History) GetAll() []HistoryItem {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// Return a copy to prevent external modification
	items := make([]HistoryItem, len(h.items))
	copy(items, h.items)
	return items
}

// Clear removes all items from history
func (h *History) Clear() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.items = make([]HistoryItem, 0)
	h.save()
}

// Remove removes a specific item from history
func (h *History) Remove(index int) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if index < 0 || index >= len(h.items) {
		return false
	}

	h.items = append(h.items[:index], h.items[index+1:]...)
	h.save()
	return true
}

// Load loads the history from persistent storage
func (h *History) Load() {
	if h.app == nil {
		return
	}

	prefs := h.app.Preferences()
	historyJSON := prefs.String("watch_history")
	
	if historyJSON == "" {
		h.items = make([]HistoryItem, 0)
		return
	}

	var items []HistoryItem
	if err := json.Unmarshal([]byte(historyJSON), &items); err != nil {
		fmt.Printf("Error loading history: %v\n", err)
		h.items = make([]HistoryItem, 0)
		return
	}

	// Sort by watched time (most recent first)
	sort.Slice(items, func(i, j int) bool {
		return items[i].WatchedAt.After(items[j].WatchedAt)
	})

	h.items = items
}

// save saves the history to persistent storage
func (h *History) save() {
	if h.app == nil {
		return
	}

	data, err := json.Marshal(h.items)
	if err != nil {
		fmt.Printf("Error saving history: %v\n", err)
		return
	}

	prefs := h.app.Preferences()
	prefs.SetString("watch_history", string(data))
}

// removeDuplicates removes duplicate entries
func (h *History) removeDuplicates(items []HistoryItem, newItem HistoryItem) []HistoryItem {
	var filtered []HistoryItem
	for _, item := range items {
		if newItem.Type == "movie" {
			// For movies, check TMDB ID
			if item.Type != "movie" || item.TMDBID != newItem.TMDBID {
				filtered = append(filtered, item)
			}
		} else {
			// For TV shows, check TMDB ID, season, and episode
			if item.Type != "tv" || item.TMDBID != newItem.TMDBID || 
			   item.Season != newItem.Season || item.Episode != newItem.Episode {
				filtered = append(filtered, item)
			}
		}
	}
	return filtered
}

// GetDisplayTitle returns a formatted display title for a history item
func (item *HistoryItem) GetDisplayTitle() string {
	if item.Type == "movie" {
		return fmt.Sprintf("🎬 %s", item.Title)
	}
	return fmt.Sprintf("📺 %s - S%dE%d: %s", item.Title, item.Season, item.Episode, item.EpisodeName)
}

// GetDisplayTime returns a human-readable time string
func (item *HistoryItem) GetDisplayTime() string {
	duration := time.Since(item.WatchedAt)
	
	if duration < time.Minute {
		return "Just now"
	} else if duration < time.Hour {
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else {
		return item.WatchedAt.Format("Jan 2, 2006")
	}
}

