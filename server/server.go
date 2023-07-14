/*
El paquete server es la pieza encargada de unir todo el sitio web,
se encarga de iniciar el servidor, manejar el enrutamiento y
comunicarse en las peticiones con la base de datos.
*/
package server

import (
	"os"

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
	admin := s.router.Group("/admin")
	api := s.router.Group("/api")
	render := s.router.Group("/render")

	// GET
	root.Get("/", s.handleIndexPage)
	root.Get("/login", s.handleLoginPage)
	admin.Get("/", onlyAdmin(s.store, s.handleAdminPage))
	admin.Get("/buscar", onlyAdmin(s.store, s.handleAdminSearchUserPage))
	admin.Get("/agregar", onlyAdmin(s.store, s.handleAdminAddUserPage))
	admin.Get("/modificar-puntaje", onlyAdmin(s.store, s.handleAdminScorePage))
	admin.Get("/eliminar", onlyAdmin(s.store, s.handleAdminDeletePage))
	admin.Get("/limpiar", onlyAdmin(s.store, s.handleAdminClearPage))

	render.Get("/ranking", s.handleRenderRanking)
	render.Get("/search/:name", s.handleRenderSearch)

	api.Get("/users", s.handleUsers)
	api.Get("/users/ranking", s.handleUsersRanking)
	api.Get("/user/:name", s.handleUserByName)

	// POST
	root.Post("/login", s.handleLoginForm)

	api.Post("/user", s.handleNewUser)

	// Delete
	api.Delete("/user/:name", s.handleDeleteUser)
	api.Delete("/users", s.handleClearUsers)

	// Put
	api.Put("/user/inc/:name", s.handleIncrementUserScore)
	api.Put("/user/dec/:name", s.handleDecrementUserScore)

	return s.router.Listen(s.listenAddr)
}

func (s *Server) handleIndexPage(ctx *fiber.Ctx) error {
	return ctx.Render("index", fiber.Map{"site_location": os.Getenv("SITE_LOCATION")})
}
