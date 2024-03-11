package main

import (
	"log"
	"net/http"
)

func main() {
	// store := NewInMemoryPlayerStore()
	pgStore := NewPgPlayerStore("localhost", "testdb", "testusr", "testpwd")
	server := NewPlayerServer(pgStore)
	log.Fatal(http.ListenAndServe(":5000", server))
}
