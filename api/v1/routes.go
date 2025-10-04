package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manthan307/nota-cms/api/v1/auth"
	"github.com/manthan307/nota-cms/api/v1/content"
	"github.com/manthan307/nota-cms/api/v1/media"
	schemasRoutes "github.com/manthan307/nota-cms/api/v1/schemas"
	db "github.com/manthan307/nota-cms/db/output"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func RegisterRoutes(app *fiber.App, queries *db.Queries, logger *zap.Logger, minioClient *minio.Client) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	//auth
	v1.Post("/auth/register", auth.RegisterHandler(queries, logger))
	v1.Post("/auth/login", auth.LoginHandler(queries, logger))

	//schemas
	schemas := v1.Group("/schemas")
	schemas.Post("/create", auth.ProtectedRoute(logger, queries, "editor"), schemasRoutes.SchemasCreateHandler(queries, logger))
	schemas.Post("/get_by_id/:id", auth.ProtectedRoute(logger, queries, "viewer"), schemasRoutes.GetSchemaByID(queries, logger))
	schemas.Post("/get_by_name/:name", auth.ProtectedRoute(logger, queries, "viewer"), schemasRoutes.GetSchemaByName(queries, logger))

	//content
	contentRoute := v1.Group("/content")
	contentRoute.Post("/create", auth.ProtectedRoute(logger, queries, "editor"), content.CreateContentHandler(queries, logger))
	contentRoute.Post("/delete/:id", auth.ProtectedRoute(logger, queries, "editor"), content.DeleteContentHandler(queries, logger))
	contentRoute.Post("/get/:id", content.GetContentHandler(queries, logger))
	contentRoute.Post("/get_all/:schema_name", content.GetAllContentsBySchemaHandler(queries, logger, true))

	//media
	mediaRoute := v1.Group("/media")
	mediaRoute.Post("/upload", auth.ProtectedRoute(logger, queries, "editor"), media.UploadMediaHandler(queries, logger, minioClient))
}
