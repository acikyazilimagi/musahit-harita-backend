package main

import (
	"fmt"
	_ "github.com/acikkaynak/musahit-harita-backend/docs"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/acikkaynak/musahit-harita-backend/repository"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	app        *fiber.App
	repository *repository.Repository
}

func (a *Application) RegisterApi() {
	// monitor endpoint for pprof
	a.app.Get("/monitor", monitor.New())

	// health check endpoint for kubernetes
	a.app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// metrics endpoint for prometheus
	a.app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// swagger docs endpoint
	route := a.app.Group("/swagger")
	route.Get("*", swagger.HandlerDefault)
}

// @title Musahit Harita API
// @version 1.0
// @description Musahit Harita API
// @contact.name Acik Kaynak
// @license.name Apache 2.0
// @license.url https://github.com/acikkaynak/musahit-harita-backend/blob/main/LICENSE
// @host localhost:80
func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: jsoniter.Marshal,
		JSONDecoder: jsoniter.Unmarshal,
	})

	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(pprof.New())
	app.Use(cache.New())

	// register repositories to fiber app
	pgStore := repository.New()

	application := &Application{
		app:        app,
		repository: pgStore,
	}

	application.RegisterApi()

	// gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)

	log.Logger().Info("application is running..")
	go func() {
		_ = <-c
		log.Logger().Info("application gracefully shutting down..")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":80"); err != nil {
		log.Logger().Panic(fmt.Sprintf("app error: %s", err.Error()))
	}

	log.Logger().Info("application cleanup tasks..")
	// close database connection
	pgStore.Close()
}
