package main

import (
	"github.com/MobileCPX/PreBaseLib/splib/click"
	_ "github.com/MobileCPX/PreMgage/initial"
	"github.com/MobileCPX/PreMgage/models/sp"
	_ "github.com/MobileCPX/PreMgage/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron"
)

func init() {
	sp.InitSetServiceConfig()

}

func main() {
	//backData.SendMo()
	//backData.SendNotification()

	//sp.UpdateMoAndNotificationTable()
	//sp.SendAdminData()
	//SendClickDataToAdmin()

	logs.SetLogger(logs.AdapterFile, `{"filename":"/home/ubuntu/Logs/mgage/uk/mgage_uk.log","level":6,"maxlines":100000000,"daily":true,"maxdays":10000}`)
	logs.Async(1e3)
	logs.EnableFuncCallDepth(true)

	//sp.UpdateMoAndNotificationTable()
	task()

	beego.Run()
}

// 定时任务
func task() {
	cr := cron.New()

	cr.AddFunc("0 24 */1 * * ?", SendClickDataToAdmin) // 一个小时存一次点击数据并且发送到Admin

	cr.Start()
}

func SendClickDataToAdmin() {
	sp.InsertHourClick()

	for _, service := range sp.ServiceData {
		click.SendHourData(service.CampID, click.PROD) // 发送有效点击数据
	}

}
