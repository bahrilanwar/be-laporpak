package main

import (
	"be-laporpak/config"
	"be-laporpak/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	// Koneksi ke database
	config.ConnectDB()

	app := fiber.New()

	// WebSocket untuk real-time update
	app.Get("/ws", websocket.New(routes.WebSocketHandler))

	// API Routes
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
