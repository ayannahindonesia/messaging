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

func TestSMSSendOTP(t *testing.T) {
	RebuildData()

	//UrlAuth := "https://sms241.xyz/sms/api_sms_otp_send_json.php"
	//UrlSmsAPI := "https://sms241.xyz/sms/api_sms_otp_send_json.php"

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
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
		"phone_number": "082297335657",
		"message":      "Kode OTP Anda adalah 997823",
	}

	// expect valid response
	auth.POST("/client/message_otp_send").
		WithJSON(payload).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	auth.POST("/client/message_otp_send").WithJSON(nil).
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

	clientToken := getClientLoginToken(e, auth, "0")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+clientToken)
	})

	// expect valid response
	auth.GET("/client/message_otp").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	//get by phone number
	obj := auth.GET("/client/message_otp").
		WithQuery("phone_number", "08").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.ContainsKey("data").NotEmpty()
	//get status success
	obj = auth.GET("/client/message_otp").
		WithQuery("status", "success").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	//assert.Equal(t, 1, obj.ContainsKey("data").Values().Length())
	//assert.Equal(t, "success", )
	obj.ContainsKey("status").Value("success")

	//negative test, invalid token
	auth2 := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer invalid")
	})
	auth2.GET("/client/message_otp").
		WithQuery("phone_number", "08").
		Expect().
		Status(http.StatusUnauthorized).
		JSON().
		Object()
}
