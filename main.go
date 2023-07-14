package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/nicolito128/nintendo-salta/server"
	"github.com/nicolito128/nintendo-salta/storage"
)

func main() {
	_, err := os.Stat(".env")
	if err == nil {
		godotenv.Load(".env")
	}

	// Definiendo la engine que se utilizará.
	// La engine no es más que la forma de
	// renderizar los archivos .html situados en /public
	engine := html.New("./public", ".html")

	// Definiendo el enrutador.
	// Se encarga de manejar las rutas del sitio y renderizar respuestas.
	router := fiber.New(fiber.Config{Views: engine})
	router.Static("/assets", "./public/assets")

	// Obteniendo el puerto de ejecución del servidor.
	// Es posible ingresarlo por defecto desde consola.
	port := flag.String("port", os.Getenv("port"), "the server port")
	flag.Parse()

	listenAddr := fmt.Sprintf(":%s", *port)

	// Creando un nuevo acceso a la base de datos.
	store := storage.NewSqliteStorage("database")

	// Iniciando la aplicación
	app := server.NewServer(listenAddr, router, store)
	log.Printf("Server running on port %s!\n", *port)
	log.Fatal(app.Start())
}
