package main

import (
	"fmt"
	"log"

	"films_library/config"
	"films_library/internal/app"

	"github.com/joho/godotenv"
)

// @title Go Film Libary REST API
// @version 1.0
// @description Golang REST API  for managing films, directors and actors in a film library database.
// @contact.name Grigory Kovalenko
// @contact.url https://github.com/CodeMaster482
// @contact.email grigorikovalenko@gmail.com
// @BasePath /api/v1
func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load .env file")
		return
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	fmt.Println(cfg)

	app.Run(cfg)
}
