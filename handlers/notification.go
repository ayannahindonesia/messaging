package handlers

import (
	"context"
	"log"
	"messaging/models"
	"net/http"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/labstack/echo"
)

func MessageNotificationSend(c echo.Context) error {
	defer c.Request().Body.Close()
	//get user id
	// user := c.Get("user")
	// token := user.(*jwt.Token)
	// claims := token.Claims.(jwt.MapClaims)
	//clientID, _ := strconv.Atoi(claims["jti"].(string))
	//get messaging model
	//messaging_model := models.Messaging{}
	start := time.Now()
	//firebase init
	ctx := context.Background()
	config := &firebase.Config{ProjectID: "asira-app-33ed7"}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// Obtain a messaging.Client from the App.
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "cEh41s_l_t4:APA91bGaE1OLrCN0P3myiSslwtddtmZMDj4uy_0YbJJ3qvt_N_f81HdxJL5juuuud18OW3zfKZqLDMbn83O1EoBBhGHvJMKupupb5CUsSaWc9A4b6bItmDEctwZ3F-5ENoJfHPZP4NMn"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Notification: &messaging.Notification{
			Title: "Olla Kevin",
			Body:  "hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai hai ",
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
	log.Println("Successfully sent message: ", response)
	elapsed := time.Since(start)
	log.Printf("Execution time took : %s", elapsed)

	return c.JSON(http.StatusOK, response)
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
