package main

import (
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	// pgStore := NewPgPlayerStore("localhost", "testdb", "testusr", "testpwd")
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s. Error: %v", dbFileName, err)
	}
	store, err := NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("an error occurred while creating file system player store, error: %v", err)
	}
	server := NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen to port 5000. Error: %v", err)
	}
}
