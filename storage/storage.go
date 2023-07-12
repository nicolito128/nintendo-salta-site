/*
Storage es el paquete encargado de otorgar las interfaces y estructuras necesarias
para comunicarse e interactuar con la base de datos existente.
*/
package storage

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Storage representa una nueva base de datos.
type Storage interface {
	DB() *gorm.DB
}

// GetAll devuelve en la variable 'dest' todos los valores del modelo encontrado.
// Estos valores se encuentran paginados según la 'page' y 'pageSize' del contexto.
func GetAll(ctx *fiber.Ctx, store Storage, dest interface{}, conds ...interface{}) *gorm.DB {
	f, err := Paginate(ctx)
	if err != nil {
		store.DB().AddError(err)
		return store.DB()
	}

	scope := store.DB().Scopes(f)
	if scope.Error != nil {
		return scope
	}

	return store.DB().Find(&dest, conds)
}
