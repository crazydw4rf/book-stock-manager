package repository

import (
	"context"
	"database/sql"

	"github.com/crazydw4rf/book-stock-manager/internal/entity"
	"github.com/crazydw4rf/book-stock-manager/internal/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rotisserie/eris"
)

type BookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) *BookRepository {
	return &BookRepository{db}
}

func (b BookRepository) Create(ctx context.Context, book *entity.Book) (*entity.Book, error) {
	err := b.db.QueryRowxContext(ctx, bookCreate, book.BookId, book.ISBN, book.Title, book.Author, book.Publisher, book.PublishedAt, book.Stock).StructScan(book)
	if err != nil {
		return nil, eris.Wrap(types.ErrDatabaseQuery, err.Error())
	}

	return book, nil
}

func (b BookRepository) GetById(ctx context.Context, bookId uuid.UUID) (*entity.Book, error) {
	book := new(entity.Book)
	err := b.db.GetContext(ctx, book, bookGetById, bookId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, eris.Wrap(types.ErrNoRows, "book not found")
		}

		return nil, eris.Wrap(types.ErrDatabaseQuery, err.Error())
	}

	return book, nil
}

func (b BookRepository) GetByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	book := new(entity.Book)
	err := b.db.GetContext(ctx, book, bookGetByISBN, isbn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, eris.Wrap(types.ErrNoRows, "book not found")
		}
		return nil, eris.Wrap(types.ErrDatabaseQuery, err.Error())
	}

	return book, nil
}

func (b BookRepository) GetMany(ctx context.Context, offset int64, limit int64) ([]*entity.Book, error) {
	books := make([]*entity.Book, limit)
	err := b.db.SelectContext(ctx, books, bookGetBooksMany, offset, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, eris.Wrap(types.ErrNoRows, "books not found")
		}

		return nil, eris.Wrap(types.ErrDatabaseQuery, err.Error())
	}

	return books, nil
}

func (b BookRepository) Update(ctx context.Context, book *entity.Book) (*entity.Book, error) {
	err := b.db.QueryRowxContext(
		ctx, bookUpdate,
		book.BookId,
		book.ISBN,
		book.Title,
		book.Author,
		book.Publisher,
		book.PublishedAt,
		book.Stock,
	).StructScan(book)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, eris.Wrap(types.ErrNoRows, "book not found")
		}

		return nil, eris.Wrap(types.ErrDatabaseQuery, err.Error())
	}

	return book, nil
}

func (b BookRepository) Delete(ctx context.Context, bookId uuid.UUID) error {
	result, err := b.db.ExecContext(ctx, bookDelete, bookId)
	if err != nil {
		return eris.Wrap(types.ErrDatabaseQuery, err.Error())
	}

	if i, _ := result.RowsAffected(); i <= 0 {
		return eris.Wrap(types.ErrNoRows, "book not found")
	}

	return nil
}
