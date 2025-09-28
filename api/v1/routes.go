package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manthan307/nota-cms/api/v1/auth"
	db "github.com/manthan307/nota-cms/db/output"
	"go.uber.org/zap"
)

func RegisterRoutes(app *fiber.App, queries *db.Queries, logger *zap.Logger) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/auth/register", auth.Register(app, queries, logger))
	v1.Post("/auth/login", auth.Login(app, queries, logger))
}
