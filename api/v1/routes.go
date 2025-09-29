package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manthan307/nota-cms/api/v1/auth"
	schemasRoutes "github.com/manthan307/nota-cms/api/v1/schemas"
	db "github.com/manthan307/nota-cms/db/output"
	"go.uber.org/zap"
)

func RegisterRoutes(app *fiber.App, queries *db.Queries, logger *zap.Logger) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	//auth
	v1.Post("/auth/register", auth.RegisterHandler(queries, logger))
	v1.Post("/auth/login", auth.LoginHandler(queries, logger))

	//schemas
	schemas := v1.Group("/schemas", auth.ProtectedRoute(logger, queries))
	schemas.Post("/create", schemasRoutes.SchemasCreateHandler(queries, logger))
	schemas.Post("/get_by_id/:id", schemasRoutes.GetSchemaByID(queries, logger))
	schemas.Post("/get_by_name/:name", schemasRoutes.GetSchemaByName(queries, logger))
}
