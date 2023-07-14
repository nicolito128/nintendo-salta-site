package server

import "github.com/gofiber/fiber/v2"

func (s *Server) handleAdminPage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/index", fiber.Map{"title": "Admin Panel"})
}

func (s *Server) handleAdminSearchUserPage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/buscar", fiber.Map{"title": "Buscar Participante"})
}

func (s *Server) handleAdminAddUserPage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/agregar", fiber.Map{"title": "Agregar Participante"})
}

func (s *Server) handleAdminScorePage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/modificar-puntaje", fiber.Map{"title": "Modificar Puntaje"})
}

func (s *Server) handleAdminDeletePage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/eliminar", fiber.Map{"title": "Eliminar Participante"})
}

func (s *Server) handleAdminClearPage(ctx *fiber.Ctx) error {
	return ctx.Render("admin/limpiar", fiber.Map{"title": "Limpiar Base de Datos"})
}
