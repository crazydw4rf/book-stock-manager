package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/crazydw4rf/book-stock-manager/internal/model"
	"github.com/crazydw4rf/book-stock-manager/internal/types"
	"github.com/crazydw4rf/book-stock-manager/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/rotisserie/eris"
)

// newHTTPError creates a new HTTPError response with additional context information
func newHTTPError(c *fiber.Ctx, status int, message string) error {
	now := time.Now().Format(time.RFC3339)
	return c.Status(status).JSON(types.HTTPError{
		Code:      status,
		Message:   message,
		Error:     http.StatusText(status),
		Timestamp: now,
		Path:      c.Path(),
	})
}

type BookController struct {
	bookUsecase *usecase.BookUsecase
}

func NewBookController(bookUsecase *usecase.BookUsecase) *BookController {
	return &BookController{bookUsecase}
}

// BookCreate membuat buku baru
//
//	@Summary		Create a new book
//	@Description	Create a new book with the provided information
//	@Tags			books
//	@Router			/books [post]
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		model.CreateBookRequest					true	"Request payload"
//	@Success		201		{object}	model.DataResponse[model.BookResponse]	"Book created successfully"
//	@Failure		500		{object}	types.HTTPError							"Internal server error"
//	@Failure		400		{object}	types.HTTPError							"Invalid request payload"
func (b BookController) BookCreate(c *fiber.Ctx) error {
	// parse request body ke dalam struct CreateBookRequest
	request := new(model.CreateBookRequest)
	if err := c.BodyParser(request); err != nil {
		log.Println("Error parsing request body:", err)
		return newHTTPError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	// panggil usecase untuk membuat buku baru
	bookResp, err := b.bookUsecase.Create(c.Context(), request)
	if err != nil {
		var fe *fiber.Error
		log.Println("Error creating book:", eris.ToString(err, true))

		if eris.As(err, &fe) {
			return newHTTPError(c, fe.Code, fe.Message)
		}

		return newHTTPError(c, fiber.StatusInternalServerError, "Failed to create book")
	}

	response := model.DataResponse[model.BookResponse]{
		Data: bookResp,
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetBookByISBN mengambil data buku berdasarkan ISBN
//
//	@Summary		Get book by ISBN
//	@Description	Get a book information by ISBN
//	@Tags			books
//	@Router			/books/isbn/{isbn} [get]
//	@Accept			json
//	@Produce		json
//	@Param			isbn	path		string									true	"ISBN"
//	@Success		200		{object}	model.DataResponse[model.BookResponse]	"Book information retrieved successfully"
//	@Failure		500		{object}	types.HTTPError							"Internal server error"
//	@Failure		404		{object}	types.HTTPError							"Book not found"
//	@Failure		400		{object}	types.HTTPError							"Invalid ISBN format or ISBN is required"
func (b BookController) GetBookByISBN(c *fiber.Ctx) error {
	isbn := c.Params("isbn")
	if isbn == "" {
		return newHTTPError(c, fiber.StatusBadRequest, "ISBN is required")
	}

	book, err := b.bookUsecase.GetByISBN(c.Context(), isbn)
	if err != nil {
		var fe *fiber.Error
		log.Println("Error getting book by ISBN:", eris.ToString(err, true))
		if eris.As(err, &fe) {
			return newHTTPError(c, fe.Code, fe.Message)
		}

		return newHTTPError(c, fiber.StatusInternalServerError, "Failed to get book by ISBN")
	}

	response := model.DataResponse[model.BookResponse]{
		Data: book,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// GetBookByID mengambil data buku berdasarkan ID
//
//	@Summary		Get book by ID
//	@Description	Get a book information by ID
//	@Tags			books
//	@Router			/books/{book_id} [get]
//	@Accept			json
//	@Produce		json
//	@Param			book_id	path		string									true	"Book ID"
//	@Success		200		{object}	model.DataResponse[model.BookResponse]	"Book information retrieved successfully"
//	@Failure		500		{object}	types.HTTPError							"Internal server error"
//	@Failure		404		{object}	types.HTTPError							"Book not found"
//	@Failure		400		{object}	types.HTTPError							"Invalid Book ID format or Book ID is required"
func (b BookController) GetBookByID(c *fiber.Ctx) error {
	bookId := c.Params("book_id")
	if bookId == "" {
		return newHTTPError(c, fiber.StatusBadRequest, "Book ID is required")
	}

	book, err := b.bookUsecase.GetById(c.Context(), bookId)
	if err != nil {
		var fe *fiber.Error
		if eris.As(err, &fe) {
			return newHTTPError(c, fe.Code, fe.Message)
		}

		return newHTTPError(c, fiber.StatusInternalServerError, "Failed to get book by ID")
	}

	response := model.DataResponse[model.BookResponse]{
		Data: book,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// GetBooks mengambil daftar buku dengan pagination
//
//	@Summary		Get books with pagination
//	@Description	Get a list of books with pagination support including navigation links
//	@Tags			books
//	@Router			/books [get]
//	@Accept			json
//	@Produce		json
//	@Param			offset	query		int											false	"Page offset (default: 0)"
//	@Param			limit	query		int											false	"Page limit (default: 10, max: 100)"
//	@Success		200		{object}	model.PaginatedResponse[model.BookResponse]	"Books information with pagination metadata and navigation links"
//	@Failure		500		{object}	types.HTTPError								"Internal server error"
//	@Failure		400		{object}	types.HTTPError								"Invalid query parameters"
func (b BookController) GetBooks(c *fiber.Ctx) error {
	pagination := new(model.PaginationRequest)
	if err := c.QueryParser(pagination); err != nil {
		return newHTTPError(c, fiber.StatusBadRequest, "Invalid query parameters")
	}

	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}
	if pagination.Offset < 0 {
		pagination.Offset = 0
	}

	if pagination.Limit > 100 {
		return newHTTPError(c, fiber.StatusBadRequest, "Maximum limit is 100")
	}

	books, total, err := b.bookUsecase.GetMany(c.Context(), pagination.Offset, pagination.Limit)
	if err != nil {
		var fe *fiber.Error
		if eris.As(err, &fe) {
			return newHTTPError(c, fe.Code, fe.Message)
		}

		return newHTTPError(c, fiber.StatusInternalServerError, "Failed to get books")
	}

	hasNext := pagination.Offset+pagination.Limit < total
	hasPrev := pagination.Offset > 0

	baseURL := c.BaseURL() + c.Route().Path
	response := model.PaginatedResponse[model.BookResponse]{
		Data: books,
		Meta: model.PaginationMeta{
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
			Total:  total,
		},
		Links: model.PaginationLinks{
			Self:  fmt.Sprintf("%s?offset=%d&limit=%d", baseURL, pagination.Offset, pagination.Limit),
			First: fmt.Sprintf("%s?offset=0&limit=%d", baseURL, pagination.Limit),
			Last:  fmt.Sprintf("%s?offset=%d&limit=%d", baseURL, (total/pagination.Limit)*pagination.Limit, pagination.Limit),
		},
	}

	if hasNext {
		response.Links.Next = fmt.Sprintf("%s?offset=%d&limit=%d", baseURL, pagination.Offset+pagination.Limit, pagination.Limit)
	}

	if hasPrev {
		response.Links.Prev = fmt.Sprintf("%s?offset=%d&limit=%d", baseURL, max(0, pagination.Offset-pagination.Limit), pagination.Limit)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Update memperbarui data buku
//
//	@Summary		Update book
//	@Description	Update book information partially. For stock field, use -1 as a sentinel value to indicate no update is intended.
//	@Tags			books
//	@Router			/books [patch]
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		model.UpdateBookRequest					true	"Request payload"
//	@Success		200		{object}	model.DataResponse[model.BookResponse]	"Book updated successfully"
//	@Failure		500		{object}	types.HTTPError							"Internal server error"
//	@Failure		404		{object}	types.HTTPError							"Book not found"
//	@Failure		400		{object}	types.HTTPError							"Invalid request payload"
func (b BookController) Update(c *fiber.Ctx) error {
	request := new(model.UpdateBookRequest)
	err := c.BodyParser(request)
	if err != nil {
		log.Println("Error updating book:", eris.ToString(err, true))
		return newHTTPError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	book, err := b.bookUsecase.Update(c.Context(), request)
	if err != nil {
		log.Println("Error updating book:", eris.ToString(err, true))
		var fe *fiber.Error
		if eris.As(err, &fe) {
			return newHTTPError(c, fe.Code, fe.Message)
		}
		return newHTTPError(c, fiber.StatusInternalServerError, "Failed to update book")
	}

	response := model.DataResponse[model.BookResponse]{
		Data: book,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// Delete menghapus data buku berdasarkan ID
//
//	@Summary		Delete book
//	@Description	Delete book by ID
//	@Tags			books
//	@Router			/books/{book_id} [delete]
//	@Accept			json
//	@Produce		json
//	@Param			book_id	path		string			true	"Book ID"
//	@Success		204		{string}	string			"Book deleted successfully"
//	@Failure		500		{object}	types.HTTPError	"Internal server error"
//	@Failure		404		{object}	types.HTTPError	"Book not found"
//	@Failure		400		{object}	types.HTTPError	"Invalid Book ID format or Book ID is required"
func (b BookController) Delete(c *fiber.Ctx) error {
	bookId := c.Params("book_id")
	if bookId == "" {
		return newHTTPError(c, fiber.StatusBadRequest, "Book ID is required")
	}

	err := b.bookUsecase.Delete(c.Context(), bookId)
	if err != nil {
		var fe *fiber.Error
		if eris.As(err, &fe) {
			return newHTTPError(c, fe.Code, fe.Message)
		}

		return newHTTPError(c, fiber.StatusInternalServerError, "Failed to delete book")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
