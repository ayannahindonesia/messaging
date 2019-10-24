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

	// Obtain a messaging.Client from the App.
ctx := context.Background()
client, err := app.Messaging(ctx)
if err != nil {
        log.Fatalf("error getting Messaging client: %v\n", err)
}

// This registration token comes from the client FCM SDKs.
registrationToken := "YOUR_REGISTRATION_TOKEN"

// See documentation on defining a message payload.
message := &messaging.Message{
        Data: map[string]string{
                "score": "850",
                "time":  "2:45",
        },
        Token: registrationToken,
}

// Send a message to the device corresponding to the provided
// registration token.
response, err := client.Send(ctx, message)
if err != nil {
        log.Fatalln(err)
}
// Response is a message ID string.
fmt.Println("Successfully sent message:", response)
	return c.JSON(http.StatusOK, "OK"
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
