package handler

import (
	"github.com/acikkaynak/musahit-harita-backend/aws/s3"
	"github.com/acikkaynak/musahit-harita-backend/repository"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func UpdateOvoData() fiber.Handler {
	return updateOvoData()
}

func updateOvoData() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		body := ctx.Request().Body()

		type UpdateOvoData struct {
			MostRecent bool   `json:"most_recent"`
			Bucket     string `json:"bucket"`
			Key        string `json:"key,omitempty"`
		}

		var data UpdateOvoData
		json.Unmarshal(body, &data)

		var s3Object s3.ObjectData
		if data.MostRecent {
			s3Object = s3.DownloadMostRecentObject(data.Bucket)
		} else {
			s3Object = s3.Download(data.Bucket, data.Key)
		}

		obi := repository.NewOvoBuildingInfo(s3Object)
		if obi != nil {
			obi.Store()
		}

		return ctx.JSON(
			fiber.Map{
				"success": true,
			},
		)
	}
}
