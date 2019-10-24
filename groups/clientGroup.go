package groups

import (
	"messaging/handlers"
	"messaging/middlewares"

	"github.com/labstack/echo"
)

func ClientGroup(e *echo.Echo) {
	g := e.Group("/client")
	middlewares.SetClientJWTmiddlewares(g, "admin")

	g.POST("/message_sms_send", handlers.MessageOTPSend)
	g.GET("/message_sms", handlers.MessageOTPList)
	e.GET("/login_admin", handlers.AdminLogin)
}
