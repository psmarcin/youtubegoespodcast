package ports

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const (
	error404Message = "I'm so sorry, didn't find page that you are looking for. Please try again in a bit."
	error500Message = "I'm so sorry, we have technical problem, please try again in a bit."
)

// errorHandler is server route handler for internal errors and not found routes
func errorHandler(ctx *fiber.Ctx) error {
	status := http.StatusNotFound
	message := error404Message

	ctx.Set("content-type", "text/html; charset=utf-8")
	err := ctx.Status(status).Render("error", fiber.Map{
		"errorCode":    status,
		"errorMessage": message,
	})

	if err != nil {
		l.WithError(err).Errorf("error while rendering template")
		return err
	}

	return nil
}
