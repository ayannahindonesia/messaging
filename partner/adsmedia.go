package partner

import (
	"bytes"
	"encoding/json"
	"errors"
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

//NOTE: expected return value from API partner
/*{
	"sending_respon":[
		{
			"globalstatus":10,
			"globalstatustext":"Success",
			"datapacket":[
				{
					"packet":{
						"number":"1",
						"sendingid":0,
						"sendingstatus":60,
						"sendingstatustext":"Invalid Number",
						"price":0
					}
				}
			]
		}
	]
}*/
type ResponseType struct {
	SendingRespon []struct {
		GlobalStatus     int    `json:"globalstatus"`
		GlobalStatusText string `json:"globalstatustext"`
		DataPacket       []struct {
			Packet struct {
				Number            string `json:"number"`
				SendingID         int    `json:"sendingid"`
				SendingStatus     int    `json:"sendingstatus"`
				SendingStatusText string `json:"sendingstatustext"`
				Price             int    `json:"price"`
			} `json:"packet"`
		} `json:"datapacket"`
	} `json:"sending_respon"`
}

//const URLendpoint string = "http://sms241.xyz/sms/api_sms_otp_send_json.php" // mahal
//const URLendpoint string = "http://sms241.xyz/sms/api_sms_masking_send_json.php" //agak mahal

func Send(con Config, debugFlag bool) (body []byte, err error) {
	//format data tuk testing dan juga debugFlag == true
	dummyData := `{
		"id": 5,
		"created_time": "2019-10-21T12:34:28.726458+07:00",
		"updated_time": "2019-10-21T12:34:28.726458+07:00",
		"client_id": 2,
		"phone_number": "%s",
		"message": "%s",
		"partner": "adsmedia",
		"raw_response": "{\"sending_respon\":[{\"globalstatus\":10,\"globalstatustext\":\"Success\",\"datapacket\":[{\"packet\":{\"number\":\"6282297335657\",\"sendingid\":1287265,\"sendingstatus\":10,\"sendingstatustext\":\"success\",\"price\":320}}]}]}",
		"status": "success",
		"send_time": "2019-10-21T12:34:28.726458+07:00"
	}`
	//paksa ambil dari query param (debugFlag)
	messaging.App.DebugMode = debugFlag

	//NOTE(RA): check run in test or normal execution
	if flag.Lookup("test.v") != nil {
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

				result := fmt.Sprintf(dummyData, PayloadSmsOTP["phone_number"], PayloadSmsOTP["message"])
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

	if debugFlag == false {
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
	} else {
		//debugFlag == true
		body = []byte(fmt.Sprintf(dummyData, "081234567890", payload))
	}
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

func GetStatusResponse(response []byte) (string, error) {

	status := "success"
	debugFlag := messaging.App.DebugMode

	//skip response use dummy result
	if debugFlag == true {
		return status, nil
	}

	var responseObj ResponseType
	//defer responseObj
	log.Printf("RES : %s", response)

	//parsing json
	err := json.Unmarshal(response, &responseObj)
	if err != nil {
		return "Invalid Number", err
	}

	//res := responseObj["sending_respon"].(map[string]interface{})
	//err = json.Unmarshal(responseObj["sending_respon"].([]byte), &res)
	log.Printf("%+v", responseObj)
	if responseObj.SendingRespon[0].DataPacket[0].Packet.SendingStatus != 10 {
		log.Printf("%+v", responseObj)
		return "failed", errors.New(responseObj.SendingRespon[0].DataPacket[0].Packet.SendingStatusText)
	}

	return status, nil
}
