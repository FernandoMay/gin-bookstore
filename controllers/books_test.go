package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rahmanfadhil/gin-bookstore/models"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}
	db.AutoMigrate(&models.Book{})
	models.DB = db
	os.Exit(m.Run())
}

func setupRouter() *gin.Engine {
	r := gin.New()
	r.GET("/books", FindBooks)
	r.GET("/books/:id", FindBook)
	r.POST("/books", CreateBook)
	r.PATCH("/books/:id", UpdateBook)
	r.DELETE("/books/:id", DeleteBook)
	return r
}

func cleanup() {
	models.DB.Exec("DELETE FROM books")
}

func TestCreateBook(t *testing.T) {
	cleanup()
	r := setupRouter()
	w := httptest.NewRecorder()
	body := `{"title":"Test Book","author":"Test Author"}`
	req, _ := http.NewRequest("POST", "/books", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["title"] != "Test Book" {
		t.Errorf("expected title 'Test Book', got %v", data["title"])
	}
	if data["author"] != "Test Author" {
		t.Errorf("expected author 'Test Author', got %v", data["author"])
	}
	if data["id"] == nil {
		t.Error("expected id to be set")
	}
}

func TestCreateBook_Invalid(t *testing.T) {
	cleanup()
	r := setupRouter()
	w := httptest.NewRecorder()
	body := `{"title":"","author":""}`
	req, _ := http.NewRequest("POST", "/books", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestFindBooks_Empty(t *testing.T) {
	cleanup()
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].([]interface{})
	if len(data) != 0 {
		t.Errorf("expected empty list, got %d items", len(data))
	}
}

func TestFindBooks(t *testing.T) {
	cleanup()
	models.DB.Create(&models.Book{Title: "Book1", Author: "Author1"})
	models.DB.Create(&models.Book{Title: "Book2", Author: "Author2"})

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("expected 2 books, got %d", len(data))
	}
}

func TestFindBook(t *testing.T) {
	cleanup()
	book := models.Book{Title: "FindMe", Author: "Finder"}
	models.DB.Create(&book)

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/books/%d", book.ID), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["title"] != "FindMe" {
		t.Errorf("expected title 'FindMe', got %v", data["title"])
	}
}

func TestFindBook_NotFound(t *testing.T) {
	cleanup()
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateBook(t *testing.T) {
	cleanup()
	book := models.Book{Title: "Original", Author: "Author"}
	models.DB.Create(&book)

	r := setupRouter()
	w := httptest.NewRecorder()
	body := `{"title":"Updated Title"}`
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/books/%d", book.ID), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["title"] != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %v", data["title"])
	}
}

func TestUpdateBook_NotFound(t *testing.T) {
	cleanup()
	r := setupRouter()
	w := httptest.NewRecorder()
	body := `{"title":"Nope"}`
	req, _ := http.NewRequest("PATCH", "/books/999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDeleteBook(t *testing.T) {
	cleanup()
	book := models.Book{Title: "DeleteMe", Author: "Author"}
	models.DB.Create(&book)

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/books/%d", book.ID), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["data"] != true {
		t.Errorf("expected true, got %v", resp["data"])
	}

	count := 0
	models.DB.Model(&models.Book{}).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 books after delete, got %d", count)
	}
}

func TestDeleteBook_NotFound(t *testing.T) {
	cleanup()
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/books/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
