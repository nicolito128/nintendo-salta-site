package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nicolito128/nintendo-salta/pkg/database"
	"github.com/nicolito128/nintendo-salta/pkg/models"
	"github.com/nicolito128/nintendo-salta/pkg/utils"
)

func GetUsers(ctx *fiber.Ctx) error {
	var users []models.User

	_, err := database.GetAll(ctx, &users)
	if err != nil {
		return ctx.JSON(fiber.Map{"status": err.StatusRequest, "error": err.WrappedError})
	}

	return ctx.JSON(utils.NewAPIResponse(ctx, http.StatusOK, users))
}

func GetRanking(ctx *fiber.Ctx) error {
	var users []models.User

	scope, err := database.Paginate(ctx)
	if err != nil {
		return ctx.JSON(fiber.Map{"status": http.StatusBadRequest, "error": err})
	}

	res := database.DB.Scopes(scope)
	if res.Error != nil {
		return ctx.JSON(fiber.Map{"status": http.StatusBadRequest, "error": res.Error})
	}

	res = res.Model(&models.User{}).Order("score DESC").Find(&users)
	if res.Error != nil {
		return ctx.JSON(fiber.Map{"status": http.StatusNotFound, "error": res.Error})
	}

	return ctx.JSON(utils.NewAPIResponse(ctx, http.StatusOK, users))
}

func GetUserByName(ctx *fiber.Ctx) error {
	name := ctx.Params("name", "")
	if name == "" {
		return ctx.JSON(fiber.Map{"status": http.StatusBadRequest, "error": "empty name param"})
	}

	var u models.User
	tx := database.DB.Model(&models.User{}).Where("name = ?", name).First(&u)
	if tx.Error != nil {
		return ctx.JSON(fiber.Map{"status": http.StatusNotFound, "error": tx.Error})
	}

	return ctx.JSON(utils.NewAPIResponse(ctx, http.StatusOK, u))
}

func PostUser(ctx *fiber.Ctx) error {
	var user models.APIUser

	err := json.Unmarshal(ctx.Body(), &user)
	if err != nil {
		return ctx.JSON(fiber.Map{"status": http.StatusBadRequest, "error": err})
	}

	u := &models.User{Name: user.Name}
	errQuery := database.Create(ctx, u)
	if errQuery != nil {
		return ctx.JSON(fiber.Map{"status": errQuery.StatusRequest, "error": errQuery.WrappedError})
	}

	ctx.Request().Header.Add("Location", fmt.Sprintf("/api/user/%s", u.Name))
	return ctx.JSON(utils.NewAPIResponse(ctx, http.StatusCreated, u))
}

func DeleteUser(ctx *fiber.Ctx) error {
	name := ctx.Params("name", "")
	if name == "" {
		return ctx.JSON(fiber.Map{"status": http.StatusBadRequest, "error": "empty name param"})
	}

	tx := database.DB.Unscoped().Where("name = ?", name).Delete(&models.User{})
	if tx.Error != nil {
		return ctx.JSON(fiber.Map{"status": http.StatusBadRequest, "error": tx.Error})
	}

	return ctx.JSON(utils.NewAPIResponse(ctx, http.StatusOK, fiber.Map{}))
}

func IncUserScore(ctx *fiber.Ctx) error {
	name := ctx.Params("name", "")
	if name == "" {
		return ctx.JSON(fiber.Map{"status": http.StatusBadRequest, "error": "empty name param"})
	}

	u := models.User{Name: name}
	tx := database.DB.Model(&models.User{}).Where("name = ?", name).First(&u)
	if tx.Error != nil {
		return ctx.JSON(fiber.Map{"status": http.StatusBadRequest, "error": tx.Error})
	}

	u.Score++
	tx = database.DB.Model(&models.User{}).Where("name = ?", name).Update("score", u.Score)
	if tx.Error != nil {
		return ctx.JSON(fiber.Map{"status": http.StatusNotFound, "error": tx.Error})
	}

	return ctx.JSON(utils.NewAPIResponse(ctx, http.StatusOK, fiber.Map{}))
}

func DecUserScore(ctx *fiber.Ctx) error {
	name := ctx.Params("name", "")
	if name == "" {
		return ctx.JSON(fiber.Map{"status": http.StatusBadRequest, "error": "empty name param"})
	}

	u := models.User{Name: name}
	tx := database.DB.Model(&models.User{}).Where("name = ?", name).First(&u)
	if tx.Error != nil {
		return ctx.JSON(fiber.Map{"status": http.StatusBadRequest, "error": tx.Error})
	}

	u.Score--
	tx = database.DB.Model(&models.User{}).Where("name = ?", name).Update("score", u.Score)
	if tx.Error != nil {
		return ctx.JSON(fiber.Map{"status": http.StatusNotFound, "error": tx.Error})
	}

	return ctx.JSON(utils.NewAPIResponse(ctx, http.StatusOK, fiber.Map{}))
}
