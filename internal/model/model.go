package model

import "github.com/crazydw4rf/book-stock-manager/internal/entity"

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
