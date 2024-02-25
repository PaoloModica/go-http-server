package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRewtrieveThem(t *testing.T) {
	t.Run("in-memory player store", func(t *testing.T) {
		store := NewInMemoryPlayerStore()
		server := PlayerServer{store}
		player := "Pepper"

		for i := 0; i < 3; i++ {
			server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
		}

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetPlayerScoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "3")
	})
	t.Run("postegres player store", func(t *testing.T) {
		store := NewPgPlayerStore("localhost", "testdb", "testusr", "testpwd")
		server := PlayerServer{store}
		player := "Pepper"

		for i := 0; i < 3; i++ {
			server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
		}

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetPlayerScoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "3")
	})
}
