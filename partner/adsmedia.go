package partner

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jarcoal/httpmock"
)

type (
	Config map[string]interface{}
)

//const URLendpoint string = "http://sms241.xyz/sms/api_sms_otp_send_json.php"
const URLendpoint string = "http://sms241.xyz/sms/api_sms_masking_send_json.php"

func Send(con Config) (body []byte, err error) {
	//NOTE(RA): check run in test or normal execution
	//var DefaultTransport = http.DefaultTransport
	if flag.Lookup("test.v") == nil {
		log.Println("normal run")

	} else {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		log.Println("run under go test")

		//DONE mockup responder
		httpmock.RegisterResponder("POST", "http://sms241.xyz/sms/api_sms_masking_send_json.php",
			func(req *http.Request) (*http.Response, error) {
				//DONE: mockup with parameters and POST
				PayloadSmsOTP := make(map[string]interface{})
				if err := json.NewDecoder(req.Body).Decode(&PayloadSmsOTP); err != nil {
					return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
				}
				result := fmt.Sprintf(`{
					"Config": {
						"apikey": "7a8f16c956ae0c1e50461972d972d228",
						"callbackurl": "",
						"datapacket": [
							{
								"number": "%s",
								"message": "%s",
								"sendingdatetime": ""
							}
						]
					},
					"Response": "{\"sending_respon\":[{\"globalstatus\":10,\"globalstatustext\":\"Success\",\"datapacket\":[{\"packet\":{\"number\":\"6282297335657\",\"sendingid\":1287265,\"sendingstatus\":10,\"sendingstatustext\":\"success\",\"price\":320}}]}]}",
					"message": "SMS Sent"
				}`, PayloadSmsOTP["phone_number"], PayloadSmsOTP["message"])
				fmt.Println(result)
				resp, err := httpmock.NewJsonResponse(http.StatusOK, result)
				if err != nil {
					return httpmock.NewStringResponse(http.StatusInternalServerError, err.Error()), nil
				}
				return resp, nil
			},
		)
	}

	//set payload
	payload, _ := json.Marshal(con)
	//send request
	request, err := http.Post(URLendpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	defer request.Body.Close()

	body, err = ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))

	return body, err
}
