package partner

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type (
	Config map[string]interface{}
)

//const URLendpoint string = "http://sms241.xyz/sms/api_sms_otp_send_json.php"
const URLendpoint string = "http://sms241.xyz/sms/api_sms_masking_send_json.php"

func Send(con Config) (body []byte, err error) {

	payload, _ := json.Marshal(con)
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
