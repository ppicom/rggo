//go:build integration
// +build integration

package cmd

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
	var apiRoot, apiKey string

	if os.Getenv("TMDB_API_KEY") == "" {
		t.Skip("TMDB_API_KEY not set, skipping integration test")
	}
	apiKey = os.Getenv("TMDB_API_KEY")

	apiRoot = "https://api.themoviedb.org"
	if os.Getenv("TMBD_API_ROOT") != "" {
		apiRoot = os.Getenv("TMBD_API_ROOT")
	}

	title := "Avatar: The Way of Water"
	ID := ""

	trendingPassed := t.Run("Trending", func(t *testing.T) {
		var out bytes.Buffer

		err := listTrendingAction(&out, apiRoot, apiKey)

		if err != nil {
			t.Fatal(err)
		}

		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			ln := scanner.Text()
			if strings.Contains(ln, title) {
				ID = strings.Fields(ln)[len(strings.Fields(ln))-1]
			}
		}

		if ID == "" {
			t.Errorf("movie %q not found in trending list", title)
		}
	})

	if !trendingPassed {
		t.Fatal("List trending failed. Stopping integration tests.")
	}

	t.Run("Movie", func(t *testing.T) {
		var out bytes.Buffer

		err := movieAction(&out, apiRoot, apiKey, ID)
		if err != nil {
			t.Fatal(err)
		}

		scanner := bufio.NewScanner(&out)
		var titlePresent bool
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), title) {
				titlePresent = true
				break
			}
		}

		if !titlePresent {
			t.Errorf("Expected title %q in movie details", title)
		}
	})
}
