package helper

import "github.com/gofiber/fiber/v2"

func Success(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func SuccessWithMeta(c *fiber.Ctx, status int, message string, data interface{}, meta interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}

func Error(c *fiber.Ctx, status int, message string, errors interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
		"errors":  errors,
	})
}