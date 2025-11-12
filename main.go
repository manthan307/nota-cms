package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/manthan307/nota-cms/api"
	v1 "github.com/manthan307/nota-cms/api/v1"
	postgres "github.com/manthan307/nota-cms/db"
	db "github.com/manthan307/nota-cms/db/output"
	"github.com/manthan307/nota-cms/logger"
	minio_pkg "github.com/manthan307/nota-cms/utils/minio"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	fx.New(
		fx.Provide(
			logger.InitLogger,
			postgres.Connect,
			func(pool *pgxpool.Pool) *db.Queries {
				return db.New(pool)
			},

			minio_pkg.InitS3,
			api.StartServer,
		),
		fx.Invoke(
			postgres.RunMigrations,
			v1.RegisterRoutes,
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()

}
