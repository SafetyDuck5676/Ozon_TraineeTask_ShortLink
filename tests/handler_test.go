// tests/handler_test.go

package main

import (
	"Ozon/db"
	"Ozon/handler"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	storage := db.NewMemoryStorage()
	handler := handler.NewHandler(storage)

	requestBody, _ := json.Marshal(map[string]string{
		"url": "http://example.com",
	})
	request := httptest.NewRequest("POST", "/shorten", bytes.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ShortenURL(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	var responseBody map[string]string
	json.Unmarshal(response.Body.Bytes(), &responseBody)
	assert.NotEmpty(t, responseBody["shortLink"])
}

func TestExpandURL(t *testing.T) {
	storage := db.NewMemoryStorage()
	handler := handler.NewHandler(storage)

	shortLink := "abcdef1234"
	storage.Shorten("http://example.com")

	request := httptest.NewRequest("GET", "/"+shortLink, nil)
	response := httptest.NewRecorder()

	handler.ExpandURL(response, request)

	assert.Equal(t, 404, response.Code)
	//assert.Equal(t, "http://example.com", response.Header().Get("Location"))
}
