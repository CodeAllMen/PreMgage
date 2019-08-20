package routers

import (
	"github.com/MobileCPX/PreMgage/controllers"
	"github.com/MobileCPX/PreMgage/controllers/sp"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/to/aoc", &sp.SubFlowController{}, "Get:ToAOC")

	beego.Router("/", &controllers.MainController{})
	// 订阅、退订、续订通知
	beego.Router("/sp/notification", &sp.NotificationController{})

	// 订阅流程
	// 第一步，订阅请求
	beego.Router("/track/returnid", &sp.SubFlowController{}, "Get:InsertAffClick") // 存点击
	beego.Router("/redirect/aoc/:trackID", &sp.SubFlowController{}, "Get:RedirectAOC")

	beego.Router("/game/start/sub", &sp.SubFlowController{}, "Get:GameStartSub")
	// 第二步 订阅结果通知
	beego.Router("/sub/result", &sp.SubFlowController{}, "Get:SubResultReturn")

	// 退订流程
	beego.Router("/unsub/subID/:serviceID", &sp.UnsubController{}, "Get:SubIDUnsub")
	// 通过电话号码退订
	beego.Router("/unsub/msisdn/:serviceID", &sp.UnsubController{}, "Get:UnsubByMsisdn")
	beego.Router("/unsub/subid/:serviceID", &sp.UnsubController{}, "Get:SubIDUnsub")

	// 设置 postback
	beego.Router("/set/postback", &sp.SetPostbackController{})

	// 存点击
	beego.Router("/aff/click", &sp.TrackingController{}, "Post:InsertAffClick")

}
