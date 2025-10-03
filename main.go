package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/manthan307/nota-cms/api"
	v1 "github.com/manthan307/nota-cms/api/v1"
	postgres "github.com/manthan307/nota-cms/db"
	db "github.com/manthan307/nota-cms/db/output"
	"github.com/manthan307/nota-cms/logger"
	"github.com/manthan307/nota-cms/utils/minio"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	fx.New(
		fx.Provide(
			logger.InitLogger,
			postgres.Connect,
			func(pool *pgxpool.Pool) *db.Queries {
				ctx := context.Background()
				conn, err := pool.Acquire(ctx)
				if err != nil {
					panic(err)
				}
				return db.New(conn.Conn())
			},
			api.StartServer,
		),
		fx.Invoke(
			postgres.RunMigrations,
			v1.RegisterRoutes,
			minio.InitBucket,
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()

}
