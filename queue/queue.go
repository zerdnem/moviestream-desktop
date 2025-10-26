package queue

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

// QueueItem represents an item in the playback queue
type QueueItem struct {
	ID           string    `json:"id"`           // Unique identifier
	Type         string    `json:"type"`         // "movie" or "tv"
	TMDBID       int       `json:"tmdb_id"`      // TMDB ID
	Title        string    `json:"title"`        // Display title
	Season       int       `json:"season"`       // For TV shows
	Episode      int       `json:"episode"`      // For TV shows
	EpisodeName  string    `json:"episode_name"` // Episode name for TV shows
	AddedAt      time.Time `json:"added_at"`     // When it was added to queue
}

// Queue manages the playback queue
type Queue struct {
	items []QueueItem
	mutex sync.RWMutex
	app   fyne.App
}

var (
	globalQueue *Queue
)

// Initialize sets up the queue system
func Initialize(app fyne.App) {
	globalQueue = &Queue{
		items: make([]QueueItem, 0),
		app:   app,
	}
	globalQueue.Load()
}

// Get returns the global queue instance
func Get() *Queue {
	if globalQueue == nil {
		return &Queue{
			items: make([]QueueItem, 0),
		}
	}
	return globalQueue
}

// Add adds an item to the queue
func (q *Queue) Add(item QueueItem) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Generate ID if not provided
	if item.ID == "" {
		item.ID = generateID(item)
	}
	item.AddedAt = time.Now()

	q.items = append(q.items, item)
	q.save()
}

// AddMovie adds a movie to the queue
func (q *Queue) AddMovie(tmdbID int, title string) {
	item := QueueItem{
		Type:   "movie",
		TMDBID: tmdbID,
		Title:  title,
	}
	q.Add(item)
}

// AddEpisode adds a TV episode to the queue
func (q *Queue) AddEpisode(tmdbID int, showName string, season, episode int, episodeName string) {
	item := QueueItem{
		Type:        "tv",
		TMDBID:      tmdbID,
		Title:       showName,
		Season:      season,
		Episode:     episode,
		EpisodeName: episodeName,
	}
	q.Add(item)
}

// Remove removes an item from the queue by ID
func (q *Queue) Remove(id string) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, item := range q.items {
		if item.ID == id {
			q.items = append(q.items[:i], q.items[i+1:]...)
			q.save()
			return true
		}
	}
	return false
}

// RemoveAt removes an item at a specific index
func (q *Queue) RemoveAt(index int) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if index < 0 || index >= len(q.items) {
		return false
	}

	q.items = append(q.items[:index], q.items[index+1:]...)
	q.save()
	return true
}

// GetNext returns and removes the next item in the queue
func (q *Queue) GetNext() (*QueueItem, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.items) == 0 {
		return nil, false
	}

	item := q.items[0]
	q.items = q.items[1:]
	q.save()
	return &item, true
}

// Peek returns the next item without removing it
func (q *Queue) Peek() (*QueueItem, bool) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if len(q.items) == 0 {
		return nil, false
	}

	return &q.items[0], true
}

// GetAll returns all items in the queue
func (q *Queue) GetAll() []QueueItem {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	// Return a copy to prevent external modification
	items := make([]QueueItem, len(q.items))
	copy(items, q.items)
	return items
}

// Clear removes all items from the queue
func (q *Queue) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = make([]QueueItem, 0)
	q.save()
}

// IsEmpty returns true if the queue is empty
func (q *Queue) IsEmpty() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return len(q.items) == 0
}

// Size returns the number of items in the queue
func (q *Queue) Size() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return len(q.items)
}

// Load loads the queue from persistent storage
func (q *Queue) Load() {
	if q.app == nil {
		return
	}

	prefs := q.app.Preferences()
	queueJSON := prefs.String("playback_queue")
	
	if queueJSON == "" {
		q.items = make([]QueueItem, 0)
		return
	}

	var items []QueueItem
	if err := json.Unmarshal([]byte(queueJSON), &items); err != nil {
		fmt.Printf("Error loading queue: %v\n", err)
		q.items = make([]QueueItem, 0)
		return
	}

	q.items = items
}

// save saves the queue to persistent storage
func (q *Queue) save() {
	if q.app == nil {
		return
	}

	data, err := json.Marshal(q.items)
	if err != nil {
		fmt.Printf("Error saving queue: %v\n", err)
		return
	}

	prefs := q.app.Preferences()
	prefs.SetString("playback_queue", string(data))
}

// GetDisplayTitle returns a formatted display title for a queue item
func (item *QueueItem) GetDisplayTitle() string {
	if item.Type == "movie" {
		return fmt.Sprintf("🎬 %s", item.Title)
	}
	return fmt.Sprintf("📺 %s - S%dE%d: %s", item.Title, item.Season, item.Episode, item.EpisodeName)
}

// generateID generates a unique ID for a queue item
func generateID(item QueueItem) string {
	if item.Type == "movie" {
		return fmt.Sprintf("movie-%d-%d", item.TMDBID, time.Now().Unix())
	}
	return fmt.Sprintf("tv-%d-s%de%d-%d", item.TMDBID, item.Season, item.Episode, time.Now().Unix())
}

