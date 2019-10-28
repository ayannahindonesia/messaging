package handlers

import (
	"context"
	"fmt"
	"log"
	messaging_notif "messaging/messaging"
	"messaging/models"
	"net/http"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func getPushNotifClient() (context.Context, *messaging.Client, error) {
	//firebase init
	//TODO: config init.go
	ctx := context.Background()
	projectID := messaging_notif.App.Config.GetStringMap(fmt.Sprintf("%s.push_notification", messaging_notif.App.ENV))
	config := &firebase.Config{ProjectID: projectID["project_id"].(string)} //"asira-app-33ed7"
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return ctx, &messaging.Client{}, err
	}
	// Obtain a messaging.Client from the App.
	client, err := app.Messaging(ctx)
	if err != nil {
		return ctx, &messaging.Client{}, err
	}
	return ctx, client, nil
}
func MessageNotificationSend(c echo.Context) error {
	defer c.Request().Body.Close()
	//get user id
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	clientID, _ := strconv.Atoi(claims["jti"].(string))
	//internal benchmark
	start := time.Now()

	//get messaging model & validate
	notification := models.Notification{}
	payloadRules := govalidator.MapData{
		"title":        []string{"required"},
		"message_body": []string{"required"},
	}
	validate := validateRequestPayload(c, payloadRules, &notification)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}
	//check payload Topic / RegistrationToken must have value
	if len(notification.Topic) != 0 && len(notification.FirebaseToken) != 0 {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Topic Or RegistrationToken must have value")
	}

	//registrationToken := notification.RegistrationToken //"cEh41s_l_t4:APA91bGaE1OLrCN0P3myiSslwtddtmZMDj4uy_0YbJJ3qvt_N_f81HdxJL5juuuud18OW3zfKZqLDMbn83O1EoBBhGHvJMKupupb5CUsSaWc9A4b6bItmDEctwZ3F-5ENoJfHPZP4NMn"

	//get FCM
	ctx, client, err := getPushNotifClient()
	if err != nil {
		returnInvalidResponse(http.StatusInternalServerError, nil, fmt.Sprintf("error initializing app: %v\n", err))
	}

	//var message messaging.Message
	var response string
	//send by RegistrationToken
	if len(notification.FirebaseToken) != 0 {
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: notification.Title,
				Body:  notification.MessageBody,
			},
			Token: notification.FirebaseToken,
		}
		// Send a message to the device corresponding to the provided
		// registration token.
		response, err = client.Send(ctx, message)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, validate, err.Error())
		}
	}

	//send by Topic
	if len(notification.Topic) != 0 {
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: notification.Title,
				Body:  notification.MessageBody,
			},
			Topic: notification.Topic,
		}
		// Send a message to the device corresponding to the provided
		// registration token.
		response, err = client.Send(ctx, message)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, validate, err.Error())
		}
	}

	//storing
	notification.ClientID = clientID
	notification.Response = response
	err = notification.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal menyimpan notifikasi")
	}
	// Response is a message ID string.
	log.Println("Successfully sent message: ", response, clientID)
	elapsed := time.Since(start)
	log.Printf("Execution time took : %s", elapsed)

	return c.JSON(http.StatusOK, notification)
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
