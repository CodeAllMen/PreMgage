package sp

import (
	"errors"
	"fmt"
	"github.com/MobileCPX/PreMgage/util"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

// Postback 网盟信息
type Postback struct {
	ID           int     `orm:"pk;auto;column(id)" json:"-"`                            // 自增ID
	CreateTime   string  `orm:"column(create_time)" json:"-"`                           // 添加时间
	UpdateTime   string  `orm:"column(update_time)" json:"-"`                           // 更新时间
	DayCap       int     `orm:"column(day_cap)" json:"day_cap"`                         // 更新时间
	AffName      string  `orm:"column(aff_name);size(30)" json:"aff_name"`              // 网盟名称
	PostbackURL  string  `orm:"column(postback_url);size(300)" json:"postback_url"`     // postback URL
	PostbackRate int     `orm:"column(postback_rate);default(50)" json:"postback_rate"` // 回传概率
	Payout       float32 `orm:"column(Payout)" json:"payout"`                           // 转化单价
	PromoterName string  `orm:"column(promoter_name)" json:"promoter_name"`             // 外放人
	PromoterID   int     `orm:"column(promoter_id)" json:"promoter_id"`                 // 外放人 ID
	CampID       int     `orm:"column(camp_id)" json:"camp_id"`                         // CampID
	OfferID      int     `orm:"column(offer_id)" json:"offer_id"`
}

// StartPostback 订阅成功后向网盟回传订阅数据
// 请求 todaySubNum 该网盟今日订阅数，  todayPostbackNum 该网盟今日回传数   根据这两个算概率，是否回传
// 返回数据 isSuccess 是否回传   code 网络请求的返回code   payout  请求成功后的花费
func StartPostback(mo *Mo, todaySubNum, todayPostbackNum int64) (isSuccess bool, code string, payout float32) {
	// postback, err := getPostbackInfoByAffName(mo.AffName, mo.ServiceName)

	postback, err := getPostbackInfoByOfferID(mo.OfferID, mo.AffName, mo.ServiceID)
	if err != nil {
		return
	}

	isPostback := postback.CheckTodayPostbackStatus(todaySubNum, todayPostbackNum)
	if isPostback {
		isSuccess, code = postback.PostbackRequest(mo)
		payout = postback.Payout
	}
	return
}

func getPostbackInfoByOfferID(offerID int, affName, serviceName string) (*Postback, error) {
	postback := new(Postback)
	o := orm.NewOrm()
	if offerID != 0 {
		err := o.QueryTable("postback").Filter("offer_id", offerID).One(postback)
		if err != nil {
			logs.Error("用户订阅成功，但是没有找到此网盟 ", affName, "OfferID", offerID)
			util.BeegoEmail(serviceName, "没有找到此 "+affName+"信息", affName+" postback回传失败", []string{})
		}
		return postback, err
	}
	return postback, errors.New("网盟为空")
}

func getPostbackInfoByAffName(affName, serviceName string) (*Postback, error) {
	postback := new(Postback)
	o := orm.NewOrm()
	if affName != "" {
		err := o.QueryTable("postback").Filter("aff_name", affName).One(postback)

		if err != nil {
			logs.Error("用户订阅成功，但是没有找到此网盟 ", affName)
			util.BeegoEmail(serviceName, "没有找到此 "+affName+"信息", affName+" postback回传失败", []string{})
		}
		return postback, err
	}
	return postback, errors.New("网盟为空")
}
func GetPostbackInfoByAffName(affName, serviceName string) (*Postback, error) {
	postback := new(Postback)
	o := orm.NewOrm()
	if affName != "" {
		err := o.QueryTable("postback").Filter("aff_name", affName).One(postback)

		if err != nil {
			logs.Error("用户订阅成功，但是没有找到此网盟 ", affName)
			util.BeegoEmail(serviceName, "没有找到此 "+affName+"信息", affName+" postback回传失败", []string{})
		}
		return postback, err
	}
	return postback, errors.New("网盟为空")
}

func (postback *Postback) CheckTodayPostbackStatus(todaySubNum, todayPostbackNum int64) (isPostback bool) {
	defer logs.Info("postbakck 状态 ", isPostback)
	if todaySubNum == 0 {
		isPostback = true
		return
	}
	currentRate := float32(todayPostbackNum) / float32(todaySubNum)
	if currentRate > float32(postback.PostbackRate)/float32(100) {
		isPostback = false
	} else {
		isPostback = true
	}
	return
}

func (postback *Postback) PostbackRequest(mo *Mo) (isSuccess bool, code string) {
	postbackURL := postback.PostbackURL
	timestamp := time.Now().Unix()
	postbackURL = strings.Replace(postbackURL, "{click_id}", mo.ClickID, -1)
	postbackURL = strings.Replace(postbackURL, "##clickid##", mo.ClickID, -1)
	postbackURL = strings.Replace(postbackURL, "{pro_id}", mo.ProID, -1)
	postbackURL = strings.Replace(postbackURL, "{other}", mo.ProID, -1)
	postbackURL = strings.Replace(postbackURL, "{pub_id}", mo.PubID, -1)
	postbackURL = strings.Replace(postbackURL, "##pub_id##", mo.ClickID, -1)
	postbackURL = strings.Replace(postbackURL, "{operator}", mo.Operator, -1)
	postbackURL = strings.Replace(postbackURL, "{auto}", strconv.Itoa(int(timestamp)), -1)
	postbackURL = strings.Replace(postbackURL, "##auto_id##", strconv.Itoa(int(timestamp)), -1)
	postbackURL = strings.Replace(postbackURL, "{payout}", fmt.Sprintf("%f", postback.Payout), -1)
	postResult, err := httplib.Get(postbackURL).String()
	if err == nil {
		// postback 成功
		isSuccess = true
		logs.Info("postback URL: ", postbackURL, " CODE: ", code)
	} else {
		logs.Error("postback Error , msisdn : " + mo.Msisdn + " aff_name : " + mo.AffName + " error " + err.Error())
	}
	code = postResult
	return
}

func GetAffNameByOfferID(offerID int64) string {
	o := orm.NewOrm()
	postback := new(Postback)
	err := o.QueryTable("postback").Filter("offer_id", offerID).One(postback)
	if err != nil {
		logs.Error("GetAffNameByOfferID 错误，offerID：", offerID, " ERROR: ", err.Error())
	}
	return postback.AffName
}

func (postback *Postback) InsertPostback() error {
	o := orm.NewOrm()
	postback.CreateTime, _ = util.GetNowTimeFormat()
	_, err := o.Insert(postback)
	if err != nil {
		logs.Error("Postback InsertPostback ERROR:", err.Error(), postback)
	}
	return err
}

func (postback *Postback) CheckOfferID(offerID int) error {
	o := orm.NewOrm()
	return o.QueryTable(PostbackTBName()).Filter("offer_id", offerID).One(postback)

}

func GetCampIDByOfferID(offerID int) int {
	o := orm.NewOrm()
	postback := new(Postback)
	err := o.QueryTable(PostbackTBName()).Filter("offer_id", offerID).One(postback)
	if err != nil {
		logs.Error("GetCampIDByOfferID  通过offerId 查询 postback失败,offerID: ", offerID)
		return 0
	}
	return postback.CampID
}

// CheckOfferIDIsExist 检查offer_id是否已经存在
func (postback *Postback) CheckOfferIDIsExist(offerID int) error {
	o := orm.NewOrm()
	err := o.QueryTable("postback").Filter("offer_id", offerID).One(postback)
	if err != nil {
		logs.Error("Postback CheckOfferIDIsExist  ERROR, 检查offer_id是否已经存在 失败")
	}

	return err
}

// Update  更新Postback
func (postBack *Postback) Update() error {
	o := orm.NewOrm()
	postBack.UpdateTime, _ = util.GetFormatTime()
	_, err := o.Update(postBack)
	if err != nil {
		logs.Error("Postback Update 插入postback数据成功")
	}
	return err
}

// Insert  插入Postback
func (postBack *Postback) Insert() error {
	o := orm.NewOrm()
	postBack.CreateTime, _ = util.GetFormatTime()
	_, err := o.Insert(postBack)
	if err != nil {
		logs.Error("Postback Insert 插入postback数据失败，", err.Error())
	}
	return err
}
