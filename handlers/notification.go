package handlers

import (
	"context"
	"encoding/json"
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

type (
	NotificationPayload struct {
		RecipientID   string            `json:"recipient_id"`
		Title         string            `json:"title"`
		MessageBody   string            `json:"message_body"`
		Data          map[string]string `json:"data"`
		Topic         string            `json:"topic"`
		FirebaseToken string            `json:"firebase_token"`
	}
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
	notificationPayload := NotificationPayload{}
	payloadRules := govalidator.MapData{
		"title":        []string{"required"},
		"recipient_id": []string{},
		"message_body": []string{},
		"data":         []string{},
	}
	validate := validateRequestPayload(c, payloadRules, &notificationPayload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//check payload Topic / RegistrationToken must have value
	if len(notificationPayload.Topic) == 0 && len(notificationPayload.FirebaseToken) == 0 {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "Topic Or FirebaseToken must have value")
	}

	//assignment object
	notification := models.Notification{}
	notification.RecipientID = notificationPayload.RecipientID
	notification.Title = notificationPayload.Title
	notification.MessageBody = notificationPayload.MessageBody
	notification.Topic = notificationPayload.Topic
	notification.FirebaseToken = notificationPayload.FirebaseToken

	//get FCM
	ctx, client, err := getPushNotifClient()
	if err != nil {
		returnInvalidResponse(http.StatusInternalServerError, nil, fmt.Sprintf("error initializing app: %v\n", err))
	}

	//var message messaging.Message
	var response string
	//marshalling object to json
	marshalledData, _ := json.Marshal(notificationPayload.Data)

	//send by RegistrationToken
	if len(notification.FirebaseToken) != 0 {
		message := &messaging.Message{
			Data: map[string]string{
				"Title":    notification.Title,
				"Body":     notification.MessageBody,
				"JsonData": string(marshalledData),
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
	Title := c.QueryParam("title")
	FirebaseToken := c.QueryParam("token")
	Topic := c.QueryParam("topic")
	layout := "2019-10-21T12:34:28.726458+07:00"
	SendTime, _ := time.Parse(layout, c.QueryParam("send_time"))
	type Filter struct {
		ID            int       `json:"id"`
		ClientID      int       `json:"client_id"`
		Title         string    `json:"title" condition:"LIKE"`
		FirebaseToken string    `json:"firebase_token" condition:"LIKE"`
		Topic         string    `json:"topic"`
		SendTime      time.Time `json:"send_time"`
	}

	notification := models.Notification{}
	result, err := notification.PagedFilterSearch(page, rows, order, sort, &Filter{
		ID:            id,
		ClientID:      ClientID,
		Title:         Title,
		FirebaseToken: FirebaseToken,
		Topic:         Topic,
		SendTime:      SendTime,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func ClientMessageNotificationList(c echo.Context) error {
	defer c.Request().Body.Close()

	///get jti
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	clientID, _ := strconv.Atoi(claims["jti"].(string))

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	order := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	// filters
	id, _ := strconv.Atoi(c.QueryParam("id"))
	Title := c.QueryParam("title")
	FirebaseToken := c.QueryParam("token")
	Topic := c.QueryParam("topic")
	layout := "2019-10-21T12:34:28.726458+07:00"
	SendTime, _ := time.Parse(layout, c.QueryParam("send_time"))
	type Filter struct {
		ID            int       `json:"id"`
		ClientID      int       `json:"client_id"`
		Title         string    `json:"title" condition:"LIKE"`
		FirebaseToken string    `json:"firebase_token" condition:"LIKE"`
		Topic         string    `json:"topic"`
		SendTime      time.Time `json:"send_time"`
	}

	notification := models.Notification{}
	result, err := notification.PagedFilterSearch(page, rows, order, sort, &Filter{
		ID:            id,
		ClientID:      clientID, //DONE id == jti
		Title:         Title,
		FirebaseToken: FirebaseToken,
		Topic:         Topic,
		SendTime:      SendTime,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}
