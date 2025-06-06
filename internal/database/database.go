package database

import (
	"database/sql"
	"errors"
	"fmt"
	"url_short/pkg"

	_ "github.com/lib/pq"
)

var ErrURLNotFound = errors.New("url not found")

type DB struct {
	db *sql.DB
}

func Init(params string) (*DB, error) {
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS url(
		id SERIAL PRIMARY KEY,
		alias TEXT UNIQUE,
		url TEXT NOT NULL,
		created_at 	TIMESTAMP DEFAULT NOW());
	`)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &DB{db: db}, nil
}

func (db *DB) SaveURL(url string) (string, error) {
	existingURL, err := db.GetAlias(url) //checking if url exists

	if err != ErrURLNotFound { // url exists
		return existingURL, nil
	}

	var id uint64
	err = db.db.QueryRow(`INSERT INTO url(url) values($1) RETURNING id`, url).Scan(&id) //inserting url
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	alias := pkg.Encode_base62(uint64(id))
	_, err = db.db.Exec(`UPDATE url SET alias = $1 WHERE id = $2`, alias, id) //inserting alias

	return alias, nil
}

func (db *DB) GetAlias(url string) (string, error) {
	var alias string
	err := db.db.QueryRow(`SELECT alias FROM url WHERE url = $1`, url).Scan(&alias)
	if err == sql.ErrNoRows {
		return "", ErrURLNotFound
	}
	return alias, nil
}

func (db *DB) GetURL(alias string) (string, error) {
	var url string
	err := db.db.QueryRow(`SELECT url FROM url WHERE alias = $1`, alias).Scan(&url)
	if err == sql.ErrNoRows {
		return "", ErrURLNotFound
	}
	return url, nil
}
