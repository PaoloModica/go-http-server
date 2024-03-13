package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrieveThem(t *testing.T) {
	type testCase struct {
		server *PlayerServer
		id     string
	}

	pgStorePlayerServer := createPlayerServerPgStore(t)
	fsStorePlayerServer, cleanDB := createPlayerServerFSStore(t)

	defer cleanDB()

	testCases := []testCase{
		{pgStorePlayerServer, "postgres store"},
		{fsStorePlayerServer, "file system store"},
	}

	player := "Pepper"
	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("get score %s", testCase.id),
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
			fmt.Sprintf("get league %s", testCase.id),
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

func createPlayerServerPgStore(t testing.TB) *PlayerServer {
	store := NewPgPlayerStore("localhost", "testdb", "testusr", "testpwd")
	return NewPlayerServer(store)
}

func createPlayerServerFSStore(t testing.TB) (*PlayerServer, func()) {
	database, cleanDatabase := createTempFile(t, `[]`)
	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)
	return NewPlayerServer(store), cleanDatabase
}
