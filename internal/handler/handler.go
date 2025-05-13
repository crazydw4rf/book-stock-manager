package handler

import (
	"github.com/crazydw4rf/book-stock-manager/internal/config"
	"github.com/crazydw4rf/book-stock-manager/internal/controller"
	"github.com/gofiber/fiber/v2"
)

const (
	BOOK_CREATE_ROUTE    = config.BASE_API_HTTP_PATH + "/books"
	BOOK_GETBYID_ROUTE   = config.BASE_API_HTTP_PATH + "/books/:book_id"
	BOOK_GETBYISBN_ROUTE = config.BASE_API_HTTP_PATH + "/books/isbn/:isbn"
	BOOK_GETMANY_ROUTE   = config.BASE_API_HTTP_PATH + "/books"
	BOOK_UPDATE_ROUTE    = config.BASE_API_HTTP_PATH + "/books"
	BOOK_DELETE_ROUTE    = config.BASE_API_HTTP_PATH + "/books/:book_id"
)

func SetupBookHandler(app *fiber.App, ctrl *controller.BookController) *fiber.App {
	app.Post(BOOK_CREATE_ROUTE, ctrl.BookCreate)
	app.Get(BOOK_GETBYID_ROUTE, ctrl.GetBookByID)
	app.Get(BOOK_GETBYISBN_ROUTE, ctrl.GetBookByISBN)
	app.Get(BOOK_GETMANY_ROUTE, ctrl.GetBooks)
	app.Patch(BOOK_UPDATE_ROUTE, ctrl.Update)
	app.Delete(BOOK_DELETE_ROUTE, ctrl.Delete)

	return app
}
