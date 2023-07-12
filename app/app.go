package app

import (
	"fmt"
	"log"
	"os"

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

	log.Printf("Application started at port %s!", os.Getenv("PORT"))
	err := App.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}
}
