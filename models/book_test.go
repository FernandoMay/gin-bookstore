package models

import (
	"encoding/json"
	"testing"
)

func TestBookStruct(t *testing.T) {
	book := Book{Title: "Test Title", Author: "Test Author"}
	if book.Title != "Test Title" {
		t.Errorf("expected 'Test Title', got '%s'", book.Title)
	}
	if book.Author != "Test Author" {
		t.Errorf("expected 'Test Author', got '%s'", book.Author)
	}
	if book.ID != 0 {
		t.Errorf("expected ID 0, got %d", book.ID)
	}
}

func TestBookJSON(t *testing.T) {
	book := Book{ID: 1, Title: "JSON Book", Author: "JSON Author"}
	data, err := json.Marshal(book)
	if err != nil {
		t.Fatal(err)
	}

	var decoded Book
	json.Unmarshal(data, &decoded)
	if decoded.Title != "JSON Book" {
		t.Errorf("expected 'JSON Book', got '%s'", decoded.Title)
	}
	if decoded.Author != "JSON Author" {
		t.Errorf("expected 'JSON Author', got '%s'", decoded.Author)
	}
	if decoded.ID != 1 {
		t.Errorf("expected ID 1, got %d", decoded.ID)
	}
}
