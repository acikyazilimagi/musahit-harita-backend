package handler

import (
	"strconv"

	"github.com/acikkaynak/musahit-harita-backend/repository/mock"
	"github.com/gofiber/fiber/v2"
)

// GetFeedDetailMock godoc
//
//	@Summary	Get Feed Detail mock
//	@Tags		Feed
//	@Produce	json
//	@Success	200	{object}	feeds.FeedDetailResponse
//	@Router		/feed/mock/{neighborhoodId} [GET]
func GetFeedDetailMock() fiber.Handler {
	return getFeedDetailMock()
}

func getFeedDetailMock() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		neighborhoodId := ctx.Params("neighborhoodId", "0")
		id, err := strconv.Atoi(neighborhoodId)
		if err != nil || id <= 0 {
			return ctx.JSON(fiber.Map{
				"error": "neighborhoodId not found",
			})
		}
		feed, err := mock.GetFeedDetail(id)
		if err != nil {
			return ctx.JSON(err)
		}
		return ctx.JSON(feed)
	}
}
