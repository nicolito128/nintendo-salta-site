package storage

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Paginate checks if the request contains the query parameters "page" and "pageSize".
// By default, it sets page to 1 and pageSize to 100.
func Paginate(ctx *fiber.Ctx) (func(db *gorm.DB) *gorm.DB, error) {
	strPage := ctx.Query("page", "0")
	strPageSize := ctx.Query("pageSize", "0")

	page, err := strconv.Atoi(strPage)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}

	size, err := strconv.Atoi(strPageSize)
	if err != nil {
		return nil, err
	}

	// Filter out invalid page sizes
	switch {
	case size > 100:
		size = 100
	case size <= 0:
		size = 10
	}

	ctx.Set("page", fmt.Sprint(page))
	ctx.Set("page_size", fmt.Sprint(size))
	return func(db *gorm.DB) *gorm.DB {
		// Where to start looking for the database rows
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}, nil
}
