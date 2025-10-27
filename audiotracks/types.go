package audiotracks

// AudioTrackResult represents an audio track search result
type AudioTrackResult struct {
	ID           string
	Language     string
	LanguageName string
	Format       string  // e.g., "aac", "mp3", "opus"
	Codec        string
	Quality      string  // e.g., "128kbps", "320kbps"
	FileName     string
	DownloadURL  string
	FileSize     int64   // Size in bytes
	Source       string  // Source of the audio track
}

// AudioTrackInfo contains information about available audio tracks
type AudioTrackInfo struct {
	Tracks []AudioTrackResult
	Error  error
}

