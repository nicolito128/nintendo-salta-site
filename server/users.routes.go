package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nicolito128/nintendo-salta/models"
	"github.com/nicolito128/nintendo-salta/storage"
)

func (s *Server) handleUsers(ctx *fiber.Ctx) error {
	var users []models.User

	tx := storage.GetAll(ctx, s.store, &users)
	if tx.Error != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": tx.Error})
	}

	ctx.Status(http.StatusOK)
	return ctx.JSON(fiber.Map{"status": "success", "data": users})
}

func (s *Server) handleUsersRanking(ctx *fiber.Ctx) error {
	var users []models.User

	scope, err := storage.Paginate(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err})
	}

	res := s.store.DB().Scopes(scope)
	if res.Error != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": res.Error})
	}

	res = res.Model(&models.User{}).Order("score DESC").Find(&users)
	if res.Error != nil {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(fiber.Map{"status": "fail", "error": res.Error})
	}

	ctx.Status(http.StatusOK)
	return ctx.JSON(fiber.Map{"status": "success", "data": users})
}

func (s *Server) handleUserByName(ctx *fiber.Ctx) error {
	name := ctx.Params("name", "")
	if name == "" {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": "empty name param"})
	}

	var u models.User
	tx := s.store.DB().Model(&models.User{}).Where("name = ?", name).First(&u)
	if tx.Error != nil {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(fiber.Map{"status": "fail", "error": tx.Error})
	}

	ctx.Status(http.StatusOK)
	return ctx.JSON(fiber.Map{"status": "success", "data": u})
}

func (s *Server) handleNewUser(ctx *fiber.Ctx) error {
	var user models.APIUser

	err := json.Unmarshal(ctx.Body(), &user)
	if err != nil {
		return ctx.JSON(fiber.Map{"status": http.StatusBadRequest, "error": err})
	}

	u := &models.User{Name: user.Name}
	tx := s.store.DB().Model(&models.User{}).Create(&u)
	if tx.Error != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": tx.Error})
	}

	ctx.Status(http.StatusCreated)
	ctx.Request().Header.Add("Location", fmt.Sprintf("/api/user/%s", u.Name))
	return ctx.JSON(fiber.Map{"status": "success", "data": u})
}

func (s *Server) handleDeleteUser(ctx *fiber.Ctx) error {
	name := ctx.Params("name", "")
	if name == "" {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": "empty name param"})
	}

	tx := s.store.DB().Unscoped().Where("name = ?", name).Delete(&models.User{})
	if tx.Error != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": tx.Error})
	}

	ctx.Status(http.StatusOK)
	return ctx.JSON(fiber.Map{"status": "success", "data": fiber.Map{"name": name}})
}

func (s *Server) handleIncrementUserScore(ctx *fiber.Ctx) error {
	name := ctx.Params("name", "")
	if name == "" {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": "empty name param"})
	}

	u := models.User{Name: name}
	tx := s.store.DB().Model(&models.User{}).Where("name = ?", name).First(&u)
	if tx.Error != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": tx.Error})
	}

	u.Score++
	tx = s.store.DB().Model(&models.User{}).Where("name = ?", name).Update("score", u.Score)
	if tx.Error != nil {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(fiber.Map{"status": "fail", "error": tx.Error})
	}

	ctx.Status(http.StatusOK)
	return ctx.JSON(fiber.Map{"status": "success", "data": u})
}

func (s *Server) handleDecrementUserScore(ctx *fiber.Ctx) error {
	name := ctx.Params("name", "")
	if name == "" {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": "empty name param"})
	}

	u := models.User{Name: name}
	tx := s.store.DB().Model(&models.User{}).Where("name = ?", name).First(&u)
	if tx.Error != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": tx.Error})
	}

	u.Score--
	tx = s.store.DB().Model(&models.User{}).Where("name = ?", name).Update("score", u.Score)
	if tx.Error != nil {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(fiber.Map{"status": "fail", "error": tx.Error})
	}

	ctx.Status(http.StatusOK)
	return ctx.JSON(fiber.Map{"status": "success", "data": u})
}
