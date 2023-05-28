package handler

import (
	"github.com/acikkaynak/musahit-harita-backend/feeds"
	"strconv"

	"github.com/acikkaynak/musahit-harita-backend/repository"
	"github.com/gofiber/fiber/v2"
)

// GetFeedDetail godoc
//
//	@Summary	Get Feed Detail
//	@Tags		Feed
//	@Produce	json
//	@Success	200	{object}	feeds.FeedDetailResponse
//	@Router		/feed/{neighborhoodId} [GET]
func GetFeedDetail(repo *repository.Repository) fiber.Handler {
	return getFeedDetail(repo)
}

func getFeedDetail(repo *repository.Repository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		neighborhoodId := ctx.Params("neighborhoodId", "0")
		id, err := strconv.Atoi(neighborhoodId)
		if err != nil || id <= 0 {
			return ctx.JSON(fiber.Map{
				"error": "neighborhoodId not found",
			})
		}
		//feeds, err := repo.GetFeedDetailFromMemory(id)
		// TODO: Get feed detail from database. For now, we are fetching it from aws s3 bucket and storing it in memory.
		//feeds, err := repo.GetFeedDetail(id)
		if err != nil {
			return ctx.JSON(err)
		}

		var feed feeds.FeedDetailResponse
		return ctx.JSON(feed)
	}
}
