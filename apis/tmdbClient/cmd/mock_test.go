package cmd

import (
	"net/http"
	"net/http/httptest"
)

var testResp = map[string]struct {
	Status      int
	Body        string
	ContentType string
}{
	"unauthorized": {
		Status:      http.StatusUnauthorized,
		ContentType: "application/json",
		Body: `{
			"status_code": 7,
			"status_message": "Invalid API key: You must be granted a valid key.",
			"success": false
		}`,
	},
	"trending": {
		Status:      http.StatusOK,
		ContentType: "application/json",
		Body: `{
			"page": 1,
			"results": [
				{
					"adult": false,
					"backdrop_path": "/ovM06PdF3M8wvKb06i4sjW3xoww.jpg",
					"id": 76600,
					"title": "Avatar: The Way of Water",
					"original_language": "en",
					"original_title": "Avatar: The Way of Water",
					"overview": "Set more than a decade after the events of the first film, learn the story of the Sully family (Jake, Neytiri, and their kids), the trouble that follows them, the lengths they go to keep each other safe, the battles they fight to stay alive, and the tragedies they endure.",
					"poster_path": "/t6HIqrRAclMCA60NsSmeqe9RmNV.jpg",
					"media_type": "movie",
					"genre_ids": [
						878,
						12,
						28
					],
					"popularity": 10224.28,
					"release_date": "2022-12-14",
					"video": false,
					"vote_average": 7.743,
					"vote_count": 6308
				},
				{
					"adult": false,
					"backdrop_path": "/i8dshLvq4LE3s0v8PrkDdUyb1ae.jpg",
					"id": 603692,
					"title": "John Wick: Chapter 4",
					"original_language": "en",
					"original_title": "John Wick: Chapter 4",
					"overview": "With the price on his head ever increasing, John Wick uncovers a path to defeating The High Table. But before he can earn his freedom, Wick must face off against a new enemy with powerful alliances across the globe and forces that turn old friends into foes.",
					"poster_path": "/vZloFAK7NmvMGKE7VkF5UHaz0I.jpg",
					"media_type": "movie",
					"genre_ids": [
						28,
						53,
						80
					],
					"popularity": 2569.508,
					"release_date": "2023-03-22",
					"video": false,
					"vote_average": 8.162,
					"vote_count": 621
				}
			],
			"total_pages": 1000,
			"total_results": 20000
		}`,
	},
	"invalidURL": {
		Status:      http.StatusNotFound,
		ContentType: "plain/text",
		Body: `<html>
		<head><title>404 Not Found</title></head>
		<body>
		<center><h1>404 Not Found</h1></center>
		<hr><center>openresty</center>
		</body>
		</html>
		`,
	},
	"movie": {
		Status:      http.StatusOK,
		ContentType: "application/json",
		Body: `{
			"adult": false,
			"backdrop_path": "/i8dshLvq4LE3s0v8PrkDdUyb1ae.jpg",
			"belongs_to_collection": {
				"id": 404609,
				"name": "John Wick Collection",
				"poster_path": "/xUidyvYFsbbuExifLkslpcd8SMc.jpg",
				"backdrop_path": "/fSwYa5q2xRkBoOOjueLpkLf3N1m.jpg"
			},
			"budget": 90000000,
			"genres": [
				{
					"id": 28,
					"name": "Action"
				},
				{
					"id": 53,
					"name": "Thriller"
				},
				{
					"id": 80,
					"name": "Crime"
				}
			],
			"homepage": "https://johnwick.movie",
			"id": 603692,
			"imdb_id": "tt10366206",
			"original_language": "en",
			"original_title": "John Wick: Chapter 4",
			"overview": "With the price on his head ever increasing, John Wick uncovers a path to defeating The High Table. But before he can earn his freedom, Wick must face off against a new enemy with powerful alliances across the globe and forces that turn old friends into foes.",
			"popularity": 2569.508,
			"poster_path": "/vZloFAK7NmvMGKE7VkF5UHaz0I.jpg",
			"production_companies": [
				{
					"id": 3528,
					"logo_path": "/cCzCClIzIh81Fa79hpW5nXoUsHK.png",
					"name": "Thunder Road",
					"origin_country": "US"
				},
				{
					"id": 23008,
					"logo_path": "/5SarYupipdiejsEqUkwu1SpYfru.png",
					"name": "87Eleven",
					"origin_country": "US"
				},
				{
					"id": 491,
					"logo_path": "/rUp0lLKa1pr4UsPm8fgzmnNGxtq.png",
					"name": "Summit Entertainment",
					"origin_country": "US"
				}
			],
			"production_countries": [
				{
					"iso_3166_1": "US",
					"name": "United States of America"
				}
			],
			"release_date": "2023-03-22",
			"revenue": 157000000,
			"runtime": 169,
			"spoken_languages": [
				{
					"english_name": "English",
					"iso_639_1": "en",
					"name": "English"
				},
				{
					"english_name": "Russian",
					"iso_639_1": "ru",
					"name": "Pусский"
				},
				{
					"english_name": "French",
					"iso_639_1": "fr",
					"name": "Français"
				},
				{
					"english_name": "Latin",
					"iso_639_1": "la",
					"name": "Latin"
				},
				{
					"english_name": "Japanese",
					"iso_639_1": "ja",
					"name": "日本語"
				},
				{
					"english_name": "Cantonese",
					"iso_639_1": "cn",
					"name": "广州话 / 廣州話"
				},
				{
					"english_name": "Spanish",
					"iso_639_1": "es",
					"name": "Español"
				}
			],
			"status": "Released",
			"tagline": "No way back. One way out.",
			"title": "John Wick: Chapter 4",
			"video": false,
			"vote_average": 8.158,
			"vote_count": 619
		}`,
	},
	"notFound": {
		Status:      http.StatusNotFound,
		ContentType: "application/json",
		Body: `{
			"success": false,
			"status_code": 34,
			"status_message": "The resource you requested could not be found."
		}`,
	},
}

func mockServer(h http.HandlerFunc) (string, func()) {
	s := httptest.NewServer(h)

	return s.URL, func() {
		s.Close()
	}
}
