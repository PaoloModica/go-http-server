package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrieveThem(t *testing.T) {
	testCases := []struct {
		server *PlayerServer
		id     string
	}{
		{NewPlayerServer(NewInMemoryPlayerStore()), "in-memory player store"},
		{NewPlayerServer(NewPgPlayerStore("localhost", "testdb", "testusr", "testpwd")), "PG player store"},
	}

	player := "Pepper"
	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("Get Score - %s", testCase.id),
			func(t *testing.T) {
				for i := 0; i < 3; i++ {
					testCase.server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
				}

				response := httptest.NewRecorder()
				testCase.server.ServeHTTP(response, newGetPlayerScoreRequest(player))
				assertStatus(t, response.Code, http.StatusOK)
				assertResponseBody(t, response.Body.String(), "3")
			},
		)
		t.Run(
			fmt.Sprintf("Get league - %s", testCase.id),
			func(t *testing.T) {
				response := httptest.NewRecorder()
				testCase.server.ServeHTTP(response, newLeagueRequest())
				league := getLeagueFromResponse(t, response.Body)
				expectedLeague := []Player{
					{"Pepper", 3},
				}
				assertStatus(t, response.Code, http.StatusOK)
				assertLeague(t, expectedLeague, league)
			},
		)
	}
}
