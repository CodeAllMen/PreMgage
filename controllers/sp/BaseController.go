package sp

import (
	"github.com/MobileCPX/PreBaseLib/common"
	"github.com/MobileCPX/PreMgage/enums"
	"github.com/MobileCPX/PreMgage/models"
	"github.com/MobileCPX/PreMgage/models/sp"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
	"time"
)

type BaseController struct {
	common.BaseController
}

// CheckError 检查是否有错 msg 定义日志信息
func (c *BaseController) CheckError(err error, errorCode enums.ErrorCode, msg ...string) {
	if err != nil {
		// 打印日志信息
		if len(msg) != 0 {
			logs.Error(msg, " ERROR: ", err.Error())
		}
		switch errorCode {
		case enums.RedirectGoogle:
			c.redirect("https://wwww.google.com")
		case enums.Error502:
			c.Ctx.ResponseWriter.WriteHeader(502)
			c.StopRun()
		default:
			c.redirect("https://wwww.google.com")
		}
	}
}

func (c *BaseController) NewInsertMo(notification *sp.Notification, affTrack *sp.AffTrack) (*sp.Mo, string) {
	mo := new(sp.Mo)
	notificationType := ""
	isExist := mo.CheckSubIDIsExist(notification.SubscriptionID)

	// 判断用户是否已经存在
	if !isExist {
		mo = new(sp.Mo)
		// 初始化MO数据
		mo.InitNewSubMO(notification, affTrack)
		// 查询次网盟今天的订阅数及postback回传数，根据概率判断次数书是否回传
		subNum, postbackNum := mo.GetOfferTodaySubInfo()
		// 根据概率判断次数书是否回传
		isSuccess, code, payout := sp.StartPostback(mo, subNum, postbackNum)
		mo.PostbackCode = code

		//if mo.AffName != "" {
		//	// 有转化后发邮件
		//	util.BeegoEmail("", "Mgage 英国  有新的转化", "网盟名称： "+mo.AffName,
		//		[]string{"wangangui@mobilecpx.com"})
		//	//[]string{"tengjiaqing@mobilecpx.com", "wangangui@mobilecpx.com"})
		//}
		// 判断是否回传成功
		if isSuccess {
			mo.Payout = payout
			mo.PostbackStatus = 1
		}
		_ = mo.InsertNewMo()
		notificationType = "SUB"
	}
	return mo, notificationType
}

func (c *BaseController) serviceCofig(serviceID string) (*sp.ServiceInfo, bool) {
	serviceCofig, isExist := sp.ServiceData[serviceID]
	return &serviceCofig, isExist
}

// setCookie
func (c *BaseController) setCookie(trackID string) string {
	// 获取cookie
	userId, ok := c.GetSecureCookie("user_cookie", "8A66b76dbd3759445fe924d28a5F6856")
	if !ok {
		userId = "PinkCity__" + trackID + "_" + "1"
	} else {
		userIdList := strings.Split(userId, "_")
		if len(userIdList) != 3 {
			userId = "PinkCity__" + trackID + "_" + "1"
		} else {
			vistTimes, err := strconv.Atoi(userIdList[2])
			if err != nil {
				c.Ctx.ResponseWriter.ResponseWriter.WriteHeader(404)
				c.StopRun()
			} else {
				userId = userIdList[0] + "_" + userIdList[1] + "_" + strconv.Itoa(vistTimes+1)
			}
		}
	}
	// 设置cookie
	c.SetSecureCookie("user_cookie", "8A66b76dbd3759445fe924d28a5F6856", userId, 61622400*time.Second)
	return userId
}

func (c *BaseController) redirect(url string) {
	if url == "" {
		url = "http://google.com"
	}
	c.Redirect(url, 302)
	c.StopRun()
}

func (c *BaseController) jsonResult(code enums.JsonResultCode, msg string, obj interface{}) {
	r := &models.JsonResult{code, msg, obj}
	c.Data["json"] = r
	c.ServeJSON()
	c.StopRun()
}

// 分割requestID to trackID  1819_sub_1550768968
func (c *BaseController) splitReuestIDToTrackID(requestID string) (trackID string) {
	result := strings.Split(requestID, "_")
	if len(result) == 3 {
		trackID = result[0]
	}
	return
}

func (c *BaseController) getService(serviceID string) *sp.ServiceInfo {
	service, isExist := sp.GetServerConfByServiceID(serviceID)
	if !isExist {
		c.redirect("https://google.com")
	}
	return service
}

func (c *BaseController) RegisteredService(serviceID, msisdn string) {
	serviceConfig, isExist := c.serviceCofig(serviceID)
	if isExist {
		registerURL := strings.Replace(serviceConfig.RegisterURL, "msisdn", msisdn, -1)
		result, _ := httplib.Get(registerURL).String()
		if result == "OK" {
			logs.Info("订阅后注册账号成功 serviceID :", serviceID, " msisdn: ", msisdn)
		} else {
			logs.Error("订阅后注册账失败 serviceID :", serviceID, " msisdn: ", msisdn)
		}
	}
}

func (c *BaseController) HandlerParameterToAffTrack(track *sp.AffTrack) *sp.AffTrack {
	// 获取网盟名称
	track.AffName = c.GetString("affName")
	if track.AffName == "" {
		track.AffName = c.GetString("aff")
		if track.AffName == "" {
			track.AffName = "aff_name"
		}
	}

	// 获取offerID
	offerID, _ := c.GetInt("offer")
	if offerID == 0 {
		offerID, _ = c.GetInt("of")
		if offerID == 0 {
			offerID, _ = c.GetInt("f")
		}
	}
	track.OfferID = offerID

	// 获取子渠道
	track.PubID = c.GetString("pubId")
	if track.PubID == "" {
		track.PubID = c.GetString("pub")
		if track.PubID == "" {
			track.PubID = c.GetString("p")
		}
	}

	// 获取click_id
	track.ClickID = c.GetString("clickId")
	if track.ClickID == "" {
		track.ClickID = c.GetString("click")
		if track.ClickID == "" {
			track.ClickID = c.GetString("c")
		}
	}

	// 获取 服务id
	track.ProID = c.GetString("pro")
	if track.ProID == "" {
		track.ProID = c.GetString("proId")
	}

	// 获取其他参数
	track.Other = c.GetString("ot")
	if track.Other == "" {
		track.Other = c.GetString("o")
	}

	// 如果offerID 不为空，则根据offerid 查询对应的网盟名称
	if track.OfferID != 0 {
		track.AffName = sp.GetAffNameByOfferID(int64(offerID))
	}

	track.Refer = c.GetString("rf")
	return track
}

func (c *BaseController) getOperatorByCampID() {

}

func (c *BaseController) getIntTrackID(trackStrID string) int {
	trackIntID, err := strconv.Atoi(trackStrID)
	if err != nil {
		c.redirect("http://google.com")
	}
	return trackIntID
}
