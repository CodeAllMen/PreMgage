package sp

import (
	"encoding/json"
	"fmt"
	"github.com/MobileCPX/PreMgage/models/sp"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
)

type UnsubController struct {
	BaseController
}

type unsubResp struct {
	Response data `json:"response"`
}
type data struct {
	MigSid string `json:"mig_sid"`
	Status string `json:"status"`
}

func (c UnsubController) SubIDUnsub() {
	logs.Info("SubIDUnsub: ", c.Ctx.Input.URI())
	resp := new(unsubResp)
	subID := c.GetString("subID")
	serviceID := c.Ctx.Input.Param(":serviceID")
	serviceConfig := c.getService(serviceID)
	unsubURL := "https://api.migpay.com/subscriptions/" + subID + "/unsubscribe.json"
	fmt.Println(unsubURL)
	req := httplib.Post(unsubURL)
	req.Header("Host", "api.migpay.com")
	req.Header("MigPay-API-Key", "ec10abf99d5ef8e97c6731d74225c17ab0f94fd7")
	req.Header("Accept", "application/json")
	body, err := req.Bytes()


	if err != nil {
		c.Ctx.WriteString("Unsub Failed")
		c.StopRun()
	}
	err = json.Unmarshal(body, resp)
	if err != nil {
		c.Ctx.WriteString("Unsub Failed")
		c.StopRun()
	}
	logs.Info("SubIDUnsub", string(body))
	if resp.Response.MigSid != "" {
		mo := new(sp.Mo)
		err = mo.GetMoByUserID(resp.Response.MigSid)
		_, err = mo.UnsubUpdateMo()
		if err == nil {
			//c.Ctx.WriteString("Unsub SUCCESS")
			c.redirect(serviceConfig.UnsubSuccessURL)
		} else {
			//c.Ctx.WriteString("Unsub ERROR")
			c.redirect(serviceConfig.UnsubFailedURL)
		}
	} else {
		c.redirect(serviceConfig.UnsubFailedURL)
	}

}

func (c UnsubController) UnsubByMsisdn() {
	logs.Info("UnsubByMsisdn: ", c.Ctx.Input.URI())
	resp := new(unsubResp)
	msisdn := c.GetString("msisdn")
	serviceID := c.Ctx.Input.Param(":serviceID")
	serviceConfig := c.getService(serviceID)

	mo := new(sp.Mo)
	err := mo.GetMoByMsisdnAndService(msisdn, serviceID)
	if err != nil {
		c.redirect(serviceConfig.UnsubFailedURL)
	}

	unsubURL := "https://api.migpay.com/subscriptions/" + mo.SubscriptionID + "/unsubscribe.json"
	fmt.Println(unsubURL)
	req := httplib.Post(unsubURL)
	req.Header("Host", "api.migpay.com")
	req.Header("MigPay-API-Key", "ec10abf99d5ef8e97c6731d74225c17ab0f94fd7")
	req.Header("Accept", "application/json")
	body, err := req.Bytes()
	if err != nil {
		c.Ctx.WriteString("Unsub Failed")
		c.StopRun()
	}
	err = json.Unmarshal(body, resp)
	if err != nil {
		c.Ctx.WriteString("Unsub Failed")
		c.StopRun()
	}
	logs.Info("SubIDUnsub", string(body))
	if resp.Response.MigSid != "" {
		mo := new(sp.Mo)
		err = mo.GetMoByUserID(resp.Response.MigSid)
		_, err = mo.UnsubUpdateMo()
		if err == nil {
			//c.Ctx.WriteString("Unsub SUCCESS")
			c.redirect(serviceConfig.UnsubSuccessURL)
		} else {
			//c.Ctx.WriteString("Unsub ERROR")
			c.redirect(serviceConfig.UnsubFailedURL)
		}
	} else {
		c.redirect(serviceConfig.UnsubFailedURL)
	}

}
