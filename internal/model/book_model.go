package model

import (
	"time"

	"github.com/google/uuid"
)

type BookResponse struct {
	BookID      uuid.UUID `json:"book_id" example:"b2a0f3c4-5d8e-4c1b-9f7e-2d3f4e5a6b7c"`
	ISBN        string    `json:"isbn" example:"978-3-16-148410-0"`
	Title       string    `json:"title" example:"Hujan"`
	Author      string    `json:"author" example:"Tere Liye"`
	Publisher   string    `json:"publisher" example:"Gramedia"`
	PublishedAt time.Time `json:"published_at" example:"2016-01-28"`
	Stock       int64     `json:"stock" example:"200"`
}

type CreateBookRequest struct {
	ISBN        string    `json:"isbn" validate:"required,isbn" example:"9783161484100"`
	Title       string    `json:"title" validate:"required" example:"Hujan"`
	Author      string    `json:"author" validate:"required" example:"Tere Liye"`
	Publisher   string    `json:"publisher" validate:"required" example:"Gramedia"`
	PublishedAt time.Time `json:"published_at" validate:"required" example:"2016-01-28"`
	Stock       int64     `json:"stock" validate:"required,gte=0" example:"200"`
}

type UpdateBookRequest struct {
	BookID      uuid.UUID `json:"book_id" validate:"required"`
	ISBN        string    `json:"isbn" validate:"omitempty,isbn" example:"9783161484100"`
	Title       string    `json:"title" validate:"omitempty" example:"Hujan"`
	Author      string    `json:"author" validate:"omitempty" example:"Tere Liye"`
	Publisher   string    `json:"publisher" validate:"omitempty" example:"Gramedia"`
	PublishedAt time.Time `json:"published_at" validate:"omitempty" example:"2016-01-28"`
	Stock       int64     `json:"stock" validate:"omitempty,gte=0" example:"200"`
}
