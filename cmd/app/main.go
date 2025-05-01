package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/crazydw4rf/book-stock-manager/docs"
	"github.com/crazydw4rf/book-stock-manager/internal/config"
	"github.com/crazydw4rf/book-stock-manager/internal/controller"
	"github.com/crazydw4rf/book-stock-manager/internal/handler"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/swagger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func newConfig() (*config.Config, error) {
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func newValidator() *validator.Validate {
	return validator.New()
}

func newDBConn(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newFiberApp() (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Use(helmet.New())

	app.Get("/docs/*", swagger.New(swagger.Config{Title: "Book Stock Manager | API Docs"}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"api_name": "book stock manager",
			"version":  "1.0.0",
		})
	})

	return app, nil
}

func startApp(lc fx.Lifecycle, app *fiber.App, cfg *config.Config) {
	lc.Append(fx.StartHook(func() {
		listenAddr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)

		log.Printf("Starting app...listen to: http://%s\n", listenAddr)

		go func() {
			err := app.Listen(listenAddr)
			if err != nil {
				log.Printf("Error!! Shutting down app...%v\n", err)
				time.Sleep(time.Second * 2)
				os.Exit(1)
			}
		}()
	}))

	lc.Append(fx.StopHook(func(ctx context.Context) {
		err := app.ShutdownWithContext(ctx)
		log.Println("Shutting down app...", err.Error())
	}))
}

func main() {
	app := fx.New(
		fx.Provide(newConfig, newFiberApp, newDBConn, newValidator),
		fx.Provide(controller.NewBookController),
		fx.Decorate(handler.SetupBookHandler),
		fx.Invoke(startApp),
	)

	app.Run()
}
