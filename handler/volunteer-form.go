package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type VolunteerFormRequest struct {
	Name            string `json:"name" validate:"required"`
	Surname         string `json:"surname" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"required,numeric,len=10"`
	KvkkAccepted    *bool  `json:"kvkkAccepted" validate:"required,eq=true"`
	neighbourhoodId int    `json:"neighbourhoodId" validate:"required"`
}

func VolunteerForm(val *validator.Validate) fiber.Handler {
	return volunteerForm(val)
}

func volunteerForm(val *validator.Validate) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		d := &VolunteerFormRequest{}
		err := ctx.BodyParser(&d)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err)
		}
		errors := validateStruct(val, d)
		if len(errors) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(errors)
		}
		return nil
	}
}
