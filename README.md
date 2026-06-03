# Bookstore REST API

Gin + GORM REST API for book management with SQLite.

## Endpoints

- `GET /books` — List all books
- `GET /books/:id` — Get a book
- `POST /books` — Create a book
- `PATCH /books/:id` — Update a book
- `DELETE /books/:id` — Delete a book

## Run

```bash
go run main.go
```

## Test

```bash
go test ./...
```

Based on [LogRocket article](https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/).
