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
		// TODO: Get feeds from database. For now, we are fetching it from aws s3 bucket and storing it in memory.
		//feeds, err := repo.GetFeeds()
		feeds, err := repo.GetFeedsFromMemory()
		if err != nil {
			return ctx.JSON(fiber.Map{
				"error": "feeds not found",
			})
		}

		return ctx.JSON(feeds)
	}
}
