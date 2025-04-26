package bootstrap

import (
	"batikin-be/internal/infra/postgresql"
	"flag"
	"os"

	"gorm.io/gorm"
)

func CommandHandler(db *gorm.DB) {
	var migrationCmd string
	var seederCmd bool

	flag.StringVar(&migrationCmd, "m", "", "Migrate database 'up' or 'down'")
	flag.BoolVar(&seederCmd, "s", false, "Seed database")
	flag.Parse()

	if migrationCmd != "" {
		postgresql.Migrate(db, migrationCmd)
		os.Exit(0)
	}
}
