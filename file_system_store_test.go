package main

import (
	"os"
	"testing"
)

func TestFileSystemPlayerStore(t *testing.T) {
	t.Run("build store with empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()
		_, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
	})
	t.Run("get league", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name": "Cleo", "Wins": 10},
				{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
		got := store.GetLeague()
		expected := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}
		assertLeague(t, expected, got)
		// read again
		got = store.GetLeague()
		assertLeague(t, expected, got)
	})
	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name": "Cleo", "Wins": 10},
				{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		expected := 33
		got := store.GetPlayerScore("Chris")
		assertScoreEqual(t, expected, got)
	})
	t.Run("store wins for existing player", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name": "Cleo", "Wins": 10},
				{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Chris")
		expected := 34
		got := store.GetPlayerScore("Chris")
		assertScoreEqual(t, expected, got)
	})
	t.Run("store wins for new player", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name": "Cleo", "Wins": 10},
				{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
		store.RecordWin("Matthew")
		expected := 1
		got := store.GetPlayerScore("Matthew")
		assertScoreEqual(t, expected, got)
	})
}

func assertScoreEqual(t testing.TB, expected, got int) {
	t.Helper()
	if expected != got {
		t.Errorf("Expected %d, got %d", expected, got)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("an unexpected error occurred: %v", err)
	}
}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file, error: %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
