package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/nicolito128/nintendo-salta/app"
	_ "github.com/nicolito128/nintendo-salta/pkg/database"
)

func main() {
	_, err := os.Stat(".env")
	if err == nil {
		godotenv.Load(".env")
	}

	app.Start()
}
