package handler

import (
	"github.com/acikkaynak/musahit-harita-backend/repository"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetFeedsOvoMock(repo *repository.Repository) fiber.Handler {
	return getFeedsOvoMock(repo)
}

func getFeedsOvoMock(repo *repository.Repository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		feeds, err := repo.GetFeedsFromMemory()
		if err != nil {
			return ctx.JSON(fiber.Map{
				"error": "feeds not found",
			})
		}

		return ctx.JSON(feeds)
	}
}

func GetFeedDetailOvoMock(repo *repository.Repository) fiber.Handler {
	return getFeedDetailOvoMock(repo)
}

func getFeedDetailOvoMock(repo *repository.Repository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		neighborhoodId := ctx.Params("neighborhoodId", "0")
		id, err := strconv.Atoi(neighborhoodId)
		if err != nil || id <= 0 {
			return ctx.JSON(fiber.Map{
				"error": "neighborhoodId not found",
			})
		}
		feeds, err := repo.GetFeedDetailFromMemory(id)
		// TODO: Get feed detail from database. For now, we are fetching it from aws s3 bucket and storing it in memory.
		//feeds, err := repo.GetFeedDetail(id)
		if err != nil {
			return ctx.JSON(err)
		}

		return ctx.JSON(feeds)
	}
}
