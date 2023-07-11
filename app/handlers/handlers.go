package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nicolito128/nintendo-salta/pkg/database"
	"github.com/nicolito128/nintendo-salta/pkg/models"
)

func OnlyAdmin(fn fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		strToken := ctx.Cookies("token")
		if strToken == "" {
			return ctx.SendStatus(http.StatusUnauthorized)
		}

		var session models.Session

		tx := database.DB.Model(&models.Session{}).Where("token = ?", strToken).First(&session)
		if tx.Error != nil {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		if (session.Expire - time.Now().UnixNano()) <= 0 {
			tx = database.DB.Model(&models.Session{}).Unscoped().Delete(session)
			if tx.Error != nil {
				return ctx.SendStatus(http.StatusInternalServerError)
			}

			return ctx.SendStatus(http.StatusBadRequest)
		}

		return fn(ctx)
	}
}
