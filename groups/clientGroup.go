package groups

import (
	"messaging/handlers"
	"messaging/middlewares"

	"github.com/labstack/echo"
)

func ClientGroup(e *echo.Echo) {
	g := e.Group("/client")
	middlewares.SetClientJWTmiddlewares(g, "client")

	g.GET("/send", handlers.SendMessage)
}
