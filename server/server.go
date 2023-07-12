/*
El paquete server es la pieza encargada de unir todo el sitio web,
se encarga de iniciar el servidor, manejar el enrutamiento y
comunicarse en las peticiones con la base de datos.
*/
package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicolito128/nintendo-salta/storage"
)

// Server representa un nuevo servidor.
type Server struct {
	listenAddr string
	router     *fiber.App
	store      storage.Storage
}

func NewServer(listenAddr string, router *fiber.App, store storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		router:     router,
		store:      store,
	}
}

// Start inicializa las rutas y corre el servidor.
func (s *Server) Start() error {
	// Groups
	root := s.router.Group("/")
	api := s.router.Group("/api")

	// GET
	root.Get("/", s.handleIndexPage)
	root.Get("/login", s.handleLoginPage)
	root.Get("/admin", onlyAdmin(s.store, s.handleAdminPage))

	api.Get("/users", s.handleUsers)
	api.Get("/users/ranking", s.handleUsersRanking)
	api.Get("/user/:name", s.handleUserByName)

	// POST
	root.Post("/login", s.handleLoginForm)

	api.Post("/user", s.handleNewUser)

	// Delete
	api.Delete("/user/:name", s.handleDeleteUser)

	// Put
	api.Put("/user/inc/:name", s.handleIncrementUserScore)
	api.Put("/user/dec/:name", s.handleDecrementUserScore)

	return s.router.Listen(s.listenAddr)
}

func (s *Server) handleIndexPage(ctx *fiber.Ctx) error {
	return ctx.Render("index", nil)
}

func (s *Server) handleAdminPage(ctx *fiber.Ctx) error {
	return ctx.Render("admin", fiber.Map{"title": "Admin Panel"})
}
