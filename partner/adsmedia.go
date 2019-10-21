package partner

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"messaging/messaging"
	"net/http"

	"github.com/jarcoal/httpmock"
)

//TODO: OOP style partner
type (
	Config map[string]interface{}
)

//const URLendpoint string = "http://sms241.xyz/sms/api_sms_otp_send_json.php" // mahal
//const URLendpoint string = "http://sms241.xyz/sms/api_sms_masking_send_json.php" //agak mahal

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
		httpmock.RegisterResponder("POST", "https://sms241.xyz/sms/api_sms_masking_send_json.php",
			func(req *http.Request) (*http.Response, error) {
				//DONE: mockup with parameters and POST
				PayloadSmsOTP := make(map[string]interface{})
				if err := json.NewDecoder(req.Body).Decode(&PayloadSmsOTP); err != nil {
					return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
				}
				result := fmt.Sprintf(`{
					"id": 5,
					"created_time": "2019-10-21T12:34:28.726458+07:00",
					"updated_time": "2019-10-21T12:34:28.726458+07:00",
					"phone_number": "%s",
					"message": "%s",
					"partner": "adsmedia",
					"raw_response": ""{\"sending_respon\":[{\"globalstatus\":10,\"globalstatustext\":\"Success\",\"datapacket\":[{\"packet\":{\"number\":\"6282297335657\",\"sendingid\":1287265,\"sendingstatus\":10,\"sendingstatustext\":\"success\",\"price\":320}}]}]}",
					"status": true,
					"send_time": "2019-10-21T12:34:28.726458+07:00"
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
	Partners := messaging.App.Config.GetStringMap(fmt.Sprintf("%s.partners.adsmedia", messaging.App.ENV))
	hostUrl := Partners["host_url"].(string)
	request, err := http.Post(hostUrl, "application/json", bytes.NewBuffer(payload))
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

func PrepareRequestData(phoneNumber string, message string) (con Config) {
	var data = []map[string]string{
		{"number": phoneNumber, "message": message, "sendingdatetime": ""}, //time.Now().String()
	}
	Partners := messaging.App.Config.GetStringMap(fmt.Sprintf("%s.partners.adsmedia", messaging.App.ENV))
	conf := Config{
		"apikey":      Partners["api_key"],
		"callbackurl": "",
		"datapacket":  data,
	}
	return conf
}
