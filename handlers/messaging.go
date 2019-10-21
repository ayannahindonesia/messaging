package handlers

import (
	"log"
	"messaging/models"
	"messaging/partner"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func SendMessage(c echo.Context) error {
	defer c.Request().Body.Close()
	//get user id
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	clientID, _ := strconv.Atoi(claims["jti"].(string))
	//get messaging model
	messaging := models.Messaging{}

	//TODO: OTP generate 6 digit value
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
	//build request data
	conf := partner.PrepareRequestData(number, message)
	//send to API partner
	response, err := partner.Send(conf)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "")
	}

	//TODO: check return value from API Partner
	messaging.Status = true
	//TODO: dinamic partner setting
	messaging.Partner = "adsmedia"
	//raw response from API partner
	messaging.RawResponse = string(response)
	//clientID from jwt
	messaging.ClientID = clientID
	//DONE: storing models
	err = messaging.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal menyimpan pesan OTP")
	}
	log.Println(string(response))
	return c.JSON(http.StatusOK, messaging)
}
