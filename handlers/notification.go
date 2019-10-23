package handlers

import (
	"log"
	"messaging/models"
	"messaging/partner"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func MessageNotificationSend(c echo.Context) error {
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
	//get response from parsing json API partner
	status, err := partner.GetStatusResponse(response)
	/*
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "Internal Error : "+err.Error())
		}
	*/
	//TODO: check return value from API Partner
	messaging.Status = status
	//error message
	if err == nil {
		messaging.ErrorMessage = ""
	} else {
		messaging.ErrorMessage = err.Error()
	}
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

func MessageNotificationList(c echo.Context) error {
	defer c.Request().Body.Close()

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	order := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	// filters
	id, _ := strconv.Atoi(c.QueryParam("id"))
	ClientID, _ := strconv.Atoi(c.QueryParam("client_id"))
	PhoneNumber := c.QueryParam("phone_number")
	Message := c.QueryParam("message")
	Partner := c.QueryParam("partner")
	Status := c.QueryParam("status")
	layout := "2019-10-21T12:34:28.726458+07:00"
	SendTime, _ := time.Parse(layout, c.QueryParam("send_time"))
	type Filter struct {
		ID          int       `json:"id"`
		ClientID    int       `json:"client_id"`
		PhoneNumber string    `json:"phone_number" condition:"LIKE"`
		Message     string    `json:"message" condition:"LIKE"`
		Partner     string    `json:"partner" condition:"LIKE"`
		Status      string    `json:"status"`
		SendTime    time.Time `json:"send_time"`
	}

	messaging := models.Messaging{}
	result, err := messaging.PagedFilterSearch(page, rows, order, sort, &Filter{
		ID:          id,
		ClientID:    ClientID,
		PhoneNumber: PhoneNumber,
		Message:     Message,
		Partner:     Partner,
		Status:      Status,
		SendTime:    SendTime,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}
