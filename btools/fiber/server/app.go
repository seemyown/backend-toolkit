package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type StaticConfig struct {
	Dir    string
	Prefix string
}

type ServerConfig struct {
	FiberConfig       *fiber.Config
	CorsConfig        *cors.Config
	CustomMiddlewares *[]func(*fiber.Ctx) error
	ClosureMiddleware *[]func(...any) fiber.Handler
	UseCompress       bool
	Static            *StaticConfig
}

func NewServer(config *ServerConfig) *fiber.App {
	app := fiber.New(*config.FiberConfig)

	if config.CorsConfig != nil {
		app.Use(cors.New(*config.CorsConfig))
	}
	if config.CustomMiddlewares != nil {
		for _, middleware := range *config.CustomMiddlewares {
			app.Use(middleware)
		}
	}
	if config.ClosureMiddleware != nil {
		for _, middleware := range *config.ClosureMiddleware {
			app.Use(middleware)
		}
	}
	if config.UseCompress {
		app.Use(compress.New())
	}

	if config.Static != nil {
		app.Static(config.Static.Dir, config.Static.Prefix)
	}
	return app
}
