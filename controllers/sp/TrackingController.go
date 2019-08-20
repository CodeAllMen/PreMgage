package sp

import (
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreMgage/models/sp"

	"github.com/astaxie/beego/logs"
	"strconv"
)

type TrackingController struct {
	BaseController
}

func (c *TrackingController) InsertAffClick() {
	track := new(sp.AffTrack)
	returnStr := ""
	defer func() {
		if returnStr == "false" {
			track.Update()
		}
	}()
	reqTrack := new(tracking.Track)
	reqTrack, err := reqTrack.BodyToTrack(c.Ctx.Request.Body)
	if err != nil {
		c.StringResult("false")
	}

	track.Track = *reqTrack

	trackID, err := track.Insert()
	if err != nil {
		c.StringResult("false")
	}
	returnStr = strconv.Itoa(int(trackID)) // 返回自增ID

	// 添加判断是否可以订阅条件
	// 获取服务配置
	serviceConf, isExist := c.serviceCofig(track.ServiceID)
	if !isExist {
		track.ClickStatus = tracking.ServiceIsError
		logs.Info("serviceID 不存在：", track.ServiceID)
		returnStr = "false"
		c.StringResult("false")
	}

	track.ServiceCode = serviceConf.ServiceCode
	track.ServiceName = serviceConf.ServiceName
	// 检查是否超过订阅限制
	todaySubNum, err1 := sp.GetTodayMoNum(track.ServiceID)
	if err1 != nil || int(todaySubNum) > serviceConf.LimitSubNum {
		track.ClickStatus = tracking.ExceededSubscriptionLimit
		logs.Info("超过订阅限制: 限制数：%d  订阅数：%d", serviceConf.LimitSubNum, todaySubNum)
		returnStr = "false"
		c.StringResult("false")
	}

	c.StringResult(returnStr)
}
