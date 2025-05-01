package entity

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	BookId      uuid.UUID `json:"book_id" db:"book_id"`
	ISBN        string    `json:"isbn" db:"isbn"`
	Title       string    `json:"title" db:"title"`
	Author      string    `json:"author" db:"author"`
	Publisher   string    `json:"publisher" db:"publisher"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	Stock       int64     `json:"stock" db:"stock"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
