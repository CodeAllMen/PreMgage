package sp

import (
	"github.com/MobileCPX/PreMgage/util"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type SubResult struct {
	ID             int64  `orm:"pk;auto;column(id)"`
	Sendtime       string `orm:"column(sendtime)"`
	SubscriptionID string `orm:"column(subscription_id)"`
	SubStatus      string `orm:"column(sub_status)"`
	UserID         string `orm:"column(user_id)"`
	MigOptin       string `orm:"column(mig_optin)"`
	Sign           string `orm:"column(sign)"`
	ServiceID      string `orm:"column(service_id)"`
	TrackID        string `orm:"column(track_id)"`
	ErrorDesc      string `orm:"column(error_desc)"`
	Msisdn         string
	Operator       string
}

func SubResultTBName() string {
	return "sub_result"
}

func init() {
	orm.RegisterModel(new(SubResult))
}

func (subResult *SubResult) Insert() error {
	o := orm.NewOrm()
	subResult.Sendtime, _ = util.GetNowTimeFormat()
	_, err := o.Insert(subResult)
	if err != nil {
		logs.Error("新插入订阅结果回调失败，ERROR: ", err.Error(), &subResult)
	}
	return err
}

func (subResult *SubResult) GetSubResultByUserID(userID string) (*SubResult, error) {
	o := orm.NewOrm()
	err := o.QueryTable(SubResultTBName()).Filter("user_id", userID).OrderBy("-id").One(subResult)
	if err != nil {
		logs.Error("GetSubResultByUserID 通过UserID查询订阅结果失败，ERROR: ", err.Error(), userID)
	}
	return subResult, err
}

func (subResult *SubResult) GetSubResultByTrackID(trackID string) (*SubResult, error) {
	o := orm.NewOrm()
	err := o.QueryTable(SubResultTBName()).Filter("track_id", trackID).OrderBy("-id").One(subResult)
	if err != nil {
		logs.Error("GetSubResultByTrackID 通过trackID查询订阅结果失败，ERROR: ", err.Error(), trackID)
	}
	return subResult, err
}
