package handlers

import (
	"be-laporpak/config"
	"be-laporpak/models"
	"be-laporpak/utils"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// WebSocket Clients
var Clients = make(map[*fiber.Ctx]bool)

// Kirim update ke WebSocket
func broadcastUpdate(user models.User) {
	data, _ := json.Marshal(fiber.Map{
		"code":    "200",
		"status":  "OK",
		"message": "User updated",
		"data":    user,
	})

	for client := range Clients {
		_ = client.Write(data)
	}
}

// Register User
func RegisterUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    "400",
			"status":  "Bad Request",
			"message": "Invalid request body",
			"data":    fiber.Map{},
		})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := config.DB.NamedExec(`INSERT INTO users (nid, password, nama_user, roles, level, team, created_at, updated_at)
	VALUES (:nid, :password, :nama_user, :roles, :level, :team, :created_at, :updated_at)`, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    "500",
			"status":  "Internal Server Error",
			"message": "Database Error",
			"data":    fiber.Map{},
		})
	}

	token, _ := utils.GenerateJWT(user.NID)

	return c.JSON(fiber.Map{
		"code":    "201",
		"status":  "Created",
		"message": "User Registered",
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

// Login User
func LoginUser(c *fiber.Ctx) error {
	input := new(models.User)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    "400",
			"status":  "Bad Request",
			"message": "Invalid request body",
			"data":    fiber.Map{},
		})
	}

	var user models.User
	err := config.DB.Get(&user, "SELECT * FROM users WHERE nid=$1", input.NID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"code":    "401",
			"status":  "Unauthorized",
			"message": "Invalid NID or password",
			"data":    fiber.Map{},
		})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{
			"code":    "401",
			"status":  "Unauthorized",
			"message": "Invalid NID or password",
			"data":    fiber.Map{},
		})
	}

	token, _ := utils.GenerateJWT(user.NID)

	return c.JSON(fiber.Map{
		"code":    "200",
		"status":  "OK",
		"message": "Login successful",
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

// Get Profile
func GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	var userData models.User
	err := config.DB.Get(&userData, "SELECT nid, nama_user, roles, level, team, created_at FROM users WHERE nid=$1", user["nid"])
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"code":    "404",
			"status":  "Not Found",
			"message": "User not found",
			"data":    fiber.Map{},
		})
	}

	return c.JSON(fiber.Map{
		"code":    "200",
		"status":  "OK",
		"message": "Profile fetched",
		"data":    userData,
	})
}

// Delete User
func DeleteUser(c *fiber.Ctx) error {
	nid := c.Params("nid")
	_, err := config.DB.Exec("DELETE FROM users WHERE nid=$1", nid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    "500",
			"status":  "Internal Server Error",
			"message": "Failed to delete user",
			"data":    fiber.Map{},
		})
	}

	return c.JSON(fiber.Map{
		"code":    "200",
		"status":  "OK",
		"message": "User deleted",
		"data":    fiber.Map{},
	})
}

// Update user
func UpdateUser(c *fiber.Ctx) error {
	nid := c.Params("nid")
	var input models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    "400",
			"status":  "Bad Request",
			"message": "Invalid request body",
			"data":    fiber.Map{},
		})
	}

	var existingUser models.User
	err := config.DB.Get(&existingUser, "SELECT * FROM users WHERE nid=$1", nid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"code":    "404",
			"status":  "Not Found",
			"message": "User not found",
			"data":    fiber.Map{},
		})
	}

	if input.NamaUser != "" {
		existingUser.NamaUser = input.NamaUser
	}

	existingUser.UpdatedAt = time.Now()

	_, err = config.DB.NamedExec(`
		UPDATE users 
		SET nama_user=:nama_user, updated_at=:updated_at
		WHERE nid=:nid
	`, &existingUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    "500",
			"status":  "Internal Server Error",
			"message": "Failed to update user",
			"data":    fiber.Map{},
		})
	}

	// Kirim data ke WebSocket
	broadcastUpdate(existingUser)

	return c.JSON(fiber.Map{
		"code":    "200",
		"status":  "OK",
		"message": "User updated successfully",
		"data":    existingUser,
	})
}
