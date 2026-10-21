package main

import (
	"fmt"
	"log"
	"os"

	"expense.com/m/adapter/routes"
	"expense.com/m/config"
	"expense.com/m/shared/token"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/lpernett/godotenv"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	// .env opsional: di Docker env disuplai lewat docker-compose
	_ = godotenv.Load()
	// .env.local dipakai untuk secret lokal yang tidak di-commit (mis. JWT_SECRET).
	_ = godotenv.Overload(".env.local")

	authCfg, err := config.LoadAuthConfig()
	if err != nil {
		log.Fatalf("Konfigurasi auth tidak valid: %v", err)
	}

	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	app := fiber.New(fiber.Config{
		ProxyHeader: "X-Forwarded-For",
	})

	// CORS — izinkan frontend (Next.js) mengakses API lintas origin/port.
	// Atur CORS_ORIGINS di .env untuk environment lain.
	app.Use(cors.New(cors.Config{
		AllowOrigins:     getEnv("CORS_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization",
		AllowCredentials: true,
	}))

	tokenManager := token.NewManager(authCfg.JWTSecret, authCfg.AccessTTL)
	routes.InitRouter(app, db, authCfg, tokenManager)

	for _, route := range app.Stack() {
		for _, r := range route {
			fmt.Printf("Method: %s, Path: %s\n", r.Method, r.Path)
		}
	}

	// Port backend dibaca dari env PORT (default 8000).
	// Karena Next.js dev juga memakai :3000, gunakan port lain
	// (mis. 8000) agar tidak bentrok.
	port := getEnv("PORT", "8000")
	log.Printf("Server berjalan di port %s", port)
	log.Fatal(app.Listen(":" + port))
}
