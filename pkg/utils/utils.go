package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"
)

func NewAPIResponse(ctx *fiber.Ctx, status int, data interface{}) fiber.Map {
	return fiber.Map{"status": status, "data": data}
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
