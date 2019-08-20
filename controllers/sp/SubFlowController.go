package sp

import (
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreMgage/enums"
	"github.com/MobileCPX/PreMgage/util"
	"github.com/astaxie/beego/httplib"
	"net/url"
	"strconv"

	//"github.com/MobileCPX/PreMgage/httpRequest"
	"github.com/MobileCPX/PreMgage/models/sp"
	"github.com/astaxie/beego/logs"
)

// LPTrackControllers 存储点击
type SubFlowController struct {
	BaseController
}

func (c *SubFlowController) ToAOC() {
	var campSubNum int64
	var err error
	track := new(sp.AffTrack)
	track.ServiceID = c.GetString("sid")
	// 处理传的参数，赋值
	track = c.HandlerParameterToAffTrack(track)

	// 存入点击信息

	logs.Info("track.OfferID", track.OfferID)
	if track.OfferID != 0 {
		trackCookieID := c.Ctx.GetCookie("CK_TRACK")
		if trackCookieID != "" {
			trackCookieMo := new(sp.Mo)
			trackIDInt, err := strconv.Atoi(trackCookieID)
			if err == nil {
				_ = trackCookieMo.GetMoByTrackID(int64(trackIDInt))
				if trackCookieMo.ID != 0 {
					logs.Info("用户已经订阅，跳转到谷歌页面，track_id: ", trackIDInt)
					c.redirect("https://google.com")
				}
			}
		}

		campID := sp.GetCampIDByOfferID(track.OfferID)
		fmt.Println(campID, "!!!!!!!!!!!!!!!!")
		if campID != 0 {
			track.CampID = campID
			mo := new(sp.Mo)
			// 获取今日订阅数量，判断是否超过订阅限制
			campSubNum, err = mo.GetCampTodaySubNum(campID)
			if err != nil {
				c.Ctx.WriteString("false")
				c.StopRun()
			}
			logs.Info("campID: ", campID, " 今日订阅数量： ", campSubNum, " 限制订阅数量：", 50)
			if campSubNum > 50 {
				c.Ctx.WriteString("false")
				c.StopRun()
			}
		} else {
			c.Ctx.WriteString("false")
			c.StopRun()
		}
	}
	trackID, err := track.Insert()

	if err != nil || int(campSubNum) >= enums.DayLimitSub {
		if int(campSubNum) >= enums.DayLimitSub {
			logs.Info(track.ServiceName+" 今日订阅数超过限制 今日订阅: ", campSubNum, " 限制：", enums.DayLimitSub)
		}
		c.Ctx.WriteString("false")
		c.StopRun()
	}

	c.Ctx.SetCookie("CK_TRACK", strconv.Itoa(int(trackID)), 100000000)

	c.redirect("/redirect/aoc/" + strconv.Itoa(int(trackID)))

}

func (c *SubFlowController) InsertAffClick() {
	var campSubNum int64
	var err error
	track := new(sp.AffTrack)
	track.ServiceID = c.GetString("service_id")
	track.ServiceName = c.GetString("service_name")
	// 处理传的参数，赋值
	track = c.HandlerParameterToAffTrack(track)

	// 存入点击信息

	logs.Info("track.OfferID", track.OfferID)
	if track.OfferID != 0 {
		trackCookieID := c.Ctx.GetCookie("CK_TRACK")
		if trackCookieID != "" {
			trackCookieMo := new(sp.Mo)
			trackIDInt, err := strconv.Atoi(trackCookieID)
			if err == nil {
				_ = trackCookieMo.GetMoByTrackID(int64(trackIDInt))
				if trackCookieMo.ID != 0 {
					logs.Info("用户已经订阅，跳转到谷歌页面，track_id: ", trackIDInt)
					c.redirect("https://google.com")
				}
			}
		}

		campID := sp.GetCampIDByOfferID(track.OfferID)
		fmt.Println(campID, "!!!!!!!!!!!!!!!!")
		if campID != 0 {
			track.CampID = campID
			mo := new(sp.Mo)
			// 获取今日订阅数量，判断是否超过订阅限制
			campSubNum, err = mo.GetCampTodaySubNum(campID)
			if err != nil {
				c.Ctx.WriteString("false")
				c.StopRun()
			}
			logs.Info("campID: ", campID, " 今日订阅数量： ", campSubNum, " 限制订阅数量：", 50)
			if campSubNum > 50 {
				c.Ctx.WriteString("false")
				c.StopRun()
			}
		} else {
			c.Ctx.WriteString("false")
			c.StopRun()
		}
	}
	trackID, err := track.Insert()

	if err != nil || int(campSubNum) >= enums.DayLimitSub {
		if int(campSubNum) >= enums.DayLimitSub {
			logs.Info(track.ServiceName+" 今日订阅数超过限制 今日订阅: ", campSubNum, " 限制：", enums.DayLimitSub)
		}
		c.Ctx.WriteString("false")
		c.StopRun()
	}

	c.Ctx.SetCookie("CK_TRACK", strconv.Itoa(int(trackID)), 100000000)
	c.Ctx.WriteString(strconv.Itoa(int(trackID)))
}

func (c *SubFlowController) RedirectAOC() {
	track := new(sp.AffTrack)
	defer func() { // 更新Track表
		if track.ClickStatus != 0 {
			_ = track.Update()
		}
	}()
	trackID := c.Ctx.Input.Param(":trackID")
	trackIntID := c.getIntTrackID(trackID)

	err := track.GetAffTrackByTrackID(int64(trackIntID))
	if err != nil {
		c.redirect("https://google.com")
	}
	track.UserAgent = c.Ctx.Input.UserAgent() //用户设备信息
	track.IP = util.GetIPAddress(c.Ctx.Request)
	_ = track.Update()

	serviceInfo, isExist := c.serviceCofig(track.ServiceID)
	if !isExist {
		track.ClickStatus = tracking.ServiceIsError
		c.redirect("https://google.com")
	}

	redirectAocURL := serviceInfo.GetRedirectAOCURL(track.TrackID)
	if serviceInfo.Version == 2 {
		redirectAocURL = serviceInfo.GetRedirectAOCURLV2(track.TrackID)
	}
	if redirectAocURL == "" {
		track.ClickStatus = tracking.ServiceIsError
		c.redirect("https://uk.google.com")
	} else {
		track.ClickStatus = tracking.ReqAoc
		c.redirect(redirectAocURL)
	}
}

func (c *SubFlowController) GameStartSub() {
	logs.Info("GameStartSub", c.Ctx.Input.URI())
	affTrack := new(sp.AffTrack) // 每次点击存入此次点击的相关数据
	affTrack.AffName = c.GetString("affName")
	affTrack.PubID = c.GetString("pubId")
	affTrack.ProID = c.GetString("proId")
	affTrack.ClickID = c.GetString("clickId")
	affTrack.ServiceID = "889"
	affTrack.ServiceName = "Fit Body"
	trackID, err := affTrack.Insert()

	// 获取今日订阅数量，判断是否超过订阅限制
	todaySubNum, err1 := sp.GetTodayMoNum(affTrack.ServiceID)
	if (err != nil || err1 != nil || int(todaySubNum) >= enums.DayLimitSub) && affTrack.AffName != "" {
		if int(todaySubNum) >= enums.DayLimitSub {
			logs.Info(affTrack.ServiceName+" 今日订阅数超过限制 今日订阅: ", todaySubNum, " 限制：", enums.DayLimitSub)
		}
		c.Ctx.ResponseWriter.ResponseWriter.WriteHeader(404)
		c.StopRun()
	}
	// 获取服务配置信息
	serviceInfo, isExist := c.serviceCofig(affTrack.ServiceID)
	if !isExist {
		logs.Error("Game 服务名称不存在，请检查服务信息，servideName: ", affTrack.ServiceName)
		c.Ctx.ResponseWriter.ResponseWriter.WriteHeader(404)
		c.StopRun()
	}

	// 获取跳转到AOC页面 URL
	redirectAocURL := serviceInfo.GetRedirectAOCURL(trackID)
	if redirectAocURL == "" {
		c.redirect("https://uk.google.com")
	} else {
		c.redirect(redirectAocURL)
	}
}

func (c *SubFlowController) Fit8TubeStartSub() {
	logs.Info("GameStartSub", c.Ctx.Input.URI())
	affTrack := new(sp.AffTrack) // 每次点击存入此次点击的相关数据
	affTrack.AffName = c.GetString("affName")
	affTrack.PubID = c.GetString("pubId")
	affTrack.ProID = c.GetString("proId")
	affTrack.ClickID = c.GetString("clickId")
	affTrack.ServiceID = "892"
	affTrack.ServiceName = "Gold Finger"
	trackID, err := affTrack.Insert()

	// 获取今日订阅数量，判断是否超过订阅限制
	todaySubNum, err1 := sp.GetTodayMoNum(affTrack.ServiceID)
	if (err != nil || err1 != nil || int(todaySubNum) >= enums.DayLimitSub) && affTrack.AffName != "" {
		if int(todaySubNum) >= enums.DayLimitSub {
			logs.Info(affTrack.ServiceName+" 今日订阅数超过限制 今日订阅: ", todaySubNum, " 限制：", enums.DayLimitSub)
		}
		c.Ctx.ResponseWriter.ResponseWriter.WriteHeader(404)
		c.StopRun()
	}
	// 获取服务配置信息
	serviceInfo, isExist := c.serviceCofig(affTrack.ServiceID)
	if !isExist {
		logs.Error("Game 服务名称不存在，请检查服务信息，servideName: ", affTrack.ServiceName)
		c.Ctx.ResponseWriter.ResponseWriter.WriteHeader(404)
		c.StopRun()
	}

	// 获取跳转到AOC页面 URL
	redirectAocURL := serviceInfo.GetRedirectAOCURL(trackID)
	if redirectAocURL == "" {
		c.redirect("https://uk.google.com")
	} else {
		c.redirect(redirectAocURL)
	}
}

// 订阅结果返回
func (c *SubFlowController) SubResultReturn() {
	// V1   http://mguk.foxseek.com/sub/result?mig_status=active&migid=0536e03d&mig_optin=0&s=889&t=60&m=07488555687&msisdn=07488555687&mig_sid=2482024&sig=904e114a0bb7bb24b31871af7723af301ad31343&operator=4
	// V2  http://merchant.example.com/redirection?service_id=4&merchant_ref=89ug98uaf9ek09i90i&status=successful&sig=<generated_sig>

	// 初始化数据
	subResult := new(sp.SubResult)
	logs.Info("SubResultReturn", c.Ctx.Input.URI())
	subResult.TrackID = c.GetString("t")
	if subResult.TrackID == "" {
		subResult.TrackID = c.GetString("merchant_ref")
	}

	track := new(sp.AffTrack)
	trackIntID := c.getIntTrackID(subResult.TrackID)

	err := track.GetAffTrackByTrackID(int64(trackIntID))
	if err != nil {
		c.redirect("https://google.com")
	}
	serviceInfo, isExist := c.serviceCofig(track.ServiceID)
	if !isExist {
		track.ClickStatus = tracking.ServiceIsError
		c.redirect("https://google.com")
	}

	if serviceInfo.Version == 2 {
		subResult.SubStatus = c.GetString("status")
	} else {
		subResult.SubscriptionID = c.GetString("mig_sid")
		subResult.SubStatus = c.GetString("mig_status")
		subResult.UserID = c.GetString("migid")
		subResult.MigOptin = c.GetString("mig_optin")
		subResult.ErrorDesc = c.GetString("mig_error")

		subResult.ServiceID = c.GetString("s")
		subResult.Msisdn = c.GetString("msisdn")
		subResult.Operator = c.GetString("operator")
	}

	subResult.Sign = c.GetString("sig")

	// 通过s 获取服务信息
	//service := c.getService(subResult.ServiceID)

	if subResult.SubStatus == "successful" || subResult.SubStatus == "pending" || subResult.SubStatus == "active" {
		c.RegisteredService(subResult.ServiceID, subResult.Msisdn)
		_ = subResult.Insert()

		contentURL := serviceInfo.ContentURL + subResult.TrackID
		// 生成随机id
		randomStr, err := httplib.Get("http://offer.globaltraffictracking.com/sub_success/req?url=" +
			url.QueryEscape(contentURL)).String()
		if err == nil && len(randomStr) > 3 {
			if randomStr[:2] == "AA" {
				//订阅成功记录订阅ID
				c.redirect("http://offer.globaltraffictracking.com/sub_track/" + randomStr + "?sub=" + subResult.TrackID)
			}
		}
		c.redirect(serviceInfo.ContentURL + subResult.TrackID)
	} else {
		c.redirect("http://www.google.com")
	}
}
