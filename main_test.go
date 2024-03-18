package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBooksOk(testing *testing.T) {
	router := BuildRouter()
	request, _ := http.NewRequest("GET", "/books", nil)
	responseWriter := httptest.NewRecorder()
	router.ServeHTTP(responseWriter, request)

	assert.Equal(testing, 200, responseWriter.Code)
	assert.Containsf(testing, responseWriter.Body.String(), "id", "1", "title", "The Alchemist", "author", "Paulo Coelho", "quantity", "10")
}

func TestCreateBookOk(testing *testing.T) {
	router := BuildRouter()
	responseWriter := httptest.NewRecorder()
	responseWriter.Header().Add("Content-Type", "application/json")
	var body = []byte(`{
		"title": "Brief history of time",
		"author": "Sthepen Hawking",
		"quantity": 1
	}`)
	request, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	router.ServeHTTP(responseWriter, request)

	assert.Equal(testing, 201, responseWriter.Code)
	assert.Containsf(testing, responseWriter.Body.String(), "id", "2", "title", "Brief history of time", "author", "Sthepen Hawking", "quantity", "1")
}

func TestCreateBookFail(testing *testing.T) {
	router := BuildRouter()
	responseWriter := httptest.NewRecorder()
	responseWriter.Header().Add("Content-Type", "application/json")
	var body = []byte(`{
		"title": "Brief history of time",
		"author": "Sthepen Hawking",
		"quantity": "1"
	}`)
	request, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	router.ServeHTTP(responseWriter, request)

	assert.Equal(testing, 400, responseWriter.Code)
	assert.Contains(testing, responseWriter.Body.String(), "Failed to parse request body.")
}
