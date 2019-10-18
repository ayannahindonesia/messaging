package handlers

import (
	"flag"
	"fmt"
	"log"
	"messaging/models"
	"messaging/partner"
	"net/http"

	"github.com/jarcoal/httpmock"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func SendMessage(c echo.Context) error {

	//NOTE(RA): check run in test or normal execution
	if flag.Lookup("test.v") == nil {
		fmt.Println("normal run")
	} else {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		fmt.Println("run under go test")
	}

	defer c.Request().Body.Close()

	messaging := models.Messaging{}

	payloadRules := govalidator.MapData{
		"phone_number": []string{"required"}, //TODO:, "id_phonenumber"
		"message":      []string{"required"},
	}
	validate := validateRequestPayload(c, payloadRules, &messaging)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//build otp request to partner
	number := messaging.PhoneNumber
	message := messaging.Message
	var data = []map[string]string{
		{"number": number, "message": message, "sendingdatetime": ""}, //time.Now().String()
	}
	conf := partner.Config{
		"apikey":      "7a8f16c956ae0c1e50461972d972d228",
		"callbackurl": "",
		"datapacket":  data,
	}
	//TODO: storing store

	//normal run
	//DONE: check dev or prod ??
	var response []byte
	if flag.Lookup("test.v") == nil {
		log.Println("normal run")
		response, err = partner.Send(conf)
		//
	} else {
		//TODO mockup responder
		httpmock.RegisterResponder("POST", "/client/send",
			func(req *http.Request) (*http.Response, error) {
				resp, err := httpmock.NewJsonResponse(200, nil)
				if err != nil {
					return httpmock.NewStringResponse(500, err.Error()), nil
				}
				return resp, nil
			},
		)

		response, err = partner.Send(conf)

		// Exact URL match
		/*			httpmock.RegisterResponder("GET", UrlAuth,
						httpmock.NewStringResponder(200, `[{"id": 1, "name": "My Great Article"}]`))

					// Regexp match (could use httpmock.RegisterRegexpResponder instead)
					httpmock.RegisterResponder("GET", `=~^https://sms241\.xyz/articles/id/\d+\z`,
						httpmock.NewStringResponder(200, `{"id": 1, "name": "My Great Article"}`))

					// do stuff that makes a request to articles

					// get count info
					httpmock.GetTotalCallCount()


					// get the amount of calls for the registered responder
					info := httpmock.GetCallCountInfo()
					info["GET https://sms241.xyz/sms/api_sms_otp_send_json.php"]
					info["GET https://mybiz.xyz/articles/id/12"]
					info[`GET =~^https://mybiz\.xyz/articles/id/\d+\z`]
		*/

	}
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "SMS Sent", "Config": conf, "Response": string(response)})
}
