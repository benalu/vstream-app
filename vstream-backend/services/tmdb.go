package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Struct internal untuk mapping JSON TMDB
type TMDBResponse struct {
	// Movie fields
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Runtime     int    `json:"runtime"`
	// TV fields
	Name            string `json:"name"`
	FirstAirDate    string `json:"first_air_date"`
	EpisodeRunTimes []int  `json:"episode_run_times"`
	// Shared
	VoteAverage  float64 `json:"vote_average"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	BackdropPath string  `json:"backdrop_path"`
	Genres       []struct {
		Name string `json:"name"`
	} `json:"genres"`
}

type LogoEntry struct {
	FilePath    string  `json:"file_path"`
	AspectRatio float64 `json:"aspect_ratio"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
}

func (r *TMDBResponse) NormalizedTitle() string {
	if r.Name != "" {
		return r.Name
	}
	return r.Title
}
func (r *TMDBResponse) NormalizedReleaseDate() string {
	if r.FirstAirDate != "" {
		return r.FirstAirDate
	}
	return r.ReleaseDate
}
func (r *TMDBResponse) NormalizedRuntime() int {
	if len(r.EpisodeRunTimes) > 0 {
		return r.EpisodeRunTimes[0]
	}
	return r.Runtime
}

func FetchMetadata(tmdbID string, contentType string) (*TMDBResponse, error) {
	apiKey := os.Getenv("TMDB_API_KEY")

	tmdbType := "movie"
	if contentType == "series" || contentType == "anime" {
		tmdbType = "tv"
	}

	// Primary: en-US (genres, title, overview, dll)
	primary, err := fetchTMDB(tmdbType, tmdbID, apiKey, "en-US")
	if err != nil {
		return nil, err
	}

	// Secondary: id-ID hanya untuk override overview
	secondary, err := fetchTMDB(tmdbType, tmdbID, apiKey, "id-ID")
	if err == nil && secondary.Overview != "" {
		primary.Overview = secondary.Overview
	}

	return primary, nil
}

func fetchTMDB(tmdbType, tmdbID, apiKey, lang string) (*TMDBResponse, error) {
	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/%s/%s?api_key=%s&language=%s",
		tmdbType, tmdbID, apiKey, lang,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB returned status: %d", resp.StatusCode)
	}

	var data TMDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func FetchLogos(tmdbID string, contentType string) ([]LogoEntry, error) {
	apiKey := os.Getenv("TMDB_API_KEY")

	tmdbType := "movie"
	if contentType == "series" || contentType == "anime" {
		tmdbType = "tv"
	}

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/%s/%s/images?api_key=%s&language=en-US",
		tmdbType, tmdbID, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB images returned status: %d", resp.StatusCode)
	}

	var data struct {
		Logos []struct {
			FilePath    string  `json:"file_path"`
			AspectRatio float64 `json:"aspect_ratio"`
			Width       int     `json:"width"`
			Height      int     `json:"height"`
			VoteAverage float64 `json:"vote_average"`
		} `json:"logos"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var result []LogoEntry
	for _, l := range data.Logos {
		if l.AspectRatio >= 2.0 && l.AspectRatio <= 14.0 {
			result = append(result, LogoEntry{
				FilePath:    "https://image.tmdb.org/t/p/original" + l.FilePath,
				AspectRatio: l.AspectRatio,
				Width:       l.Width,
				Height:      l.Height,
			})
			if len(result) >= 4 {
				break
			}
		}
	}

	return result, nil
}
