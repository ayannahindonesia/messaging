package groups

import (
	"messaging/handlers"
	"messaging/middlewares"

	"github.com/labstack/echo"
)

func AdminGroup(e *echo.Echo) {
	g := e.Group("/admin")
	middlewares.SetClientJWTmiddlewares(g, "admin")

	g.POST("/message_sms_send", handlers.MessageOTPSend)
	g.GET("/message_sms", handlers.MessageOTPList)
	e.POST("/login_admin", handlers.AdminLogin)
	// config info
	// g.GET("/info", handlers.AsiraAppInfo)

}
