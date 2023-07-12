package server

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nicolito128/nintendo-salta/models"
	"github.com/nicolito128/nintendo-salta/storage"
)

// onlyAdmin es un handler que actua de intermediario al momento de intentar ingresar a la web.
// Si los usuarios no pueden proporcionar un token válido el handler te devolverá un código de
// error "http.StatusUnauthorized".
func onlyAdmin(store storage.Storage, fn fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		strToken := ctx.Cookies("token")
		if strToken == "" {
			return ctx.SendStatus(http.StatusUnauthorized)
		}

		var session models.Session

		tx := store.DB().Model(&models.Session{}).Where("token = ?", strToken).First(&session)
		if tx.Error != nil {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		if (session.Expire - time.Now().UnixNano()) <= 0 {
			tx = store.DB().Model(&models.Session{}).Unscoped().Delete(session)
			if tx.Error != nil {
				return ctx.SendStatus(http.StatusInternalServerError)
			}

			return ctx.SendStatus(http.StatusBadRequest)
		}

		return fn(ctx)
	}
}
