package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrConnection      = errors.New("connection failed")
	ErrNotFound        = errors.New("not found")
	ErrInvalidResponse = errors.New("invalid response")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrNotNumber       = errors.New("not a number")
)

type trendingResponse struct {
	Page         int        `json:"page"`
	Results      []trending `json:"results"`
	TotalPages   int        `json:"total_pages"`
	TotalResults int        `json:"total_results"`
}

type trending struct {
	ID            int     `json:"id"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	VoteAverage   float64 `json:"vote_average"`
}

func getTrending(apiRoot, apiKey string) ([]trending, error) {
	u := fmt.Sprintf("%s/3/trending/movie/week?api_key=%s", apiRoot, apiKey)
	r, err := newClient().Get(u)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}

	if r.StatusCode != http.StatusOK {
		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		if r.StatusCode == http.StatusUnauthorized {
			err = ErrUnauthorized
		}

		var errResponse struct {
			StatusMessage string `json:"status_message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&errResponse); err != nil {
			return nil, fmt.Errorf("cannot read body: %w", err)
		}

		return nil, fmt.Errorf("%w: %s", err, errResponse.StatusMessage)
	}

	var response trendingResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("cannot read body: %w", err)
	}

	return response.Results, nil
}

type movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Tagline  string `json:"tagline"`
	Overview string `json:"overview"`
}

func getMovie(apiRoot, apiKey string, ID int) (movie, error) {
	u := fmt.Sprintf("%s/3/movie/%d?api_key=%s", apiRoot, ID, apiKey)
	r, err := newClient().Get(u)
	if err != nil {
		return movie{}, fmt.Errorf("%w: %s", ErrConnection, err)
	}

	if r.StatusCode != http.StatusOK {
		return movie{}, handleError(r)
	}

	var m movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		return movie{}, fmt.Errorf("cannot read body: %w", err)
	}

	return m, nil
}

func handleError(r *http.Response) error {
	var err = ErrInvalidResponse
	if r.StatusCode == http.StatusNotFound {
		err = ErrNotFound
	}
	if r.StatusCode == http.StatusUnauthorized {
		err = ErrUnauthorized
	}

	var errResponse struct {
		StatusMessage string `json:"status_message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&errResponse); err != nil {
		return fmt.Errorf("cannot read body: %w", err)
	}

	return fmt.Errorf("%w: %s", err, errResponse.StatusMessage)
}

func newClient() *http.Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	return c
}
