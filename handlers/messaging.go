package handlers

import (
	"messaging/models"
	"messaging/partner"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func SendMessage(c echo.Context) error {
	defer c.Request().Body.Close()

	messaging := models.Messaging{}

	payloadRules := govalidator.MapData{
		"phone_number": []string{"required"}, //TODO:, "id_phonenumber"
		"message":      []string{"required"},
	}
	validate := validateRequestPayload(c, payloadRules, &messaging)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//build otp request to partner
	number := messaging.PhoneNumber
	message := messaging.Message
	var data = []map[string]string{
		{"number": number, "message": message, "sendingdatetime": time.Now().String()},
	}
	conf := partner.Config{
		"apikey":      "7a8f16c956ae0c1e50461972d972d228",
		"callbackurl": "",
		"datapacket":  data,
	}
	//TODO: storing store
	response, err := partner.Send(conf)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "SMS Sent", "Config": conf, "Response": string(response)})
}
