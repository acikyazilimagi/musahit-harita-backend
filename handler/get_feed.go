package handler

import (
	"github.com/acikkaynak/musahit-harita-backend/repository"
	"github.com/gofiber/fiber/v2"
)

// GetFeed godoc
//
//	@Summary	Get Feeds
//	@Tags		Feed
//	@Produce	json
//	@Success	200	{object}	feeds.Response
//	@Router		/feeds [GET]
func GetFeed(repo *repository.Repository) fiber.Handler {
	return getFeed(repo)
}

func getFeed(repo *repository.Repository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		feeds, err := repo.GetFeeds()
		if err != nil {
			return ctx.JSON(err)
		}

		return ctx.JSON(feeds)
	}
}
