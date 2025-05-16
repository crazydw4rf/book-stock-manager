package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/crazydw4rf/book-stock-manager/docs"
	"github.com/crazydw4rf/book-stock-manager/internal/config"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME)

	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	log.Println("Connected to database successfully")

	return db, nil
}

func newFiberApp() (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	app.Use(cors.New())
	app.Use(helmet.New())

	if config.API_DOCS_ENABLED {
		app.Get("/docs/*", swagger.New(swagger.Config{
			Title:           "Book Stock Manager | API Docs",
			TryItOutEnabled: false,
		}))
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"api_name": "book stock manager",
			"version":  config.APP_VERSION,
			"env":      config.APP_ENV,
		})
	})

	return app, nil
}

func startApp(lc fx.Lifecycle, app *fiber.App, cfg *config.Config) {
	lc.Append(fx.StartHook(func() {
		listenAddr := fmt.Sprintf("%s:%d", cfg.APP_HOST, cfg.APP_PORT)

		log.Printf("Starting app in %s mode...\n", config.APP_ENV)
		log.Printf("API version: %s\n", config.APP_VERSION)
		log.Printf("Listening on: http://%s:%d\n", cfg.APP_HOST, cfg.APP_PORT)
		log.Printf("Swagger docs available at: http://%s:%d/docs/\n", cfg.APP_HOST, cfg.APP_PORT)

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
		log.Println("Received shutdown signal, gracefully shutting down...")
		err := app.ShutdownWithContext(ctx)
		if err != nil {
			log.Printf("Error during shutdown: %v\n", err)
		}
		log.Println("Application has been shutdown")
	}))
}
