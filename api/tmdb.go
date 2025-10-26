package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	TMDBAPIKey  = "a46c50a0ccb1bafe2b15665df7fad7e1"
	TMDBBaseURL = "https://api.themoviedb.org/3"
)

type Movie struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Overview     string  `json:"overview"`
	ReleaseDate  string  `json:"release_date"`
	PosterPath   string  `json:"poster_path"`
	BackdropPath string  `json:"backdrop_path"`
	VoteAverage  float64 `json:"vote_average"`
}

type TVShow struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Overview     string  `json:"overview"`
	FirstAirDate string  `json:"first_air_date"`
	PosterPath   string  `json:"poster_path"`
	VoteAverage  float64 `json:"vote_average"`
	NumSeasons   int     `json:"number_of_seasons"`
	NumEpisodes  int     `json:"number_of_episodes"`
}

type SearchMovieResponse struct {
	Results []Movie `json:"results"`
}

type SearchTVResponse struct {
	Results []TVShow `json:"results"`
}

type Season struct {
	SeasonNumber int       `json:"season_number"`
	Name         string    `json:"name"`
	Episodes     []Episode `json:"episodes"`
}

type Episode struct {
	EpisodeNumber int    `json:"episode_number"`
	Name          string `json:"name"`
	Overview      string `json:"overview"`
	StillPath     string `json:"still_path"`
}

type TVDetails struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Overview     string   `json:"overview"`
	FirstAirDate string   `json:"first_air_date"`
	PosterPath   string   `json:"poster_path"`
	BackdropPath string   `json:"backdrop_path"`
	VoteAverage  float64  `json:"vote_average"`
	Seasons      []Season `json:"seasons"`
}

// SearchMovies searches for movies by name
func SearchMovies(query string) ([]Movie, error) {
	params := url.Values{}
	params.Add("api_key", TMDBAPIKey)
	params.Add("query", query)
	params.Add("language", "en-US")

	resp, err := http.Get(fmt.Sprintf("%s/search/movie?%s", TMDBBaseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SearchMovieResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

// SearchTVShows searches for TV shows by name
func SearchTVShows(query string) ([]TVShow, error) {
	params := url.Values{}
	params.Add("api_key", TMDBAPIKey)
	params.Add("query", query)
	params.Add("language", "en-US")

	resp, err := http.Get(fmt.Sprintf("%s/search/tv?%s", TMDBBaseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SearchTVResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

// GetMovieDetails gets detailed information about a movie
func GetMovieDetails(movieID int) (*Movie, error) {
	params := url.Values{}
	params.Add("api_key", TMDBAPIKey)
	params.Add("language", "en-US")

	resp, err := http.Get(fmt.Sprintf("%s/movie/%d?%s", TMDBBaseURL, movieID, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var movie Movie
	if err := json.Unmarshal(body, &movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

// GetTVDetails gets detailed information about a TV show including seasons and episodes
func GetTVDetails(tvID int) (*TVDetails, error) {
	params := url.Values{}
	params.Add("api_key", TMDBAPIKey)
	params.Add("language", "en-US")

	resp, err := http.Get(fmt.Sprintf("%s/tv/%d?%s", TMDBBaseURL, tvID, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tvShow TVDetails
	if err := json.Unmarshal(body, &tvShow); err != nil {
		return nil, err
	}

	return &tvShow, nil
}

// GetSeasonDetails gets episodes for a specific season
func GetSeasonDetails(tvID, seasonNum int) (*Season, error) {
	params := url.Values{}
	params.Add("api_key", TMDBAPIKey)
	params.Add("language", "en-US")

	resp, err := http.Get(fmt.Sprintf("%s/tv/%d/season/%d?%s", TMDBBaseURL, tvID, seasonNum, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var season Season
	if err := json.Unmarshal(body, &season); err != nil {
		return nil, err
	}

	return &season, nil
}

// GetPosterURL returns the full URL for a poster image
func GetPosterURL(posterPath string) string {
	if posterPath == "" {
		return ""
	}
	return "https://image.tmdb.org/t/p/w500" + posterPath
}

// GetBackdropURL returns the full URL for a backdrop image
func GetBackdropURL(backdropPath string) string {
	if backdropPath == "" {
		return ""
	}
	return "https://image.tmdb.org/t/p/w1280" + backdropPath
}

