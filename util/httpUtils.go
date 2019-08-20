package util

import (
	"bytes"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"log"
	"net/http"
)

func HttpPostRequest(requestData []byte, requestURL, contentType string) (responseData []byte) {
	body := bytes.NewBuffer([]byte(requestData))

	//res, err := http.Post(requestURL, "text/xml;charset=utf-8", body)
	res, err := http.Post(requestURL, contentType, body)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer res.Body.Close()
	responseData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		logs.Error("HTTP请求失败，url：", requestURL, "data: ", string(requestData), " type : ", contentType)
		log.Fatal(err)
		return
	}
	return
}

