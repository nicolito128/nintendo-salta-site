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

		var err error
		var exists int64
		var session models.Session

		store.DB().Model(&models.Session{}).Where("token = ?", strToken).Count(&exists)
		if exists == 0 {
			return ctx.SendStatus(http.StatusUnauthorized)
		}

		err = store.DB().Model(&models.Session{}).Where("token = ?", strToken).First(&session).Error
		if err != nil {
			return ctx.SendStatus(http.StatusUnauthorized)
		}

		if (session.Expire - time.Now().UnixNano()) <= 0 {
			err = store.DB().Model(&models.Session{}).Unscoped().Delete(session).Error
			if err != nil {
				return ctx.SendStatus(http.StatusUnauthorized)
			}

			return ctx.SendStatus(http.StatusUnauthorized)
		}

		return fn(ctx)
	}
}
