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

func TestMessageOTPSend(t *testing.T) {
	RebuildData()

	//UrlAuth := "https://sms241.xyz/sms/api_sms_otp_send_json.php"
	//UrlSmsAPI := "https://sms241.xyz/sms/api_sms_otp_send_json.php"

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
		"phone_number": "081234567890",
		"message":      "Kode OTP Anda adalah 997823",
	}

	// expect valid response
	auth.POST("/client/message_sms_send").
		WithJSON(payload).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	//test query params debug = true
	obj := auth.POST("/client/message_sms_send").
		WithQuery("debug", true).
		WithJSON(payload).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.ContainsKey("status").ValueEqual("status", "success")

	//test error
	auth.POST("/client/message_sms_send").WithJSON(nil).
		Expect().
		Status(http.StatusUnprocessableEntity).
		JSON().
		Object()
}

func TestMessageOTPList(t *testing.T) {
	RebuildData()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	clientToken := getAdminLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+clientToken)
	})

	// expect valid response
	auth.GET("/admin/message_sms").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	obj := auth.GET("/admin/message_sms").
		WithQuery("phone_number", "08").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.ContainsKey("data").NotEmpty()
	//log.Printf("%+v", obj)

	auth2 := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer invalid")
	})
	auth2.GET("/admin/message_sms").
		WithQuery("phone_number", "08").
		Expect().
		Status(http.StatusUnauthorized).
		JSON().
		Object()
}
