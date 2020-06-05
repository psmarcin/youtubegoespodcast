package api

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/requestid"
	"net/http"
)

const (
	error404Message = "I'm so sorry, didn't find page that you are looking for. Please try again in a bit."
	error500Message = "I'm so sorry, we have technical problem, please try again in a bit."
)

func errorHandler(c *fiber.Ctx) {
	status := http.StatusNotFound
	rId := requestid.Get(c)
	message := error404Message
	err := c.Error()
	if err != nil {
		status = http.StatusInternalServerError
		message = error500Message
	}

	c.Set("content-type", "text/html; charset=utf-8")
	c.Status(status)
	_ = templates.ExecuteTemplate(c.Fasthttp.Response.BodyWriter(), "error.tmpl", map[string]interface{}{
		"requestID":    rId,
		"errorCode":    status,
		"errorMessage": message,
	})
}
