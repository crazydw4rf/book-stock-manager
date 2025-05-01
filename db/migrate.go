package main

import (
	"fmt"
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
		panic(err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	m, err := migrate.New("file://"+migrateDir, dsn)
	if err != nil {
		panic(err)
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
		panic(err)
	}
}
