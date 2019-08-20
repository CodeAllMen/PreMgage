package util

import (
	"reflect"
	"sort"
	"time"

	"github.com/MobileCPX/PreNTH/conf"
)

// GetFormatHoursTime 获取当前格式化时间  格式为 2006-01-02 15
func GetFormatHoursTime() string {
	time.LoadLocation("UTC")
	newFormat := time.Now().UTC().Format("2006-01-02 15")
	return newFormat
}

// GetLastMonth 获取上个月月份
func GetLastMonth(num int) string {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	mouth := thisMonth.AddDate(0, num, 0).Format("2006-01")
	return mouth
}

// GetDateList 获取时间列表
func GetDateList(startDate, endDate string) (dateList []string) {
	time.LoadLocation("UTC")
	d, _ := time.ParseDuration("24h")
	start, _ := time.Parse("2006-01-02", startDate)
	for i := 1; ; i++ {
		if start.Format("2006-01-02") <= endDate {
			dateList = append(dateList, start.Format("2006-01-02"))
			start = start.Add(d)
		} else {
			break
		}
	}
	return
}

// GetAffPrice 获取每个转化的价格
func GetAffPrice(date, affName, clickType string) (price float32) {
	price = 6.8
	if affName == "" || affName == "test_affName" {
		price = 0.0
	}
	return
}

// GetOperatorPrice 获取运营商分成价格
func GetOperatorPrice(operator string) (price float32) {
	priceMap := map[string]float32{"20402": 2.199 * 1.17, "20404": 2.484 * 1.17, "20408": 2.365 * 1.17,
		"20416": 2.192 * 1.17}
	return priceMap[operator]
}

// Duplicate 列表去重
func Duplicate(a []string) (ret []string) {
	sort.Strings(a)
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).String())
	}
	return ret
}

// GetServiceType 根据服务短码获取国家及服务类型
func GetServiceType(serviceCode string) string {
	var serviceType string
	switch serviceCode {
	case "NL030070":
		serviceType = conf.NLRedlightvideosName
	case "NL030076":
		serviceType = conf.NLHotvideoName
	case "NL030077":
		serviceType = conf.NLGogamehubName
	case "NL030080":
		serviceType = conf.NLIfunnyName
	case "NL030081":
		serviceType = conf.NLPorn4KName
	}
	return serviceType
}
