package controller

import (
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
//	@Router			/books/{offset}/{limit} [get]
//	@Accept			json
//	@Produce		json
//	@Param			offset	path		int					true	"Page offset"
//	@Param			limit	path		int					true	"Page limit"
//	@Success		200		{array}		model.BookResponse	"Books information"
//	@Failure		500		{string}	string				"Failed to get books"
//	@Failure		400		{string}	string				"Invalid query parameters"
func (b BookController) GetBooks(c *fiber.Ctx) error {
	p := struct {
		Offset int `params:"offset"`
		Limit  int `params:"limit"`
	}{}

	err := c.ParamsParser(&p)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	if p.Offset < 0 || p.Limit <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Offset and limit must be greater than 0")
	}

	books, err := b.bookUsecase.GetMany(c.Context(), int64(p.Offset), int64(p.Limit))
	if err != nil {
		var fe *fiber.Error
		if eris.As(err, &fe) {
			return fe
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get books")
	}

	return c.Status(fiber.StatusOK).JSON(books)
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
	bookId := c.Params("id")
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
