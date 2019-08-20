package util

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego/logs"
)

// CmHTTPRequest http请求，根据cmid获取用户手机号、发起扣费和退订请求
func CmHTTPRequest(types string, requestData []byte) (body []byte, err error) {
	requestURL := cmRequestURL(types)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(requestData))
	req.Header.Set("Content-Encoding", "UTF-8")
	req.Header.Set("Content-Type", "text/xml")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.Error("cmHTTPRequest " + types + string(requestData) + " err: " + err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	return
}

func cmRequestURL(requestType string) (requestURL string) {
	switch requestType {
	case "getMsisdn":
		requestURL = "https://secure.cm.nl/contentbillingapi/GetMsisdn.ashx"
	case "billingRequest":
		requestURL = "http://billing.cm.nl/ContentBilling/Gateway/Request.ashx"
	case "unsub":
		requestURL = "http://DCBcancellation.cm.nl/unsubscribehandler.ashx"
	case "get_cmid":
		requestURL = "http://mcb.cmtelecom.nl"
	}
	return
}
