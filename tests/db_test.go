// tests/db_test.go

package main

import (
	"Ozon/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	storage := db.NewMemoryStorage()

	shortLink, err := storage.Shorten("http://example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, shortLink)

	originalURL, err := storage.Expand(shortLink)
	assert.NoError(t, err)
	assert.Equal(t, "http://example.com", originalURL)
}

func TestPostgresStorage(t *testing.T) {
	storage := db.NewPostgresStorage()

	shortLink, err := storage.Shorten("http://example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, shortLink)

	originalURL, err := storage.Expand(shortLink)
	assert.NoError(t, err)
	assert.Equal(t, "http://example.com", originalURL)
}
