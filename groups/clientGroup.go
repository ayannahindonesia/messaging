package groups

import (
	"messaging/handlers"
	"messaging/middlewares"

	"github.com/labstack/echo"
)

func ClientGroup(e *echo.Echo) {
	g := e.Group("/client")
	middlewares.SetClientJWTmiddlewares(g, "client")

	g.POST("/message_otp_send", handlers.MessageOTPSend)
	g.GET("/message_otp", handlers.MessageOTPList)
}
