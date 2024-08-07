package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"video-hosting/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New" // operation
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS videos(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url VARCHAR(255) UNIQUE,
    	path VARCHAR(255) UNIQUE,
    	name VARCHAR(255),
		author VARCHAR(255) NOT NULL,
		likes INTEGER,
		dislikes INTEGER);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveVideo(url string, videoName string, author string) (int64, error) {
	const op = "storage.sqlite.SaveVideo"
	stmt, err := s.db.Prepare("INSERT INTO videos(url, name, author, path) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(url, videoName, author, url)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetPath(url string) (string, error) {
	const op = "storage.sqlite.GetPath"
	stmt, err := s.db.Prepare("SELECT path FROM videos WHERE url = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var path string
	err = stmt.QueryRow(url).Scan(&path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return path, nil
}