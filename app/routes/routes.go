package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nicolito128/nintendo-salta/app/handlers"
	api_routes "github.com/nicolito128/nintendo-salta/app/routes/api"
	"github.com/nicolito128/nintendo-salta/pkg/database"
	"github.com/nicolito128/nintendo-salta/pkg/models"
	"github.com/nicolito128/nintendo-salta/pkg/utils"
)

func Initialize(app *fiber.App) {
	// Index
	app.Get("/", Index)
	app.Get("/admin", handlers.OnlyAdmin(AdminPanel))
	app.Get("/login", Login)
	app.Post("/login", PostLogin)

	// API routes
	api := app.Group("/api")
	api.Get("/users", api_routes.GetUsers)
	api.Get("/users/ranking", api_routes.GetRanking)
	api.Get("/user/:name", api_routes.GetUserByName)
	api.Post("/user", api_routes.PostUser)
	api.Delete("/user/:name", api_routes.DeleteUser)
	api.Put("/user/inc/:name", api_routes.IncUserScore)
	api.Put("/user/dec/:name", api_routes.DecUserScore)
}

func Index(ctx *fiber.Ctx) error {
	return ctx.Render("index", nil)
}

func Login(ctx *fiber.Ctx) error {
	return ctx.Render("login", fiber.Map{"title": "Login"})
}

func PostLogin(ctx *fiber.Ctx) error {
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

	tx := database.DB.Model(&models.Session{}).Create(&session)
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

func AdminPanel(ctx *fiber.Ctx) error {
	return ctx.Render("admin", fiber.Map{"title": "Admin Panel"})
}
