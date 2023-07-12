package server

import (
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nicolito128/nintendo-salta/models"
	"github.com/nicolito128/nintendo-salta/pkg/utils"
)

func (s *Server) handleLoginPage(ctx *fiber.Ctx) error {
	return ctx.Render("login", fiber.Map{"title": "Login"})
}

func (s *Server) handleLoginForm(ctx *fiber.Ctx) error {
	name := ctx.FormValue("name", "")
	pass := ctx.FormValue("password", "")
	secret := os.Getenv("AUTHENTICATION_SECRET")

	if name == "" || len(name) < 4 {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if pass != secret {
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	token := utils.GenerateSecureToken(32)

	var session models.Session
	session.Name = name
	session.Token = token
	session.Expire = time.Now().UnixNano() + int64(time.Hour*5)

	tx := s.store.DB().Model(&models.Session{}).Create(&session)
	if tx.Error != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Render("login", fiber.Map{
		"title":  "Login",
		"name":   name,
		"token":  token,
		"expire": session.Expire,
	})
}
