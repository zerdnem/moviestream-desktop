package subtitles

import (
	"fmt"
	"moviestream-gui/settings"
)

// Manager handles subtitle search and download from multiple sources
type Manager struct {
	openSubtitles *OpenSubtitlesClient
	addic7ed      *Addic7edClient
}

// NewManager creates a new subtitle manager
func NewManager() *Manager {
	return &Manager{
		openSubtitles: NewOpenSubtitlesClient(),
		addic7ed:      NewAddic7edClient(),
	}
}

// SearchSubtitles searches for subtitles across all available sources
func (m *Manager) SearchSubtitles(title string, tmdbID int, season, episode int) ([]SubtitleResult, error) {
	userSettings := settings.Get()
	language := userSettings.SubtitleLanguage
	
	// Try OpenSubtitles first
	fmt.Println("Searching OpenSubtitles...")
	results, err := m.openSubtitles.SearchByTitle(title, language, season, episode)
	
	// If OpenSubtitles returns results, use them
	if err == nil && len(results) > 0 {
		fmt.Printf("✓ Found %d subtitles from OpenSubtitles\n", len(results))
		return results, nil
	}
	
	// If OpenSubtitles failed or returned no results, inform user
	if err != nil {
		fmt.Printf("⚠ OpenSubtitles search failed: %v\n", err)
	} else {
		fmt.Println("⚠ No results from OpenSubtitles")
	}
	
	// Try with English if language is not English
	if language != "en" {
		fmt.Println("Trying OpenSubtitles with English language...")
		results, err = m.openSubtitles.SearchByTitle(title, "en", season, episode)
		if err == nil && len(results) > 0 {
			fmt.Printf("✓ Found %d English subtitles from OpenSubtitles\n", len(results))
			return results, nil
		}
	}
	
	// Try Addic7ed as fallback (only for TV shows)
	if season > 0 && episode > 0 {
		fmt.Println("Trying Addic7ed as fallback...")
		results, err = m.addic7ed.SearchByTitle(title, language, season, episode)
		if err == nil && len(results) > 0 {
			fmt.Printf("✓ Found %d subtitles from Addic7ed\n", len(results))
			return results, nil
		}
		
		// Try Addic7ed with English if language is not English
		if language != "en" {
			fmt.Println("Trying Addic7ed with English language...")
			results, err = m.addic7ed.SearchByTitle(title, "en", season, episode)
			if err == nil && len(results) > 0 {
				fmt.Printf("✓ Found %d English subtitles from Addic7ed\n", len(results))
				return results, nil
			}
		}
		
		if err != nil {
			fmt.Printf("⚠ Addic7ed search failed: %v\n", err)
		}
	}
	
	// No results from any source
	return nil, fmt.Errorf("no subtitles found from any source. Please try again later")
}

// DownloadSubtitle downloads a subtitle from the given result
func (m *Manager) DownloadSubtitle(result SubtitleResult) (string, error) {
	// Determine source by ID format:
	// Addic7ed: starts with "/" (e.g., "/updated/1/12345/0")
	// OpenSubtitles: numeric only (e.g., "12345")
	
	if len(result.ID) > 0 && result.ID[0] == '/' {
		// Addic7ed subtitle
		fmt.Println("Downloading from Addic7ed...")
		return m.addic7ed.DownloadSubtitle(result.ID, result.FileName)
	}
	
	// OpenSubtitles subtitle (numeric ID)
	fmt.Println("Downloading from OpenSubtitles...")
	return m.openSubtitles.DownloadSubtitle(result.ID, result.FileName)
}

