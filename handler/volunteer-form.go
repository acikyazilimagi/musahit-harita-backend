package handler

import (
	"github.com/acikkaynak/musahit-harita-backend/model"
	"github.com/acikkaynak/musahit-harita-backend/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type VolunteerFormRequest struct {
	Name           string `json:"name" validate:"required"`
	Surname        string `json:"surname" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Phone          string `json:"phone" validate:"required,numeric,len=10"`
	KvkkAccepted   *bool  `json:"kvkkAccepted" validate:"required,eq=true"`
	NeighborhoodId int    `json:"neighborhoodId" validate:"required"`
}

func (r *VolunteerFormRequest) ToModel() model.VolunteerDoc {
	return model.VolunteerDoc{
		Name:           r.Name,
		Surname:        r.Surname,
		Email:          r.Email,
		Phone:          r.Phone,
		KvkkAccepted:   *r.KvkkAccepted,
		NeighborhoodId: r.NeighborhoodId,
	}
}

func VolunteerForm(val *validator.Validate, repo *repository.Repository) fiber.Handler {
	return volunteerForm(val, repo)
}

func volunteerForm(val *validator.Validate, repo *repository.Repository) fiber.Handler {
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
		id, error := repo.ApplyVolunteer(d.ToModel())
		if error != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
				"message": error.Error(),
			})
		}
		if id == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
				"message": "Volunteer already exists",
			})
		}
		return ctx.Status(fiber.StatusCreated).JSON(map[string]interface{}{
			"id": id,
		})
	}
}
