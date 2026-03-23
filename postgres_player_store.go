package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// This method make the following assumptions:
//
// - DB_HOST, DB_PORT, DB_USER, DB_PASSWORD & DB_NAME is the value of your Postgres server and the database you want to connect
//
// - In the database, a table called players exist, with name as primary key, varchar(255) and score as integer, check(score >= 0), default 0
//
// - If recreate is true, the table players is refreshed
func NewPostgresPlayerStore(recreate bool) *PostgresPlayerStore {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	ps := &PostgresPlayerStore{db: db}

	if recreate {
		errReset := ps.resetDB()

		if errReset != nil {
			panic(errReset)
		}
	}

	return ps
}

type PostgresPlayerStore struct {
	db *sql.DB
}

func (i *PostgresPlayerStore) GetPlayerScore(name string) int {
	var score int

	err := i.db.QueryRow(
		"SELECT score FROM players WHERE name = $1",
		name,
	).Scan(&score)

	if err == sql.ErrNoRows {
		return 0
	}

	if err != nil {
		panic(err)
	}

	return score
}

func (i *PostgresPlayerStore) RecordWin(name string) {
	result, err := i.db.Exec(
		"UPDATE players SET score = score + 1 WHERE name = $1",
		name,
	)

	if err != nil {
		panic(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	if rowsAffected == 0 {
		_, err := i.db.Exec(
			"INSERT INTO players(name, score) VALUES($1, 1)",
			name,
		)
		if err != nil {
			panic(err)
		}
	}
}

func (i *PostgresPlayerStore) GetLeague() []Player {
	var league []Player

	rows, err := i.db.Query("SELECT name, score FROM players")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Player

		err := rows.Scan(&p.Name, &p.Wins)
		if err != nil {
			panic(err)
		}

		league = append(league, p)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return league
}

// helpers

func (i *PostgresPlayerStore) resetDB() error {
	_, err := i.db.Exec("TRUNCATE TABLE players RESTART IDENTITY")
	return err
}
