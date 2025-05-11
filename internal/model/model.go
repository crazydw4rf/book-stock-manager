package model

import "github.com/crazydw4rf/book-stock-manager/internal/entity"

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// PaginationRequest represents the query parameters for pagination
type PaginationRequest struct {
	Offset int64 `json:"offset" query:"offset" validate:"min=0" example:"0"`
	Limit  int64 `json:"limit" query:"limit" validate:"min=1,max=100" example:"10"`
}

// PaginationMeta represents metadata about the pagination
type PaginationMeta struct {
	Offset int64 `json:"offset" example:"0"`
	Limit  int64 `json:"limit" example:"10"`
	Total  int64 `json:"total" example:"100"`
}

// PaginatedResponse is a generic wrapper for paginated responses
type PaginatedResponse[T any] struct {
	Data  []T           `json:"data"`
	Meta  PaginationMeta `json:"meta"`
	Links PaginationLinks `json:"links"`
}

// PaginationLinks provides HATEOAS links for pagination
type PaginationLinks struct {
	Self  string `json:"self" example:"/api/v1/books?offset=0&limit=10"`
	Next  string `json:"next,omitempty" example:"/api/v1/books?offset=10&limit=10"`
	Prev  string `json:"prev,omitempty" example:"/api/v1/books?offset=0&limit=10"`
	First string `json:"first" example:"/api/v1/books?offset=0&limit=10"`
	Last  string `json:"last" example:"/api/v1/books?offset=90&limit=10"`
}

func BookToResponse(book *entity.Book) BookResponse {
	return BookResponse{
		BookID:      book.BookId,
		ISBN:        book.ISBN,
		Title:       book.Title,
		Author:      book.Author,
		Publisher:   book.Publisher,
		PublishedAt: book.PublishedAt,
		Stock:       book.Stock,
	}
}
