package sendData

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
)

type SpNotification struct {
	NotificationType string ` json:"notification_type"`
	SubscriptionID   string `json:"subscription_id"`
	Sendtime         string `json:"sendtime"`
	TransactionID    string `json:"transaction_id"`
	ServiceID        string `json:"service_id"`

	Operator        string `json:"operator"`
	Msisdn          string `json:"msisdn"`
	CampID          int  `json:"camp_id"`
	OfferID         int  `json:"offer_id"`
	AffName         string `json:"aff_name"`
	PubID           string `json:"pub_id"`
	ClickID         string `json:"click_id"`
	PostbackStatus  int    `json:"postback_status"`
	PostbackMessage string `json:"postback_message"`
}

func (data *SpNotification) SendData() {
	jsonStr, _ := json.Marshal(data)
	httpPostJson(jsonStr)
}

func httpPostJson(jsonStr []byte) {

	url := "http://offer.globaltraffictracking.com/sp/data"
	//url := "http://127.0.0.1:8081/sp/data"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	logs.Info(string(jsonStr), "#####", string(body))
}
