package handlers

import (
	"messaging/partner"
	"net/http"

	"github.com/labstack/echo"
)

func SendMessage(c echo.Context) error {
	defer c.Request().Body.Close()
	var data = []map[string]string{
		{"number": "08981073502", "message": "hallo", "sendingdatetime": ""},
	}

	conf := partner.Config{
		"apikey":      "7a8f16c956ae0c1e50461972d972d228",
		"callbackurl": "",
		"datapacket":  data,
	}
	err := partner.Send(conf)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Sms Send", "Config": conf})
}
