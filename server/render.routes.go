package server

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nicolito128/nintendo-salta/models"
	"github.com/nicolito128/nintendo-salta/storage"
	"gorm.io/gorm"
)

func (s *Server) handleRenderRanking(ctx *fiber.Ctx) error {
	var users []models.User

	scope, err := storage.Paginate(ctx)
	if err != nil {
		return ctx.SendString("")
	}

	res := s.store.DB().Scopes(scope)
	if res.Error != nil {
		return ctx.SendString("")
	}

	res = res.Model(&models.User{}).Order("score DESC").Find(&users)
	if res.Error != nil {
		return ctx.SendString("")
	}

	var response string
	for i := 0; i < len(users); i++ {
		response += "<tr>"
		response += fmt.Sprintf(`
			<td>%d</td>
			<td>%s</td>
			<td>%d</td>
		`, i+1, users[i].Name, users[i].Score)
		response += "</tr>"
	}

	return ctx.SendString(response)
}

func (s *Server) handleRenderSearch(ctx *fiber.Ctx) error {
	var user models.User
	name := ctx.Params("name")
	if name == "" || len(name) <= 1 {
		return ctx.SendString("Busqueda no válida.")
	}

	var exists int64
	s.store.DB().Model(&models.User{}).Where("name = ?", name).Count(&exists)
	if exists == 0 {
		return ctx.SendString("El participante no existe.")
	}

	err := s.store.DB().Model(&models.User{}).Where("name = ?", name).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.SendString("Sin resultados.")
	}

	var response = fmt.Sprintf(`
		<table class="users-table">
			<thead>
				<tr>
					<th>ID</th>
					<th>Nombre</th>
					<th>Puntaje</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<td>%d</td>
					<td>%s</td>
					<td>%d</td>
				</tr>
			</tbody>
		</table>
	`, user.ID, user.Name, user.Score)

	return ctx.SendString(response)
}
