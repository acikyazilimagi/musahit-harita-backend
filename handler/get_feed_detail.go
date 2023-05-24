package handler

import (
	"strconv"

	"github.com/acikkaynak/musahit-harita-backend/repository"
	"github.com/gofiber/fiber/v2"
)

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
        feeds, err := repo.GetFeedDetail(id)
		if err != nil {
			return ctx.JSON(err)
		}

		return ctx.JSON(feeds)
    }
}