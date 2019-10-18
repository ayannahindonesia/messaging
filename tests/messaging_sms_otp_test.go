package tests

import (
	"flag"
	"fmt"
	"log"
	"messaging/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/jarcoal/httpmock"
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

	//NOTE(RA): check run in test or normal execution
	//var DefaultTransport = http.DefaultTransport
	if flag.Lookup("test.v") == nil {
		log.Println("normal run")

	} else {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		log.Println("run under go test")

		//TODO mockup responder
		httpmock.RegisterResponder("POST", "/client/send",
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, `{
					"Config": {
						"apikey": "7a8f16c956ae0c1e50461972d972d228",
						"callbackurl": "",
						"datapacket": [
							{
								"message": "Kode OTP Anda adalah 99782377",
								"number": "082297335657",
								"sendingdatetime": ""
							}
						]
					},
					"Response": "{\"sending_respon\":[{\"globalstatus\":10,\"globalstatustext\":\"Success\",\"datapacket\":[{\"packet\":{\"number\":\"6282297335657\",\"sendingid\":1287265,\"sendingstatus\":10,\"sendingstatustext\":\"success\",\"price\":320}}]}]}",
					"message": "SMS Sent"
				}`)
				if err != nil {
					return httpmock.NewStringResponse(503, err.Error()), nil
				}
				return resp, nil
			},
		)
	}

	// expect valid response
	auth.POST("/client/send").
		WithJSON(payload).
		Expect().
		Status(http.StatusOK).
		JSON().
		String()

	log.Println(payload)
}
