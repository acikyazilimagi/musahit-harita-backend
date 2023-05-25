package main

import (
	"fmt"
	"github.com/acikkaynak/musahit-harita-backend/aws/s3"
	"github.com/acikkaynak/musahit-harita-backend/middleware/auth"
	"github.com/acikkaynak/musahit-harita-backend/middleware/cache"
	"os"
	"os/signal"
	"syscall"

	redisStore "github.com/acikkaynak/musahit-harita-backend/cache"
	_ "github.com/acikkaynak/musahit-harita-backend/docs"
	"github.com/acikkaynak/musahit-harita-backend/handler"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/acikkaynak/musahit-harita-backend/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Application struct {
	app        *fiber.App
	repository *repository.Repository
	validator  *validator.Validate
}

func (a *Application) RegisterApi() {
	a.app.Get("/", handler.RedirectSwagger)
	// monitor endpoint for pprof
	a.app.Get("/monitor", monitor.New())

	// health check endpoint for kubernetes
	a.app.Get("/healthz", handler.HealthCheck)

	// metrics endpoint for prometheus
	a.app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	a.app.Get("/feeds/", handler.GetFeed(a.repository))
	a.app.Get("/mock/feeds", handler.GetFeedMock())
	a.app.Get("/mock/feed/:neighborhoodId", handler.GetFeedDetailMock())
	a.app.Get("/feed/:neighborhoodId", handler.GetFeedDetail(a.repository))
	a.app.Post("/volunteer-form", handler.VolunteerForm(a.validator, a.repository))

	a.app.Post("/ovo-data", handler.UpdateOvoData())

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
	app.Use(auth.New())
	app.Use(cache.New())

	// gzip compression handled by nginx
	//app.Use(compress.New(compress.Config{
	//	Next: func(c *fiber.Ctx) bool {
	//		return c.Path() == "/swagger/index.html"
	//	},
	//	Level: compress.LevelBestCompression,
	//}))

	// register repositories to fiber app
	pgStore := repository.New()

	// store ovo building info
	obi := repository.NewOvoBuildingInfo(s3.DownloadMostRecentObject(os.Getenv("S3_BUCKET_NAME")))
	obi.Store()

	// register redis to fiber app
	cach, err := redisStore.NewRedisStore()
	if err != nil {
		log.Logger().Panic(fmt.Sprintf("redis error: %s", err.Error()))
		os.Exit(1)
	}

	application := &Application{
		app:        app,
		repository: pgStore,
		validator:  validator.New(),
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
	// close redis connection
	cach.Close()
}
