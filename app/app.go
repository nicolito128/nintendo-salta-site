package app

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/nicolito128/nintendo-salta/app/routes"
)

var App *fiber.App

func Start() {
	// Initialize standard Go html template engine
	engine := html.New("./public", ".html")

	App = fiber.New(fiber.Config{
		Views: engine,
	})

	App.Static("/", "./public/static")

	// Routes loading
	routes.Initialize(App)

	log.Println("Application started at port 3000!")
	err := App.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
