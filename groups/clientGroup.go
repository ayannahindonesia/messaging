package groups

import (
	"messaging/handlers"
	"messaging/middlewares"

	"github.com/labstack/echo"
)

func ClientGroup(e *echo.Echo) {
	g := e.Group("/client")
	middlewares.SetClientJWTmiddlewares(g, "client")

	g.POST("/admin_login", handlers.AdminLogin)
	g.POST("/message_sms_send", handlers.MessageOTPSend)
	g.POST("/message_notification_send", handlers.MessageNotificationSend)
	g.GET("/message_notification", handlers.MessageNotificationList)
	//TODO: endpoint list
}
