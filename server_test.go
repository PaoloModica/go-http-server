package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   League
}

func (s *StubPlayerStore) GetPlayerScore(player string) int {
	score := s.scores[player]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  30,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)
	t.Run("returns 404 on missing player", func(t *testing.T) {
		request := newGetPlayerScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetPlayerScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		expected_result := "20"
		got := response.Body.String()

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, expected_result)
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetPlayerScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		expected_result := "30"
		got := response.Body.String()

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, expected_result)
	})

	t.Run("it records win", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)
		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin, want %d", len(store.winCalls), 1)
		}
		if store.winCalls[0] != player {
			t.Errorf("Error in recording player win, expected %q, got %q", player, store.winCalls[0])
		}
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("it returns accepted on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table in JSON format", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}
		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		league := getLeagueFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseHeader(t, response, "application/json")
		assertLeague(t, wantedLeague, league)
	})
}

func newGetPlayerScoreRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)
	return request
}

func newPostWinRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Expected %q, got %q", want, got)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Expected %d, got %d", want, got)
	}
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league League) {
	t.Helper()

	league, err := NewLeague(body)

	if err != nil {
		t.Fatalf("Unable to parse response %q into slice of Players %v", body, err)
	}

	return
}

func assertLeague(t testing.TB, expected, got []Player) {
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func assertResponseHeader(t testing.TB, response *httptest.ResponseRecorder, expected string) {
	if response.Result().Header.Get("content-type") != expected {
		t.Errorf("Expected %s, got %v", expected, response.Result().Header)
	}
}
