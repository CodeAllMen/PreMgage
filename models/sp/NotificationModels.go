package sp

import (
	"github.com/MobileCPX/PreMgage/util"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// 成功订阅通知   /sp/notification?period_from=2019-03-21T11%3A36%3A57Z&s=889&t=60&period_to=2019-03-28T11%3A36%3A57Z&migid=0536e03d&sig=18d6a1413c07eaf019e2f2451116bc82e412e19b&mig_sid=2482024&mig_stage=CustomerSubscriptionCreated
// GoNotification 订阅，续订、退订通知
type Notification struct {
	ID               int64  `orm:"pk;auto;column(id)"`
	NotificationType string `orm:"column(notification_type)"`
	Sendtime         string `orm:"column(sendtime)"`
	SubscriptionID   string `orm:"column(subscription_id)"`
	TransactionID    string `orm:"column(transaction_id)"`
	SubStage         string `orm:"column(mig_stage)"` // 订阅阶段，有订阅，续订，退订等阶段，通过此来判断通知类型
	Msisdn           string `orm:"column(msisdn)"`
	UserID           string `orm:"column(user_id)"` // 用户的唯一标识(类似电话号码）
	MigSid           string
	Sign             string `orm:"column(sign)"`

	CurrentPeriodEnd   string // 开始时间
	CurrentPeriodStart string // 结束时间

	TrackID   string
	Operator  string
	ServiceID string
}

func (notification *Notification) Insert() error {
	o := orm.NewOrm()
	nowTime, _ := util.GetNowTimeFormat()
	notification.Sendtime = nowTime
	_, err := o.Insert(notification)
	if err != nil {
		logs.Error("Notification Insert 数据失败，ERROR: ", err.Error())
	}
	return err
}

// GetIdentifyNotificationByTrackID 根据trackID 获取通知信息
func (notification *Notification) GetIdentifyNotificationByTrackID(trackID string) error {
	o := orm.NewOrm()
	err := o.QueryTable("notification").Filter("request_id__istartswith", trackID+"_identify").
		OrderBy("-id").One(notification)
	if err != nil {
		logs.Error("GetIdentiryNotification ERROR", err.Error())
	}
	return err
}

func (notification *Notification) GetUnsubIdentiryNotification(trackID string) error {
	o := orm.NewOrm()
	err := o.QueryTable("notification").Filter("request_id", trackID).
		OrderBy("-id").One(notification)
	if err != nil {
		logs.Error("GetIdentiryNotification ERROR", err.Error())
	}
	return err
}
