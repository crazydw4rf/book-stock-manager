package main

import (
	_ "github.com/crazydw4rf/book-stock-manager/docs"
	"github.com/crazydw4rf/book-stock-manager/internal/controller"
	"github.com/crazydw4rf/book-stock-manager/internal/handler"
	"github.com/crazydw4rf/book-stock-manager/internal/repository"
	"github.com/crazydw4rf/book-stock-manager/internal/usecase"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(newConfig, newFiberApp, newDBConn, newValidator),
		fx.Provide(repository.NewBookRepository, usecase.NewBookUsecase),
		fx.Provide(controller.NewBookController),
		fx.Decorate(handler.SetupBookHandler),
		fx.Invoke(startApp),
	)

	app.Run()
}
