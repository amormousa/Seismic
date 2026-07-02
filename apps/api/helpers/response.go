package helpers

import "github.com/gofiber/fiber/v2"

/*
	Note: https://github.com/Majoramari/Hirelink/tree/main/src/errors
	I'm copying my way of doing things from Hirelink, I'm not sure if it's the best way.
	in go but I love it.
*/

// APIResponse is the consistent shape every single endpoint
// returns, so the frontend never has to guess the response
// structure differently per endpoint.
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success sends a 200 response with the standard success shape.
func Success(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response with the given status code
// and the standard error shape.
func Error(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}
