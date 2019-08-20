package sp

import (
	"fmt"
	"github.com/MobileCPX/PreMgage/models/sendData"
	"github.com/astaxie/beego/orm"
)

func SendAdminData() {
	o := orm.NewOrm()
	mos := new([]Mo)
	o.QueryTable(MoTBName()).OrderBy("-id").All(mos)
	fmt.Println(len(*mos))

	for _, mo := range *mos {
		sendNoti := new(sendData.SpNotification)
		sendNoti.OfferID = mo.OfferID
		sendNoti.SubscriptionID = mo.SubscriptionID
		sendNoti.ServiceID = mo.ServiceID
		sendNoti.ClickID = mo.ClickID
		sendNoti.CampID = mo.CampID
		sendNoti.PubID = mo.PubID
		sendNoti.PostbackStatus = mo.PostbackStatus
		sendNoti.PostbackMessage = mo.PostbackCode
		//sendNoti.TransactionID = notification.Sign
		sendNoti.AffName = mo.AffName
		sendNoti.Msisdn = mo.Msisdn
		sendNoti.Operator = mo.Operator
		sendNoti.Sendtime = mo.Subtime
		sendNoti.NotificationType = "SUB"
		sendNoti.SendData()
	}

	notis := new([]Notification)
	o.QueryTable(NotificationTBName()).OrderBy("-id").All(notis)
	for _, oneNoti := range *notis {
		if oneNoti.NotificationType != "" && oneNoti.NotificationType != "SUB" {
			mo := new(Mo)
			mo.GetMoBySubscriptionID(oneNoti.SubscriptionID)
			sendNoti := new(sendData.SpNotification)
			sendNoti.OfferID = mo.OfferID
			sendNoti.SubscriptionID = mo.SubscriptionID
			sendNoti.ServiceID = mo.ServiceID
			sendNoti.ClickID = mo.ClickID
			sendNoti.CampID = mo.CampID
			sendNoti.PubID = mo.PubID
			sendNoti.PostbackStatus = mo.PostbackStatus
			sendNoti.PostbackMessage = mo.PostbackCode
			sendNoti.TransactionID = oneNoti.Sign
			sendNoti.AffName = mo.AffName
			sendNoti.Msisdn = mo.Msisdn
			sendNoti.Operator = mo.Operator
			sendNoti.Sendtime = oneNoti.Sendtime
			sendNoti.NotificationType = oneNoti.NotificationType
			sendNoti.SendData()
		}
	}
}
