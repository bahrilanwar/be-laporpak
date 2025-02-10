package routes

import (
	"be-laporpak/handlers"
	"be-laporpak/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// WebSocket Handler
func WebSocketHandler(c *websocket.Conn) {
	handlers.Clients[c] = true
	defer delete(handlers.Clients, c)

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}

// SetupRoutes mengatur semua route dalam aplikasi
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/user")

	// Public routes (tanpa proteksi JWT)
	api.Post("/register", handlers.RegisterUser) // Registrasi user
	api.Post("/login", handlers.LoginUser)       // Login user

	// Middleware JWT untuk proteksi route di bawah ini
	api.Use(middleware.JWTProtected())

	// Protected routes
	api.Post("/logout", handlers.LogoutUser)        // Logout user
	api.Get("/profile", handlers.GetProfile)        // Ambil profil user
	api.Put("/update/:nid", handlers.UpdateUser)    // Update data user
	api.Delete("/delete/:nid", handlers.DeleteUser) // Hapus user
}
