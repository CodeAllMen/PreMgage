package backData

import (
	"encoding/json"
	"github.com/MobileCPX/PreBaseLib/splib/admindata"
	"github.com/MobileCPX/PreMgage/models/sp"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

//
//
func SendNotification() {
	o := orm.NewOrm()
	notifis := new([]sp.Notification)
	notificationTypeList := []string{"SUCCESS_MT", "FAILED_MT", "UNSUB"}
	o.QueryTable("notification").Filter("sendtime__gt", "2019-07-01").Filter("notification_type__in", notificationTypeList).OrderBy("id").All(notifis)

	for _, one := range *notifis {

		sendNoti := new(admindata.Notification)
		notificationType := one.NotificationType

		sendNoti.ServiceID = one.ServiceID
		sendNoti.Operator = one.Operator

		sendNoti.SubscriptionID = one.SubscriptionID

		sendNoti.TransactionID = one.Sign

		sendNoti.Sendtime = one.Sendtime

		sendNoti.NotificationType = notificationType

		data, _ := json.Marshal(sendNoti)
		logs.Info(string(data))
		sendNoti.SendData(admindata.PROD)

	}

}
