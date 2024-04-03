// db.go

package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	Shorten(url string) (string, error)
	Expand(shortLink string) (string, error)
}

type MemoryStorage struct {
	urlMap map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urlMap: make(map[string]string),
	}
}

func (m *MemoryStorage) Shorten(url string) (string, error) {
	shortLink := generateShortLink()
	m.urlMap[shortLink] = url
	return shortLink, nil
}

func (m *MemoryStorage) Expand(shortLink string) (string, error) {
	originalURL, ok := m.urlMap[shortLink]
	if !ok {
		return "", fmt.Errorf("short link not found")
	}
	return originalURL, nil
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() *PostgresStorage {
	// Load environment variables from .env file
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection details from environment variables
	connStr := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return &PostgresStorage{db: db}
}

func (p *PostgresStorage) Shorten(url string) (string, error) {
	shortLink := generateShortLink()
	_, err := p.db.Exec("INSERT INTO links (short, original) VALUES ($1, $2)", shortLink, url)
	if err != nil {
		return "", fmt.Errorf("failed to insert into database: %v", err)
	}
	return shortLink, nil
}

func (p *PostgresStorage) Expand(shortLink string) (string, error) {
	var originalURL string
	err := p.db.QueryRow("SELECT original FROM links WHERE short = $1", shortLink).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("short link not found")
		}
		return "", fmt.Errorf("failed to query database: %v", err)
	}
	return originalURL, nil
}

func generateShortLink() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	const linkLength = 10
	b := make([]byte, linkLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
