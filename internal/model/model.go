package model

import (
	"strings"

	"github.com/crazydw4rf/book-stock-manager/internal/entity"
)

type DataResponse[T any] struct {
	Data T `json:"data"`
}

// PaginationRequest merepresentasikan parameter kueri untuk paginasi
type PaginationRequest struct {
	Offset int64 `json:"offset" query:"offset" validate:"min=0" example:"0"`
	Limit  int64 `json:"limit" query:"limit" validate:"min=1,max=100" example:"10"`
}

// PaginationMeta merepresentasikan metadata tentang hasil paginasi
type PaginationMeta struct {
	Offset int64 `json:"offset" example:"0"`
	Limit  int64 `json:"limit" example:"10"`
	Total  int64 `json:"total" example:"100"`
}

type PaginatedResponse[T any] struct {
	Data  []T             `json:"data"`
	Meta  PaginationMeta  `json:"meta"`
	Links PaginationLinks `json:"links"`
}

// PaginationLinks menyediakan tautan navigasi yang sesuai dengan HATEOAS untuk paginasi
// Tautan-tautan ini memungkinkan klien untuk menavigasi koleksi tanpa membangun URL secara manual
type PaginationLinks struct {
	Self  string `json:"self" example:"/api/v1/books?offset=0&limit=10"`
	Next  string `json:"next,omitempty" example:"/api/v1/books?offset=10&limit=10"`
	Prev  string `json:"prev,omitempty" example:"/api/v1/books?offset=0&limit=10"`
	First string `json:"first" example:"/api/v1/books?offset=0&limit=10"`
	Last  string `json:"last" example:"/api/v1/books?offset=90&limit=10"`
}

// BookToResponse mengkonversi entity.Book menjadi model BookResponse
// Fungsi ini menangani semua transformasi data yang diperlukan seperti memotong spasi dari ISBN
func BookToResponse(book *entity.Book) BookResponse {
	book.ISBN = strings.TrimSpace(book.ISBN)

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
