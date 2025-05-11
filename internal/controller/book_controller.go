package controller

import (
	"fmt"
	"log"

	"github.com/crazydw4rf/book-stock-manager/internal/model"
	"github.com/crazydw4rf/book-stock-manager/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/rotisserie/eris"
)

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
//	@Param			payload	body		model.CreateBookRequest	true	"Request payload"
//	@Success		201		{object}	model.BookResponse		"Book created"
//	@Failure		500		{string}	string					"Failed to create book"
//	@Failure		400		{string}	string					"Invalid request payload"
func (b BookController) BookCreate(c *fiber.Ctx) error {
	// parse request body ke dalam struct CreateBookRequest
	request := new(model.CreateBookRequest)
	if err := c.BodyParser(request); err != nil {
		log.Println("Error parsing request body:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	// panggil usecase untuk membuat buku baru
	bookResp, err := b.bookUsecase.Create(c.Context(), request)
	if err != nil {
		var fe *fiber.Error
		log.Println("Error creating book:", eris.ToString(err, true))

		if eris.As(err, &fe) {
			return fe
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create book")
	}

	return c.Status(fiber.StatusCreated).JSON(bookResp)
}

// GetBookByISBN mengambil data buku berdasarkan ISBN
//
//	@Summary		Get book by ISBN
//	@Description	Get a book information by ISBN
//	@Tags			books
//	@Router			/books/isbn/{isbn} [get]
//	@Accept			json
//	@Produce		json
//	@Param			isbn	path		string				true	"ISBN"
//	@Success		200		{object}	model.BookResponse	"Book information"
//	@Failure		500		{string}	string				"Failed to get book by ISBN"
//	@Failure		404		{string}	string				"Book not found"
//	@Failure		400		{string}	string				"ISBN is required"
func (b BookController) GetBookByISBN(c *fiber.Ctx) error {
	isbn := c.Params("isbn")
	if isbn == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ISBN is required")
	}

	book, err := b.bookUsecase.GetByISBN(c.Context(), isbn)
	if err != nil {
		var fe *fiber.Error
		log.Println("Error getting book by ISBN:", eris.ToString(err, true))
		if eris.As(err, &fe) {
			return fe
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get book by ISBN")
	}

	return c.Status(fiber.StatusOK).JSON(book)
}

// GetBookByID mengambil data buku berdasarkan ID
//
//	@Summary		Get book by ID
//	@Description	Get a book information by ID
//	@Tags			books
//	@Router			/books/{book_id} [get]
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Book ID"
//	@Success		200	{object}	model.BookResponse	"Book information"
//	@Failure		500	{string}	string				"Failed to get book by ID"
//	@Failure		404	{string}	string				"Book not found"
//	@Failure		400	{string}	string				"Book ID is required"
func (b BookController) GetBookByID(c *fiber.Ctx) error {
	bookId := c.Params("book_id")
	if bookId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Book ID is required")
	}

	book, err := b.bookUsecase.GetById(c.Context(), bookId)
	if err != nil {
		var fe *fiber.Error
		if eris.As(err, &fe) {
			return fe
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get book by ID")
	}

	return c.Status(fiber.StatusOK).JSON(book)
}

// GetBooks mengambil daftar buku dengan pagination
//
//	@Summary		Get books with pagination
//	@Description	Get a list of books with pagination
//	@Tags			books
//	@Router			/books [get]
//	@Accept			json
//	@Produce		json
//	@Param			offset	query		int											false	"Page offset (default: 0)"
//	@Param			limit	query		int											false	"Page limit (default: 10)"
//	@Success		200		{object}	model.PaginatedResponse[model.BookResponse]	"Books information with pagination metadata"
//	@Failure		500		{string}	string										"Failed to get books"
//	@Failure		400		{string}	string										"Invalid query parameters"
func (b BookController) GetBooks(c *fiber.Ctx) error {
	// Parse pagination request
	pagination := new(model.PaginationRequest)
	if err := c.QueryParser(pagination); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	// Set default values if not provided
	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}
	if pagination.Offset < 0 {
		pagination.Offset = 0
	}

	// Validate pagination parameters
	if pagination.Limit > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "Maximum limit is 100")
	}

	books, total, err := b.bookUsecase.GetMany(c.Context(), pagination.Offset, pagination.Limit)
	if err != nil {
		var fe *fiber.Error
		if eris.As(err, &fe) {
			return fe
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get books")
	}

	// Calculate if there are previous and next pages
	hasNext := pagination.Offset+pagination.Limit < total
	hasPrev := pagination.Offset > 0

	// Build full request URL
	baseURL := c.BaseURL() + c.Route().Path

	// Create response with pagination metadata
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

	// Add next link if there are more items
	if hasNext {
		response.Links.Next = fmt.Sprintf("%s?offset=%d&limit=%d", baseURL, pagination.Offset+pagination.Limit, pagination.Limit)
	}

	// Add previous link if not on the first page
	if hasPrev {
		response.Links.Prev = fmt.Sprintf("%s?offset=%d&limit=%d", baseURL, max(0, pagination.Offset-pagination.Limit), pagination.Limit)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Update memperbarui data buku
//
//	@Summary		Update book
//	@Description	Update book information partially
//	@Tags			books
//	@Router			/books [patch]
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		model.UpdateBookRequest	true	"Request payload"
//	@Success		200		{object}	model.BookResponse		"Updated book"
//	@Failure		500		{string}	string					"Failed to update book"
//	@Failure		404		{string}	string					"Book not found"
//	@Failure		400		{string}	string					"Invalid request payload"
func (b BookController) Update(c *fiber.Ctx) error {
	request := new(model.UpdateBookRequest)
	err := c.BodyParser(request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	book, err := b.bookUsecase.Update(c.Context(), request)
	if err != nil {
		var fe *fiber.Error
		if eris.As(err, &fe) {
			return fe
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update book")
	}

	return c.Status(fiber.StatusOK).JSON(book)
}

// Delete menghapus data buku berdasarkan ID
//
//	@Summary		Delete book
//	@Description	Delete book by ID
//	@Tags			books
//	@Router			/books/{id} [delete]
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Book ID"
//	@Success		204	{string}	string	"No Content"
//	@Failure		500	{string}	string	"Failed to delete book"
//	@Failure		404	{string}	string	"Book not found"
//	@Failure		400	{string}	string	"Book ID is required"
func (b BookController) Delete(c *fiber.Ctx) error {
	bookId := c.Params("book_id")
	if bookId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Book ID is required")
	}

	err := b.bookUsecase.Delete(c.Context(), bookId)
	if err != nil {
		var fe *fiber.Error
		if eris.As(err, &fe) {
			return fe
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete book")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// max returns the maximum of two int64 values
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
