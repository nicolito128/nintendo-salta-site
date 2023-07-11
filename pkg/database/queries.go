package database

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type QueryError struct {
	StatusRequest int
	WrappedError  error
}

func (qe QueryError) Error() string {
	return qe.WrappedError.Error()
}

func GetAll(ctx *fiber.Ctx, dest interface{}, conds ...interface{}) (*gorm.DB, *QueryError) {
	f, err := Paginate(ctx)
	if err != nil {
		return nil, &QueryError{http.StatusBadRequest, err}
	}

	scope := DB.Scopes(f)
	if scope.Error != nil {
		return nil, &QueryError{http.StatusBadRequest, scope.Error}
	}

	rows := scope.Find(dest, conds)
	if rows.Error != nil {
		return nil, &QueryError{http.StatusBadRequest, rows.Error}
	}

	return rows, nil
}

func Get(dest interface{}, conds ...interface{}) *QueryError {
	row := DB.First(dest, conds)
	if row.Error != nil {
		return &QueryError{http.StatusBadRequest, row.Error}
	}

	return nil
}

func GetByID(ctx *fiber.Ctx, dest interface{}, conds ...interface{}) *QueryError {
	idParam := ctx.Query("id", "")
	if strings.TrimSpace(idParam) == "" {
		return &QueryError{http.StatusBadRequest, errors.New("empty id")}
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return &QueryError{http.StatusBadRequest, err}
	}

	return Get(dest, id, conds)
}

func Create(ctx *fiber.Ctx, dest interface{}) *QueryError {
	res := DB.Create(dest)
	if res.Error != nil {
		return &QueryError{http.StatusBadRequest, res.Error}
	}

	return nil
}

func Delete(dest interface{}, conds ...interface{}) *QueryError {
	row := DB.Unscoped().Delete(dest, conds)
	if row.Error != nil {
		return &QueryError{http.StatusBadRequest, row.Error}
	}

	return nil
}

func DeleteByID(ctx *fiber.Ctx, dest interface{}, conds ...interface{}) *QueryError {
	id := ctx.Query("id", "")
	if strings.TrimSpace(id) == "" {
		return &QueryError{http.StatusBadRequest, errors.New("empty id")}
	}

	return Delete(dest, conds)
}
