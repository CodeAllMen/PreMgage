package sp

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func UpdateMoAndNotificationTable() {
	o := orm.NewOrm()
	mos := new([]Mo)
	_, _ = o.QueryTable(MoTBName()).All(mos)
	for _, mo := range *mos {
		if mo.TrackID == 0 {
			var track AffTrack
			if mo.ClickID != "" {
				_ = o.QueryTable("aff_track").Filter("click_id", mo.ClickID).One(&track)
				if track.TrackID != 0 {
					mo.TrackID = track.TrackID
					_, _ = o.Update(&mo)
				}
			}
		}
	}
	notifs := new([]Notification)
	_, _ = o.QueryTable("notification").All(notifs)

	for _, oneNotif := range *notifs {
		if oneNotif.TrackID != "" && oneNotif.SubscriptionID == "" {
			var mo Mo
			_ = o.QueryTable("mo").Filter("track_id", oneNotif.TrackID).One(&mo)
			logs.Info("mo:", mo.SubscriptionID)
			if mo.SubscriptionID != "" {
				oneNotif.SubscriptionID = mo.SubscriptionID
				o.Update(&oneNotif)
			}
		}
	}
}

//func SendData() {
//	o := orm.NewOrm()
//	notifs := new([]Notification)
//	o.QueryTable("notification").Filter("")
//}
