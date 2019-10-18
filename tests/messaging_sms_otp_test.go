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
	auth.POST("/client/send").
		WithJSON(payload).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	auth.POST("/client/send").WithJSON(nil).
		Expect().
		Status(http.StatusUnprocessableEntity).
		JSON().
		Object()
}
