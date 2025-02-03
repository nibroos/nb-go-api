package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/nb-go-api/service/internal/controller/rest"
	"github.com/nibroos/nb-go-api/service/internal/repository"
	"github.com/nibroos/nb-go-api/service/internal/service"
	"gorm.io/gorm"
)

func SetupIdentifierRoutes(identifiers fiber.Router, gormDB *gorm.DB, sqlDB *sqlx.DB) {
	identifierRepo := repository.NewIdentifierRepository(gormDB, sqlDB)
	identifierService := service.NewIdentifierService(identifierRepo)
	identifierController := rest.NewIdentifierController(identifierService)

	// prefix /identifiers

	identifiers.Post("/index-identifier", identifierController.ListIdentifiers)
	identifiers.Post("/show-identifier", identifierController.GetIdentifierByID)
	identifiers.Post("/create-identifier", identifierController.CreateIdentifier)
	identifiers.Post("/update-identifier", identifierController.UpdateIdentifier)
	identifiers.Post("/delete-identifier", identifierController.DeleteIdentifier)
	identifiers.Post("/restore-identifier", identifierController.RestoreIdentifier)
	identifiers.Post("/auth-index-identifier", identifierController.ListIdentifiersByAuthUser)
	identifiers.Post("/auth-show-identifier", identifierController.GetIdentifierByAuthUser)
	identifiers.Post("/auth-create-identifier", identifierController.CreateIdentifierByAuthUser)
	identifiers.Post("/auth-update-identifier", identifierController.UpdateIdentifierByAuthUser)
	identifiers.Post("/auth-delete-identifier", identifierController.DeleteIdentifierByAuthUser)
	identifiers.Post("/auth-restore-identifier", identifierController.RestoreIdentifierByAuthUser)
}
