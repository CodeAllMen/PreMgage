package util

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/MobileCPX/PreCM/conf"
	"github.com/astaxie/beego/logs"
)

//GetFormatTime 获取当前时间及日期
func GetFormatTime() (nowTime, nowDate string) {
	time.LoadLocation("UTC")
	//h, _ := time.ParseDuration("1h")
	nowTime = time.Now().UTC().Format("2006-01-02 15:04:05")
	nowDate = time.Now().UTC().Format("2006-01-02")
	return
}

func Md5Decode(str string) string {
	encodeStr := str + conf.Secretkey
	data := []byte(encodeStr)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	fmt.Println(md5str1)
	return md5str1
}

////GetIpAddress 获取用户ip地址
//func GetIpAddress(r *http.Request) string {
//	hdr := r.Header
//	hdrRealIp := hdr.Get("X-Real-Ip")
//	hdrForwardedFor := hdr.Get("X-Forwarded-For")
//	if hdrRealIp == "" && hdrForwardedFor == "" {
//		return ipAddrFromRemoteAddr(r.RemoteAddr)
//	}
//	if hdrForwardedFor != "" {
//		// X-Forwarded-For is potentially a list of addresses separated with ","
//		parts := strings.Split(hdrForwardedFor, ",")
//		for i, p := range parts {
//			parts[i] = strings.TrimSpace(p)
//		}
//		// TODO: should return first non-local address
//		return parts[0]
//	}
//	return hdrRealIp
//}
//
//func ipAddrFromRemoteAddr(s string) string {
//	idx := strings.LastIndex(s, ":")
//	if idx == -1 {
//		return s
//	}
//	return s[:idx]
//}

// ServiceRegisterRequest 内容站注册
func ServiceRegisterRequest(msisdn, types, period, subID string) (serviceURL string) {
	urls := ""
	coins := "0"
	if types == "register" {
		if period == "game_w" || period == "game_d" {
			serviceURL = "http://www.gogamehub.com"
			urls = fmt.Sprintf("http://www.gogamehub.com/addusername?username=%s&coins=%s&sign=go4movil&subId=%s", msisdn, coins, subID)
		} else if period == "video_w" || period == "video_d" {
			serviceURL = "http://www.redlightvideos.com/gm/es?sign=" + msisdn // 加subID，自动登录
			urls = fmt.Sprintf("http://www.redlightvideos.com/addsubs?phone=%s&sign=go4movil&subId=%s", msisdn, subID)
		}
	} else if types == "addconis" {
		urls = fmt.Sprintf("http://www.gogamehub.com/addusername?username=%s&coins=%s&sign=go4movil&subId=%s", msisdn, coins, subID)
	} else if types == "delete" {
		urls = fmt.Sprintf("http://www.gogamehub.com/addusername?username=%s&coins=%s&sign=go4movil&subId=%s", msisdn, coins, subID)
	}
	resp, err := http.Get(urls)
	if err == nil {
		logs.Info(fmt.Sprintf("HttpRequest Success %s service %s msisdn: %s  subId: %s  conins: %s", types, period, msisdn, subID, coins))
		resp.Body.Close()
	} else {
		logs.Error(fmt.Sprintf("HttpRequest Failed %s service %s msisdn: %s  subId: %s   conins: %s   error: %s ", types, period, msisdn, subID, coins, err.Error()))
	}
	return
}

// GetServiceURL 获取内容站url
func GetServiceURL(serviceID, cmid string) (serviceURL string) {
	switch serviceID {
	case conf.ServiceID:
		serviceURL = "http://www.redlightvideos.com/cm/nl?uiid=" + cmid
	}
	return
}

// GetLPPageURL 获取内容站url
func GetLPPageURL(serviceID string) (serviceURL string) {
	switch serviceID {
	case conf.ServiceID:
		serviceURL = "http://cm.allcpx.com/lp/sub?affName=Self"
	}
	return
}
