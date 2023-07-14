package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleAdminPage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/index", fiber.Map{"title": "Admin Panel", "site_location": os.Getenv("SITE_LOCATION")})
}

func (s *Server) handleAdminSearchUserPage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/buscar", fiber.Map{"title": "Buscar Participante", "site_location": os.Getenv("SITE_LOCATION")})
}

func (s *Server) handleAdminAddUserPage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/agregar", fiber.Map{"title": "Agregar Participante", "site_location": os.Getenv("SITE_LOCATION")})
}

func (s *Server) handleAdminScorePage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/modificar-puntaje", fiber.Map{"title": "Modificar Puntaje", "site_location": os.Getenv("SITE_LOCATION")})
}

func (s *Server) handleAdminDeletePage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/eliminar", fiber.Map{"title": "Eliminar Participante", "site_location": os.Getenv("SITE_LOCATION")})
}

func (s *Server) handleAdminClearPage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/limpiar", fiber.Map{"title": "Limpiar Base de Datos", "site_location": os.Getenv("SITE_LOCATION")})
}
