package tests

import (
	"messaging/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/jarcoal/httpmock"
)

func TestSMSSendOTP(t *testing.T) {
	RebuildData()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//UrlAuth := "https://sms241.xyz/sms/api_sms_otp_send_json.php"
	//UrlSmsAPI := "https://sms241.xyz/sms/api_sms_otp_send_json.php"

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	clientToken := getClientLoginToken(e, auth, "1")

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+clientToken)
	})

	// valid response
	auth.POST("/client/send").
		Expect().
		Status(http.StatusOK).JSON().Object()
	/*
		// Exact URL match
		httpmock.RegisterResponder("GET", UrlAuth,
			httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))

		// Regexp match (could use httpmock.RegisterRegexpResponder instead)
		httpmock.RegisterResponder("GET", `=~^https://sms241\.xyz/articles/id/\d+\z`,
			httpmock.NewStringResponder(200, `{"id": 1, "name": "My Great Article"}`))

		// do stuff that makes a request to articles

		// get count info
		httpmock.GetTotalCallCount()
		// return an article related to the request with the help of regexp submatch (\d+)
		httpmock.RegisterResponder("GET", `=~^https://api\.mybiz\.com/articles/id/(\d+)\z`,
			func(req *http.Request) (*http.Response, error) {
				// Get ID from request
				id := httpmock.MustGetSubmatchAsUint(req, 1) // 1=first regexp submatch
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"id":   id,
					"name": "My Great Article",
				})
			},
		)

		// get the amount of calls for the registered responder
		info := httpmock.GetCallCountInfo()
		info["GET https://sms241.xyz/sms/api_sms_otp_send_json.php"]
		info["GET https://mybiz.xyz/articles/id/12"]
		info[`GET =~^https://mybiz\.xyz/articles/id/\d+\z`]
	*/

}
