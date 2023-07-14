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

	scope, err := storage.Paginate(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	res := s.store.DB().Scopes(scope)
	if res.Error != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": res.Error.Error()})
	}

	err = res.Model(&models.User{}).Find(&users).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	var count int64
	err = s.store.DB().Model(&models.User{}).Count(&count).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	ctx.Set("X-Total-Count", fmt.Sprint(count))
	ctx.Status(http.StatusOK)
	return ctx.JSON(fiber.Map{"status": "success", "data": users})
}

func (s *Server) handleUsersRanking(ctx *fiber.Ctx) error {
	var users []models.User

	scope, err := storage.Paginate(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	res := s.store.DB().Scopes(scope)
	if res.Error != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": res.Error.Error()})
	}

	err = res.Model(&models.User{}).Order("score DESC").Find(&users).Error
	if res.Error != nil {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	var count int64
	err = s.store.DB().Model(&models.User{}).Count(&count).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	ctx.Set("X-Total-Count", fmt.Sprint(count))
	ctx.Status(http.StatusOK)
	return ctx.JSON(fiber.Map{"status": "success", "data": users})
}

func (s *Server) handleUserByName(ctx *fiber.Ctx) error {
	name := ctx.Params("name", "")
	if name == "" {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": "empty name param"})
	}

	var exists int64
	s.store.DB().Model(&models.User{}).Where("name = ?", name).Count(&exists)
	if exists == 0 {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(fiber.Map{"status": "fail", "error": "the user doesn't exists."})
	}

	var u models.User
	err := s.store.DB().Model(&models.User{}).Where("name = ?", name).First(&u).Error
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	ctx.Status(http.StatusOK)
	return ctx.JSON(fiber.Map{"status": "success", "data": u})
}

func (s *Server) handleNewUser(ctx *fiber.Ctx) error {
	var user models.APIUser

	err := json.Unmarshal(ctx.Body(), &user)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err})
	}

	var exists int64
	err = s.store.DB().Model(&models.User{}).Where("name = ?", user.Name).Count(&exists).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	if exists == 1 {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": "the user already exist"})
	}

	u := &models.User{Name: user.Name}
	err = s.store.DB().Model(&models.User{}).Create(&u).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
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

	var exists int64
	err := s.store.DB().Model(&models.User{}).Where("name = ?", name).Count(&exists).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err})
	}

	if exists == 0 {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(fiber.Map{"status": "fail", "error": "the user doesn't exist."})
	}

	u := models.User{Name: name}
	err = s.store.DB().Model(&models.User{}).Where("name = ?", name).First(&u).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	u.Score++
	err = s.store.DB().Model(&models.User{}).Where("name = ?", name).Update("score", u.Score).Error
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
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

	var exists int64
	err := s.store.DB().Model(&models.User{}).Where("name = ?", name).Count(&exists).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	if exists == 0 {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(fiber.Map{"status": "fail", "error": "the user doesn't exist."})
	}

	u := models.User{Name: name}
	err = s.store.DB().Model(&models.User{}).Where("name = ?", name).First(&u).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	u.Score--
	if u.Score < 0 {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(fiber.Map{"status": "fail", "error": "the score cannot be less than 0."})
	}

	err = s.store.DB().Model(&models.User{}).Where("name = ?", name).Update("score", u.Score).Error
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	ctx.Status(http.StatusOK)
	return ctx.JSON(fiber.Map{"status": "success", "data": u})
}

func (s *Server) handleClearUsers(ctx *fiber.Ctx) error {
	var err error
	var users []models.User

	err = s.store.DB().Model(&models.User{}).Find(&users).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	err = s.store.DB().Model(&models.User{}).Unscoped().Delete(&users).Error
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"status": "success", "data": users})
}
