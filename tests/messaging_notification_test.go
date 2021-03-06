package tests

import (
	"flag"
	"fmt"
	"messaging/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestMessageNotificationSend(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken) //adminBasicToken
	})

	clientToken := getClientLoginToken(e, auth, "0")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+clientToken)
	})
	//NOTE(RA): check run in test or normal execution
	if flag.Lookup("test.v") == nil {
		fmt.Println("normal run")
	} else {
		fmt.Println("run under go test")
	}

	payload := map[string]interface{}{
		"recipient_id":   "borrower-2312",
		"title":          "Your Loan Applied has been Aproved",
		"message_body":   "{\"status\":\"aproved\",\"product\":\"1276216\"}",
		"firebase_token": "cEh41s_l_t4:APA91bGaE1OLrCN0P3myiSslwtddtmZMDj4uy_0YbJJ3qvt_N_f81HdxJL5juuuud18OW3zfKZqLDMbn83O1EoBBhGHvJMKupupb5CUsSaWc9A4b6bItmDEctwZ3F-5ENoJfHPZP4NMn",
	}

	// expect valid response
	auth.POST("/client/message_notification_send").
		WithJSON(payload).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	auth.POST("/client/message_notification_send").WithJSON(nil).
		Expect().
		Status(http.StatusUnprocessableEntity).
		JSON().
		Object()
}
