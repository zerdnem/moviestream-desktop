package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// GetStreamURLWithBrowser uses headless Chrome to capture the stream URL and subtitles from network requests
func GetStreamURLWithBrowser(tmdbID int, contentType string, season, episode int) (*StreamInfo, error) {
	var embedURL string
	
	if contentType == "tv" {
		embedURL = fmt.Sprintf("%s/tv/%d/%d/%d", MoviesBaseURL, tmdbID, season, episode)
	} else {
		embedURL = fmt.Sprintf("%s/movie/%d", MoviesBaseURL, tmdbID)
	}

	fmt.Printf("DEBUG: Loading page with headless browser: %s\n", embedURL)

	// Create context with timeout
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout for the whole operation
	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	// Collect m3u8 URLs and VTT subtitles
	m3u8URLs := make([]string, 0)
	vttURLs := make(map[string]string) // label -> URL
	
	// Enable network events
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			resp := ev.Response
			// Look for m3u8 streams
			if strings.Contains(resp.URL, ".m3u8") {
				fmt.Printf("DEBUG: Found m3u8 URL: %s\n", resp.URL)
				m3u8URLs = append(m3u8URLs, resp.URL)
			}
			// Look for VTT subtitles
			if strings.Contains(resp.URL, ".vtt") {
				fmt.Printf("DEBUG: Found VTT subtitle: %s\n", resp.URL)
				// Try to extract language from URL
				label := "Unknown"
				if strings.Contains(strings.ToLower(resp.URL), "english") || strings.Contains(strings.ToLower(resp.URL), "eng") {
					label = "English"
				}
				vttURLs[label] = resp.URL
			}
		}
	})

	// Run the browser automation
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(embedURL),
		chromedp.Sleep(8*time.Second), // Wait for page to load and make requests
	)

	if err != nil {
		return nil, fmt.Errorf("browser automation failed: %v", err)
	}

	// Return the best quality stream (usually the first or last one)
	var streamURL string
	if len(m3u8URLs) > 0 {
		// Prefer "playlist.m3u8" if available, otherwise take the last one
		for _, url := range m3u8URLs {
			if strings.Contains(url, "playlist.m3u8") || strings.Contains(url, "cGxheWxpc3Q") {
				fmt.Printf("DEBUG: Selected playlist m3u8: %s\n", url)
				streamURL = url
				break
			}
		}
		if streamURL == "" {
			streamURL = m3u8URLs[len(m3u8URLs)-1]
			fmt.Printf("DEBUG: Selected last m3u8: %s\n", streamURL)
		}
	} else {
		return nil, fmt.Errorf("no m3u8 streams found in network requests")
	}

	// Build subtitle tracks
	var subtitles []SubtitleTrack
	for label, url := range vttURLs {
		subtitles = append(subtitles, SubtitleTrack{Label: label, URL: url})
		fmt.Printf("DEBUG: Added subtitle track: %s - %s\n", label, url)
	}

	if len(subtitles) > 0 {
		fmt.Printf("✓ Found %d subtitle track(s)\n", len(subtitles))
	} else {
		fmt.Printf("⚠ No subtitles found\n")
	}

	return &StreamInfo{
		StreamURL:    streamURL,
		SubtitleURLs: subtitles,
	}, nil
}

