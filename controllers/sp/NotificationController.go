package sp

import (
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/admindata"
	"github.com/MobileCPX/PreMgage/models/sp"
	"github.com/MobileCPX/PreMgage/util"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
)

type NotificationController struct {
	BaseController
}

///sp/notification?mig_stage=UserMsisdnNotFound&sig=ef4b95bf6ef9bca287bec3fff3559300ba213590&s=889&t=64
// 成功订阅通知   /sp/notification?period_from=2019-03-21T11%3A36%3A57Z&s=889&t=60&period_to=2019-03-28T11%3A36%3A57Z&migid=0536e03d&sig=18d6a1413c07eaf019e2f2451116bc82e412e19b&mig_sid=2482024&mig_stage=CustomerSubscriptionCreated
// 扣费通知  /sp/notification?period_from=2019-03-21T11%3A36%3A57Z&s=889&t=60&period_to=2019-03-28T11%3A36%3A57Z&migid=0536e03d&sig=fe0cc667820ff0aac813d4e82e9673be59fa3795&mig_stage=CustomerDirectBillingSuccessful^

func (c *NotificationController) Get() {
	logs.Info("notification", c.Ctx.Input.URI())
	notification := new(sp.Notification)
	nowTime, _ := util.GetNowTimeFormat()
	trackID := c.GetString("t")
	if trackID == "" {
		trackID = c.GetString("merchant_ref")
	}

	track := new(sp.AffTrack)
	trackIDInt, err := strconv.Atoi(trackID)
	if err == nil {
		_ = track.GetAffTrackByTrackID(int64(trackIDInt))
	}
	notification.TrackID = trackID

	service, _ := sp.GetServerConfByServiceID(track.ServiceID)
	if service.Version == 2 {
		notification.TransactionID = c.GetString("nid")

		notification.UserID = c.GetString("msisdn_alias") // 用户唯一标识ID
		notification.SubStage = c.GetString("event_name")
		notification.Operator = c.GetString("operator")
		notification.CurrentPeriodStart = c.GetString("current_period_start")
		notification.CurrentPeriodEnd = c.GetString("current_period_end")

	} else {

		notification.UserID = c.GetString("migid") // 用户唯一标识ID
		notification.SubStage = c.GetString("mig_stage")
		notification.Sign = c.GetString("sig")
		notification.ServiceID = c.GetString("s")
	}

	if c.GetString("mig_sid") != "" {
		notification.SubscriptionID = c.GetString("mig_sid")
	}

	mo := new(sp.Mo)
	// 不是订阅通知，需要先查询出订阅信息
	if notification.SubStage == "CustomerDirectBillingSuccessful" ||
		notification.SubStage == "CustomerSubscriptionUnsubscribed" ||
		notification.SubStage == "CustomerSubscriptionUnsubscribeRequested" {

		err = mo.GetMoByTrackID(int64(trackIDInt))
		if mo.ID == 0 { // 如果没有查询到mo数据，说明扣费通知早于听月通知，等10秒之后再次查询mo数据
			logs.Info("扣费通知没有查到MO 数据", trackID)
			for i := 0; i <= 10; i++ {
				err = mo.GetMoByTrackID(int64(trackIDInt))
				logs.Info("第")
				if mo.ID != 0 {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}

		// 找不到需要重新插入订阅数据
		if err != nil {
			subResult := new(sp.SubResult)
			// 订阅结果表通过track_id 查询SubscriptionID
			_, _ = subResult.GetSubResultByTrackID(notification.TrackID)
			fmt.Println(subResult)
			if subResult.ID != 0 {
				notification.SubscriptionID = subResult.SubscriptionID
				mo.SubscriptionID = subResult.SubscriptionID
				mo.ServiceID = track.ServiceID
				if mo.SubscriptionID == "" {
					mo.SubscriptionID = mo.ServiceID + "_" + trackID
				}

				mo.ServiceName = track.ServiceName
				mo.SubStatus = 1
				_ = mo.InsertNewMo()
			}
		}
		notification.SubscriptionID = mo.SubscriptionID
	}

	switch notification.SubStage {
	case "CustomerSubscriptionCreated": // 订阅成功
		subResult := new(sp.SubResult)
		_, _ = subResult.GetSubResultByTrackID(notification.TrackID)
		if notification.SubscriptionID == "" {
			notification.SubscriptionID = track.ServiceID + "-" + trackID
		}

		notification.Msisdn = subResult.Msisdn
		mo, notification.NotificationType = c.NewInsertMo(notification, track)
	case "CustomerDirectBillingSuccessful": // 扣费成功
		notification.NotificationType, _ = mo.SuccessMTUpdateMO()
	case "CustomerDirectBillingFailed": // 扣费失败
		notification.NotificationType, _ = mo.FailedMTUpdateMo()
	case "CustomerSubscriptionUnsubscribed", "CustomerSubscriptionUnsubscribeRequested": // 退订
		notification.NotificationType, _ = mo.UnsubUpdateMo()
	}

	_ = notification.Insert()
	logs.Info(*mo)

	if notification.NotificationType != "" && mo.ServiceID != "" {
		fmt.Println("1111111111111111")
		sendNoti := new(admindata.Notification)

		if mo.OfferID != 0 {
			postback := new(sp.Postback)
			_ = postback.CheckOfferID(mo.OfferID)
			sendNoti.PromoterID = postback.PromoterID
		}

		sendNoti.PostbackPrice = mo.PostbackPrice

		sendNoti.OfferID = mo.OfferID
		sendNoti.SubscriptionID = mo.SubscriptionID
		sendNoti.ServiceID = mo.ServiceID
		sendNoti.ClickID = mo.ClickID
		sendNoti.Msisdn = mo.Msisdn
		sendNoti.CampID = mo.CampID
		sendNoti.PubID = mo.PubID
		sendNoti.PostbackStatus = mo.PostbackStatus
		sendNoti.PostbackMessage = mo.PostbackCode
		sendNoti.TransactionID = notification.Sign
		sendNoti.AffName = mo.AffName
		if sendNoti.AffName == "" {
			sendNoti.AffName = "未知"
		}
		sendNoti.Operator = mo.Operator

		sendNoti.Sendtime = nowTime
		sendNoti.NotificationType = notification.NotificationType
		sendNoti.SendData(admindata.PROD)
	}

	c.Ctx.WriteString("OK")
}
