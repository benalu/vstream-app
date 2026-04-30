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
	Images struct {
		Logos []struct {
			FilePath string `json:"file_path"`
		} `json:"logos"`
	} `json:"images"`
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

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/%s/%s?api_key=%s&language=id-ID&append_to_response=images",
		tmdbType, tmdbID, apiKey,
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
