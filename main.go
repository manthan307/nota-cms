package main

import (
	"github.com/joho/godotenv"
	"github.com/manthan307/nota-cms/api"
	v1 "github.com/manthan307/nota-cms/api/v1"
	postgres "github.com/manthan307/nota-cms/db"
	db "github.com/manthan307/nota-cms/db/output"
	"github.com/manthan307/nota-cms/logger"
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
			db.New,
			api.StartServer,
		),
		fx.Invoke(
			v1.RegisterRoutes,
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()

}
