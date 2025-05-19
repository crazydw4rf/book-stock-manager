package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/crazydw4rf/book-stock-manager/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Usage: go run migrate.go <dir> <up|down>")
		os.Exit(1)
	}

	var (
		migrateDir    = args[1]
		migrateAction = args[2]
	)

	cfg, err := config.InitConfig()
	if err != nil {
		log.Println("Error loading config:", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME)

	m, err := migrate.New("file://"+migrateDir, dsn)
	if err != nil {
		log.Println("Error creating migration instance:", err)
		os.Exit(1)
	}

	switch migrateAction {
	case "down":
		err = m.Down()
	case "up":
		err = m.Up()
	default:
		fmt.Println("Invalid argument. Use 'up' or 'down'.")
		os.Exit(1)
	}

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No changes to apply.")
			os.Exit(0)
		}

		log.Println("Error applying migration:", err)
		os.Exit(1)
	}
}
