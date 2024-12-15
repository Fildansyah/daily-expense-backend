package main

import (
	"fmt"
	"log"

	"expense.com/m/adapter/routes"
	"expense.com/m/config"
	"github.com/gofiber/fiber/v2"
	"github.com/lpernett/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := config.InitDatabase()

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	app := fiber.New()

	routes.InitRouter(app, db)

	for _, route := range app.Stack() {
		for _, r := range route {
			fmt.Printf("Method: %s, Path: %s\n", r.Method, r.Path)
		}
	}

	log.Fatal(app.Listen(":3000"))
}
