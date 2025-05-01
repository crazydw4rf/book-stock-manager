package controller

import (
	"log"

	"github.com/crazydw4rf/book-stock-manager/internal/entity"
	"github.com/crazydw4rf/book-stock-manager/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BookController struct {
	db *sqlx.DB
	v  *validator.Validate
}

func NewBookController(db *sqlx.DB, v *validator.Validate) *BookController {
	return &BookController{db, v}
}

// Create adalah fungsi untuk membuat buku baru
//
//	@Summary		Create a new book
//	@Description	Create a new book
//	@Tags			books
//	@Router			/books [post]
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		model.CreateBookRequest			true	"Book data"
//	@Success		201		{object}	model.WebResponse[entity.Book]	"Book created"
//	@Failure		500		{string}	string							"Internal server error"
//	@Failure		400		{string}	string							"Bad request"
func (b BookController) Create(c *fiber.Ctx) error {
	// parsing request body
	book := new(model.CreateBookRequest)
	if err := c.BodyParser(book); err != nil {
		log.Println("Error parsing request body:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	// validasi payload
	if err := b.v.Struct(book); err != nil {
		log.Println("Validation error:", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid book data")
	}

	// membuat uuid untuk book_id
	bookId, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate book ID")
	}

	// query untuk insert data ke database
	row := b.db.QueryRowx(createBook, bookId, book.ISBN, book.Title, book.Author, book.Publisher, book.PublishedAt, book.Stock)
	if row.Err() != nil {
		log.Println("Database insertion error:", row)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to insert book data")
	}

	book_ := new(entity.Book)
	err = row.StructScan(book_)
	if err != nil {
		log.Println("Error scanning row:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to scan book data")
	}

	return c.Status(fiber.StatusCreated).JSON(model.WebResponse[*entity.Book]{Data: book_})
}

func (b BookController) GetByISBN(c *fiber.Ctx) error {
	isbn := c.Params("isbn")
	if isbn == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ISBN is required")
	}

	// TODO: implement get book by ISBN

	panic("implement me")
}

func (b BookController) GetByID(c *fiber.Ctx) error {
	bookId := c.Params("id")
	if bookId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Book ID is required")
	}

	// TODO: implement get book by ID

	panic("implement me")
}

func (b BookController) GetBooks(c *fiber.Ctx) error {
	// TODO: implement get all books
	panic("implement me")
}

func (b BookController) Update(c *fiber.Ctx) error {
	// TODO: implement update book
	panic("implement me")
}

func (b BookController) Delete(c *fiber.Ctx) error {
	// TODO: implement delete book
	panic("implement me")
}
