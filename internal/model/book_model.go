package model

import (
	"github.com/google/uuid"
)

type CreateBookRequest struct {
	ISBN        string `json:"isbn" validate:"required,isbn"`
	Title       string `json:"title" validate:"required"`
	Author      string `json:"author" validate:"required"`
	Publisher   string `json:"publisher" validate:"required"`
	PublishedAt string `json:"published_at" validate:"required"`
	Stock       int64  `json:"stock" validate:"required"`
}

type UpdateBookRequest struct {
	BookID uuid.UUID `json:"book_id" validate:"required"`
	CreateBookRequest
}
