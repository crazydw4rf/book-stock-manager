package usecase

import (
	"context"

	"github.com/crazydw4rf/book-stock-manager/internal/entity"
	"github.com/crazydw4rf/book-stock-manager/internal/model"
	"github.com/crazydw4rf/book-stock-manager/internal/repository"
	"github.com/crazydw4rf/book-stock-manager/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
)

type BookUsecase struct {
	bookRepo  *repository.BookRepository
	validator *validator.Validate
}

func NewBookUsecase(bookRepo *repository.BookRepository, validator *validator.Validate) *BookUsecase {
	return &BookUsecase{bookRepo, validator}
}

func (b BookUsecase) Create(ctx context.Context, bookReq *model.CreateBookRequest) (model.BookResponse, error) {
	err := b.validator.Struct(bookReq)
	if err != nil {
		return model.BookResponse{}, eris.Wrap(fiber.NewError(fiber.StatusBadRequest, "Invalid request payload"), err.Error())
	}

	bookId, err := uuid.NewV7()
	if err != nil {
		return model.BookResponse{}, eris.Errorf("Failed to generate book ID: %v", err)
	}

	book := &entity.Book{
		BookId:      bookId,
		ISBN:        bookReq.ISBN,
		Title:       bookReq.Title,
		Author:      bookReq.Author,
		Publisher:   bookReq.Publisher,
		PublishedAt: bookReq.PublishedAt,
		Stock:       bookReq.Stock,
	}

	book, err = b.bookRepo.Create(ctx, book)
	if err != nil {
		return model.BookResponse{}, eris.Wrap(fiber.NewError(fiber.StatusInternalServerError, "Failed to create book"), eris.ToString(err, true))
	}

	return model.BookToResponse(book), nil
}

func (b BookUsecase) GetById(ctx context.Context, bookId string) (model.BookResponse, error) {
	id, err := uuid.Parse(bookId)
	if err != nil {
		return model.BookResponse{}, eris.Wrap(fiber.NewError(fiber.StatusBadRequest, "Invalid book ID"), err.Error())
	}

	book, err := b.bookRepo.GetById(ctx, id)
	if err != nil {
		if eris.Is(err, types.ErrNoRows) {
			return model.BookResponse{}, fiber.NewError(fiber.StatusNotFound, "Book not found")
		}

		return model.BookResponse{}, eris.Wrap(fiber.NewError(fiber.StatusInternalServerError, "Failed to get book"), eris.ToString(err, true))
	}

	return model.BookToResponse(book), nil
}

func (b BookUsecase) GetByISBN(ctx context.Context, isbn string) (model.BookResponse, error) {
	err := b.validator.VarCtx(ctx, isbn, "isbn")
	if err != nil {
		return model.BookResponse{}, fiber.NewError(fiber.StatusBadRequest, "Invalid ISBN format")
	}

	book, err := b.bookRepo.GetByISBN(ctx, isbn)
	if err != nil {
		if eris.Is(err, types.ErrNoRows) {
			return model.BookResponse{}, fiber.NewError(fiber.StatusNotFound, "Book not found")
		}

		return model.BookResponse{}, eris.Wrap(fiber.NewError(fiber.StatusInternalServerError, "Failed to get book"), eris.ToString(err, true))
	}

	return model.BookToResponse(book), nil
}

func (b BookUsecase) GetMany(ctx context.Context, offset int64, limit int64) ([]model.BookResponse, int64, error) {
	if limit <= 0 {
		return nil, 0, eris.Wrap(fiber.NewError(fiber.StatusBadRequest, "Limit must be greater than 0"), "Invalid limit")
	}

	books, err := b.bookRepo.GetMany(ctx, offset, limit)
	if err != nil {
		if eris.Is(err, types.ErrNoRows) {
			return nil, 0, fiber.NewError(fiber.StatusNotFound, "No books found")
		}
		return nil, 0, eris.Wrap(fiber.NewError(fiber.StatusInternalServerError, "Failed to get books"), eris.ToString(err, true))
	}

	// Get total count for pagination
	total, err := b.bookRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, 0, eris.Wrap(fiber.NewError(fiber.StatusInternalServerError, "Failed to get total count"), eris.ToString(err, true))
	}

	booksResp := make([]model.BookResponse, len(books))
	for i, book := range books {
		booksResp[i] = model.BookToResponse(book)
	}

	return booksResp, total, nil
}

func (b BookUsecase) Update(ctx context.Context, request *model.UpdateBookRequest) (model.BookResponse, error) {
	err := b.validator.Struct(request)
	if err != nil {
		return model.BookResponse{}, eris.Wrap(fiber.NewError(fiber.StatusBadRequest, "Invalid request payload"), err.Error())
	}

	book := &entity.Book{
		BookId:      request.BookID,
		ISBN:        request.ISBN,
		Title:       request.Title,
		Author:      request.Author,
		Publisher:   request.Publisher,
		PublishedAt: request.PublishedAt,
		Stock:       request.Stock,
	}

	updatedBook, err := b.bookRepo.Update(ctx, book)
	if err != nil {
		if eris.Is(err, types.ErrNoRows) {
			return model.BookResponse{}, fiber.NewError(fiber.StatusNotFound, "Book not found")
		}

		return model.BookResponse{}, eris.Wrap(fiber.NewError(fiber.StatusInternalServerError, "Failed to update book"), eris.ToString(err, true))
	}

	return model.BookToResponse(updatedBook), nil
}

func (b BookUsecase) Delete(ctx context.Context, bookId string) error {
	id, err := uuid.Parse(bookId)
	if err != nil {
		return eris.Wrap(fiber.NewError(fiber.StatusBadRequest, "Invalid book ID"), err.Error())
	}

	err = b.bookRepo.Delete(ctx, id)
	if err != nil {
		if eris.Is(err, types.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "Book not found")
		}

		return eris.Wrap(fiber.NewError(fiber.StatusInternalServerError, "Failed to delete book"), err.Error())
	}

	return nil
}
