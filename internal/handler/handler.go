package handler

import (
	"github.com/crazydw4rf/book-stock-manager/internal/controller"
	"github.com/gofiber/fiber/v2"
)

const (
	BASE_PATH   = "/api/v1"
	BOOK_CREATE = BASE_PATH + "/books"
)

func SetupBookHandler(app *fiber.App, ctrl *controller.BookController) *fiber.App {
	app.Post(BOOK_CREATE, ctrl.Create)

	return app
}
