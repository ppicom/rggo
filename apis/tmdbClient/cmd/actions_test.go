//go:build !integration
// +build !integration

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestTrendingAction(t *testing.T) {
	testCases := []struct {
		name     string
		expError error
		expOut   string
		resp     struct {
			Status      int
			Body        string
			ContentType string
		}
		closeServer bool
	}{
		{name: "Results",
			expError: nil,
			expOut: `1.            Avatar: The Way of Water              ID: 76600
              Vote average:           7.743000
2.            John Wick: Chapter 4                  ID: 603692
              Vote average:           8.162000
`,
			resp:        testResp["trending"],
			closeServer: false,
		},
		{name: "Unauthorized",
			expError:    ErrUnauthorized,
			expOut:      "",
			resp:        testResp["unauthorized"],
			closeServer: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resp.Status)
				w.Header().Set("Content-Type", tc.resp.ContentType)
				fmt.Fprintln(w, tc.resp.Body)
			})
			defer cleanup()

			if tc.closeServer {
				cleanup()
			}

			var out bytes.Buffer

			err := listTrendingAction(&out, url, "")

			if tc.expError != nil {
				if err == nil {
					t.Fatalf("Expected error %q, got nil", tc.expError)
				}

				if !errors.Is(err, tc.expError) {
					t.Errorf("Expected error %q, got %q", tc.expError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error %q", err)
			}

			if tc.expOut != out.String() {
				t.Errorf("Expected output %q, got %q", tc.expOut, out.String())
			}
		})
	}
}

func TestMovieAction(t *testing.T) {
	testCases := []struct {
		name     string
		expError error
		expOut   string
		arg      string
		resp     struct {
			Status      int
			Body        string
			ContentType string
		}
		closeServer bool
	}{
		{
			name:     "Results",
			expError: nil,
			arg:      "603692",
			expOut: `John Wick: Chapter 4              No way back. One way out.
With the price on his head ever increasing, John Wick uncovers a path to defeating The High Table. But before he can earn his freedom, Wick must face off against a new enemy with powerful alliances across the globe and forces that turn old friends into foes.
`,
			resp:        testResp["movie"],
			closeServer: false,
		},
		{
			name:        "NotFound",
			expError:    ErrNotFound,
			arg:         "123123123123123",
			resp:        testResp["notFound"],
			closeServer: false,
		},
		{
			name:        "InvalidURL",
			expError:    ErrConnection,
			arg:         "603692",
			resp:        testResp["invalidURL"],
			closeServer: false,
		},
		{
			name:        "ServerDown",
			expError:    ErrConnection,
			arg:         "603692",
			resp:        testResp["invalidURL"],
			closeServer: true,
		},
		{
			name:     "InvalidID",
			expError: ErrNotNumber,
			arg:      "a",
			resp:     testResp["notFound"],
		},
	}

	for _, tc := range testCases {
		URL, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.resp.Status)
			w.Header().Set("Content-Type", tc.resp.ContentType)
			fmt.Fprintln(w, tc.resp.Body)
		})
		defer cleanup()

		if tc.closeServer {
			cleanup()
		}

		var out bytes.Buffer

		err := movieAction(&out, URL, "", tc.arg)

		if tc.expError != nil {
			if err == nil {
				t.Fatalf("Expected error %q, got nil", tc.expError)
			}

			if !errors.Is(err, tc.expError) {
				t.Errorf("Expected error %q, got %q", tc.expError, err)
			}
			return
		}

		if err != nil {
			t.Errorf("Unexpected error %q", err)
		}

		if tc.expOut != out.String() {
			t.Errorf("Expected output %q, got %q", tc.expOut, out.String())
		}
	}
}
