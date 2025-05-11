package handler

import (
	"github.com/crazydw4rf/book-stock-manager/internal/controller"
	"github.com/gofiber/fiber/v2"
)

const (
	BASE_PATH      = "/api/v1"
	BOOK_CREATE    = BASE_PATH + "/books"
	BOOK_GETBYID   = BASE_PATH + "/books/:book_id"
	BOOK_GETBYISBN = BASE_PATH + "/books/:isbn"
	BOOK_GETMANY   = BASE_PATH + "/books"
	BOOK_UPDATE    = BASE_PATH + "/books"
	BOOK_DELETE    = BASE_PATH + "/books/:book_id"
)

func SetupBookHandler(app *fiber.App, ctrl *controller.BookController) *fiber.App {
	app.Post(BOOK_CREATE, ctrl.BookCreate)
	app.Get(BOOK_GETBYID, ctrl.GetBookByID)
	app.Get(BOOK_GETBYISBN, ctrl.GetBookByISBN)
	app.Post(BOOK_GETMANY, ctrl.GetBooks)
	app.Patch(BOOK_UPDATE, ctrl.Update)
	app.Delete(BOOK_DELETE, ctrl.Delete)

	return app
}
