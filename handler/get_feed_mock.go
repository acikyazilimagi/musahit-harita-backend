package handler

import (
	"github.com/acikkaynak/musahit-harita-backend/repository/mock"
	"github.com/gofiber/fiber/v2"
)

// GetFeedMock godoc
//
//	@Summary	Get Feeds mock
//	@Tags		Feed
//	@Produce	json
//	@Success	200	{object}	feeds.Response
//	@Router		/feeds/mock [GET]
func GetFeedMock() fiber.Handler {
	return getFeedMock()
}

func getFeedMock() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		feeds, err := mock.GetFeeds()
		if err != nil {
			return ctx.JSON(err)
		}

		return ctx.JSON(feeds)
	}
}
