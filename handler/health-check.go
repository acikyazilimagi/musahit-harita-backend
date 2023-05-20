package handler

import "github.com/gofiber/fiber/v2"

// HealthCheck is the health check endpoint for kubernetes
// @Summary Health Check
// @Description Health Check
// @Tags Health Check
// @Accept */*
// @Produce json
// @Success 200 {string} nil
// @Router /healthz [get]

func HealthCheck(ctx *fiber.Ctx) error {
	return ctx.SendStatus(200)
}
