package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Struct internal untuk mapping JSON TMDB
type TMDBResponse struct {
	Title        string  `json:"title"`
	ReleaseDate  string  `json:"release_date"`
	Runtime      int     `json:"runtime"`
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

func FetchMetadata(tmdbID string) (*TMDBResponse, error) {
	apiKey := os.Getenv("TMDB_API_KEY")
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?api_key=%s&language=id-ID&append_to_response=images", tmdbID, apiKey)

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
