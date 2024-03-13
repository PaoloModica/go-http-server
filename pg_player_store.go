package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewPgPlayerStore(baseUrl, db, user, pwd string) *PgPlayerStore {
	return &PgPlayerStore{baseUrl, db, user, pwd}
}

type PgPlayerStore struct {
	baseUrl string
	db      string
	user    string
	pwd     string
}

func (p *PgPlayerStore) connect() *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", p.user, p.pwd, p.baseUrl, p.db)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("An error occurred while trying to connect to %s database. Error: %s", p.db, err.Error())
	}
	return db
}

func (p *PgPlayerStore) RecordWin(name string) {
	db := p.connect()
	defer db.Close()
	row := db.QueryRow(`INSERT INTO players as p (name, wins) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET wins = p.wins + 1 RETURNING *;`, name, 1)
	if row.Err() != nil {
		log.Fatalf("An error occurred while inserting data into database. Error: %s", row.Err().Error())
	}
}

func (p *PgPlayerStore) GetPlayerScore(name string) int {
	var score int

	db := p.connect()
	defer db.Close()

	row := db.QueryRow(`SELECT p.wins FROM players p WHERE p.name = $1`, name)
	if row.Err() != nil {
		log.Fatalf("An error occurred while retrieving data from database. Error: %s", row.Err().Error())
	}

	row.Scan(&score)
	return score
}

func (p *PgPlayerStore) GetLeague() League {
	var league []Player

	db := p.connect()
	defer db.Close()

	rows, err := db.Query(`SELECT p.name, p.wins FROM players p ORDER BY wins`)
	if err != nil {
		log.Fatalf("An error occurred while retrieving data from database. Error: %s", err.Error())
	}
	for rows.Next() {
		var name string
		var wins int
		err := rows.Scan(&name, &wins)

		if err != nil {
			log.Fatalf("An error occurred while parsing data retrieved from DB. Error: %s", err.Error())
			break
		}

		league = append(league, Player{name, wins})
	}
	if err := rows.Close(); err != nil {
		log.Fatalf("An error occurred while closing rows: %s", err.Error())
	}

	return league
}
