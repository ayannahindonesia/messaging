package groups

import (
	"messaging/handlers"
	"messaging/middlewares"

	"github.com/labstack/echo"
)

func AdminGroup(e *echo.Echo) {
	g := e.Group("/admin")
	middlewares.SetClientJWTmiddlewares(g, "admin")

	g.GET("/message_sms", handlers.MessageOTPList)
	// config info
	// g.GET("/info", handlers.AsiraAppInfo)

}
