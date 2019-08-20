package sp

import (
	"fmt"
	"github.com/MobileCPX/PreMgage/util"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// Config 内容站配置
type Config struct {
	Service map[string]ServiceInfo
}

// 完整的服务配置
type ServiceInfo struct {
	ServiceCode     string `yaml:"service_code" url:"s"`
	ServiceID       string `yaml:"service_id"` // 服务ID 对应Mgage 的Shop ID
	ServiceName     string `yaml:"service_name"`
	MerchantID      string `yaml:"merchant_id"`
	Version         int    `yaml:"version"`                  // API 版本
	ItemDescription string `yaml:"item_description" url:"d"` // AOC 页面的顶部文本描述
	Sign            string `yaml:"sign" url:"sig"`           // 服务的Authorisation Signature
	Amount          string `yaml:"amount" url:"ra"`          // 扣费金额
	Period          string `yaml:"period" url:"rp"`          // 扣费周期 单位是秒  天、周、月
	ContractLength  string `yaml:"contract_length" url:"cl"` // 设置合同长度，单位秒数，如果用户退订了，可以通过设置无限时长来防止重复订阅
	ApiKey          string `yaml:"api_key"`

	AOCURL          string `yaml:"aoc_url"`
	ContentURL      string `yaml:"content_url" `
	ReturnURL       string `yaml:"return_url" `
	RegisterURL     string `yaml:"register_url"`
	LpURL           string `yaml:"lp_url" `
	WelcomePageURL  string `yaml:"welcome_page_url" `
	UnsubURL        string `yaml:"unsub_result_url" `
	ShareKey        string `yaml:"share_key"`
	UnsubSuccessURL string `yaml:"unsub_success_url"`
	UnsubFailedURL  string `yaml:"unsub_failed_url"`
	LimitSubNum     int    `yaml:"limit_sub_num"`
	CampID          int    `yaml:"camp_id"`
}

var ServiceData = make(map[string]ServiceInfo)

func InitSetServiceConfig() {
	filename, _ := filepath.Abs("resource/config/conf.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	config := new(Config)
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		panic(err)
	}
	ServiceData = config.Service
	logs.Info("service_CONFIG", ServiceData)
}

// 订阅请求跳转到AOC页面的结构体
type redirectAOC struct {
	//baseConfig // 基础配置信息
	ItemDescription string `yaml:"item_description" url:"d"` // AOC 页面的顶部文本描述
	ServiceID       string `yaml:"service_id" url:"s"`       // 服务ID 对应Mgage 的Shop ID
	Sign            string `yaml:"sign" url:"sig"`           // 服务的Authorisation Signature
	Amount          string `yaml:"amount" url:"ra"`          // 扣费金额
	Period          string `yaml:"period" url:"rp"`          // 扣费周期 单位是秒  天、周、月
	ContractLength  string `yaml:"contract_length" url:"cl"` // 设置合同长度，单位秒数，如果用户退订了，可以通过设置无限时长来防止重复订阅

	TransactionID int64 `url:"t"` // 交易ID
}

// API 版本1
func (service *ServiceInfo) GetRedirectAOCURL(trackID int64) (redirectAocURL string) {
	// 初始化请求数据

	redirectParm := "cl=" + service.ContractLength + "&d=" + service.ItemDescription + "&ra=" +
		service.Amount + "&rp=" + service.Period + "&s=" + service.ServiceCode + "&t=" +
		strconv.Itoa(int(trackID))
	encodeStr := service.ShareKey + ":" + redirectParm + ":" + service.ShareKey
	sign := util.SHA1Encode(encodeStr)
	logs.Info(encodeStr)

	redirectAocURL = service.AOCURL + "?" + redirectParm + "&sig=" + sign
	fmt.Println(redirectAocURL)
	return
}

// API 版本2
func (service *ServiceInfo) GetRedirectAOCURLV2(trackID int64) (redirectAocURL string) {
	// 初始化请求数据
	redirectParm := "amount=" + service.Amount + "&description=" + service.ItemDescription + "&merchant_id=" +
		service.MerchantID + "&merchant_ref=" + strconv.Itoa(int(trackID)) + "&service_id=" + service.ServiceCode + "&subscription_length=weekly"

	sign := util.HmacSHA1(service.ShareKey, redirectParm)
	logs.Info(sign)

	redirectAocURL = service.AOCURL + "?" + redirectParm + "&sig=" + sign
	return
}

// 通过查询Sign 注册服务
func RegistereServerBySign(serviceID, userNmae string) {
	service := ServiceData[serviceID]
	fmt.Println(service.RegisterURL)
	registerURL := strings.Replace(service.RegisterURL, "{user_name}", userNmae, -1)
	_, err := httplib.Get(registerURL).String()
	if err != nil {
		// error
		logs.Error("用户订阅成功后注册账号失败 ERROR: ", err.Error())
	}
}

//func GetServiceInfoBySign(sign string) *ServiceInfo {
//	for _, service := range ServiceData {
//		if service.Sign == sign {
//			return &service
//		}
//	}
//	return nil
//}

func GetServerConfByServiceID(serviceID string) (*ServiceInfo, bool) {
	conf, isExist := ServiceData[serviceID]
	return &conf, isExist
}
