package backData

import (
	"encoding/json"
	"github.com/MobileCPX/PreBaseLib/splib/admindata"
	"github.com/MobileCPX/PreMgage/models/sp"
	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego/orm"
)

//var ServiceCamp = map[string]int{"PinkCity4K-MEO": 5, "PinkCity4K-NOS": 6, "Fit8Tube-MEO": 34, "Fit8Tube-NOS": 35}

func SendMo() {
	//UpdateMO()
	o := orm.NewOrm()
	mos := new([]sp.Mo)
	o.QueryTable("mo").OrderBy("id").Filter("subtime__gt", "2019-07-01").All(mos)
	for _, mo := range *mos {
		sendNoti := new(admindata.Notification)
		promoterID := 1
		if mo.OfferID != 0 {
			postback := new(sp.Postback)
			_ = postback.CheckOfferID(mo.OfferID)
			if postback.PromoterName == "张艳阳" {
				promoterID = 2
			}
			sendNoti.CampID = postback.CampID
		}

		sendNoti.PostbackPrice = 8

		sendNoti.Operator = mo.Operator
		sendNoti.OfferID = int(mo.OfferID)
		sendNoti.SubscriptionID = mo.SubscriptionID
		sendNoti.ServiceID = mo.ServiceID
		sendNoti.ClickID = mo.ClickID
		sendNoti.Msisdn = mo.Msisdn
		sendNoti.PubID = mo.PubID
		sendNoti.PostbackStatus = mo.PostbackStatus
		sendNoti.PostbackMessage = mo.PostbackCode
		sendNoti.TransactionID = ""
		sendNoti.AffName = mo.AffName
		if sendNoti.AffName == "" {
			sendNoti.AffName = "未知"
		}

		sendNoti.PromoterID = promoterID
		sendNoti.Sendtime = mo.Subtime
		sendNoti.NotificationType = "SUB"
		data, _ := json.Marshal(sendNoti)
		logs.Info(string(data))
		sendNoti.SendData(admindata.PROD)
	}
}
