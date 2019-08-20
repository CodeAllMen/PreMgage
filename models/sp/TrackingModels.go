package sp

import (
	"errors"
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/click"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreMgage/util"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
)

// AffTrack 网盟点击追踪
type AffTrack struct {
	TrackID  int64  `orm:"pk;auto;column(track_id)"`  //自增ID
	Sendtime string `orm:"column(sendtime);size(30)"` // 点击时间
	//AffName     string `orm:"column(aff_name);size(30)"`   // 网盟名称
	//PubID       string `orm:"column(pub_id);size(100)"`    // 子渠道
	//ProID       string `orm:"column(pro_id);size(30)"`     // 服务id（可有可无）
	//ClickID     string `orm:"column(click_id);size(100)"`  // 点击
	//ServiceID   string `orm:"column(service_id);size(30)"` // 服务类型
	RequestID string `orm:"column(request_id)"`
	//ServiceName string `orm::column(service_name)"`
	//IP          string `orm:"column(ip);size(20)"` // 用户IP地址
	//UserAgent   string `orm:"column(user_agent)"`  // 用户user_agent
	//Refer       string `orm:"column(refer)"`       // 网页来源
	CanvasID string `orm:"column(canvas_id)"` // 帆布ID
	CookieID string `orm:"column(cookie_id)"` // CookieID

	//ClickTime   string `orm:"column(click_time)" json:"click_time"`
	//ServiceID   string `orm:"column(service_id)" json:"service_id"`
	//ServiceName string `orm:"column(service_name)" json:"service_name"`
	//OfferID     int    `orm:"column(offer_id)"  json:"offer_id"`          // offerID  Admin 平台对接的OfferID
	//CampID      int    `orm:"column(camp_id)" json:"camp_id"`             // CampID  Admin 平台对接的CampID
	//AffName     string `orm:"column(aff_name);size(35)" json:"aff_name"`  // 网盟名称
	//PubID       string `orm:"column(pub_id)" json:"pub_id"`               // 子渠道
	//ClickID     string `orm:"column(click_id);size(350)" json:"click_id"` // clickID
	//ProID       string `orm:"column(pro_id)" json:"pro_id"`
	//IP          string `orm:"column(ip);size(25)" json:"ip"`                  // 用户的IP地址
	//UserAgent   string `orm:"column(user_agent);size(500)" json:"user_agent"` // 用户的UserAgent
	//IpAs        string `orm:"column(ip_as);size(50)" json:"ip_as"`            // 用户IP对应的AS   例如AS16512
	//Refer       string `orm:"column(refer);size(500)" json:"refer"`           // refer 来源
	//Other       string `orm:"column(other)" json:"other"`

	tracking.Track

	//
	//OfferID   int64  `orm:"column(offer_id)"`
	//CampID    int64  `orm:"column(camp_id)"`
	//OtherData string `orm:"column(other_data)"`
}

func (track *AffTrack) TableName() string {
	return "aff_track"
}

func (track *AffTrack) Insert() (int64, error) {
	o := orm.NewOrm()
	track.Sendtime, _ = util.GetNowTimeFormat()
	track.ClickTime = track.Sendtime
	trackID, err := o.Insert(track)
	logs.Info(trackID, "1111111")
	if err != nil {
		logs.Error("新插入点击错误 ", err.Error())
	}
	fmt.Println(track)
	return trackID, err
}

func (track *AffTrack) InsertTable() (int64, error) {
	o := orm.NewOrm()
	trackID, err := o.Insert(track)
	fmt.Println(trackID)
	return trackID, err
}

func (track *AffTrack) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(track)
	if err != nil {
		logs.Error("AffTrack Update 更新点击数据失败，ERROR ", err.Error())
	}
	return err
}

func (track *AffTrack) GetAffTrackByTrackID(trackID int64) error {
	o := orm.NewOrm()
	track.TrackID = trackID
	err := o.Read(track)
	if err != nil {
		logs.Error("通过trackID 查询点击信息失败，未找到此trackID： ", trackID)
	}
	return err
}

func (track *AffTrack) GetAffTrackByRequestID(requestID string) error {
	o := orm.NewOrm()
	err := o.QueryTable("aff_track").Filter("request_id", requestID).One(track)
	if err != nil {
		logs.Error("通过RequestID 查询点击信息失败，未找到此RequestID： ", requestID)
	}
	return err
}

func GetServiceIDByTrackID(trackID string) (*AffTrack, error) {
	o := orm.NewOrm()
	track := new(AffTrack)
	trackIDInt, err := strconv.Atoi(trackID)
	if err != nil {
		logs.Error("GetServiceIDByTrackID track string to int 错误，ERROR: ", err.Error(), " trackID: ", trackID)
		return track, errors.New("track string to int error")
	}

	track.TrackID = int64(trackIDInt)
	err = o.Read(track)
	if err != nil {
		logs.Error("GetServiceIDByTrackID 通过trackID 查询aff_track 表失败，ERROR: ", err.Error(), " trackID: ", trackID)
		return track, errors.New("没有查询到数据")
	}
	return track, err
}

func InsertHourClick() {
	o := orm.NewOrm()
	hourClick := new(click.HourClick)
	nowTime, _ := util.GetNowTimeFormat()
	nowHour := nowTime[:13]
	fmt.Println(nowHour)
	hourTime := hourClick.GetNewestClickDateTime()
	if hourTime == "" {
		hourTime = "2019-07-01"
	}

	totalHourClick := new([]click.HourClick)
	//SQL := fmt.Sprintf("SELECT left(sendtime,13) as hour_time,postback_price, (case service_id when '889-Vodafone' "+
	//	"THEN 3 WHEN '889-Three' THEN 4 WHEN '892-Vodafone' THEN 11 WHEN '892-Three' THEN 12 ELSE 0 END) as"+
	//	" camp_id, offer_id,aff_name,pub_id,count(1) as click_num ,click_status, promoter_id "+
	//	"from aff_track   where service_id <> ''  and left(sendtime,13)>'%s' and left(sendtime,13)<'%s' group by "+
	//	"left(sendtime,13),offer_id,aff_name,pub_id,"+
	//	"service_id,pro_id ,promoter_id,postback_price,click_status order by left(sendtime,13)", hourTime, nowHour)

	SQL := fmt.Sprintf("SELECT left(sendtime,13) as hour_time,postback_price, "+
		" camp_id, offer_id,aff_name,pub_id,count(1) as click_num ,click_status, promoter_id "+
		"from aff_track   where service_id <> ''  and left(sendtime,13)>'%s' and left(sendtime,13)<'%s' group by "+
		"left(sendtime,13),offer_id,aff_name,pub_id,"+
		"service_id,pro_id ,promoter_id,camp_id,postback_price,click_status order by left(sendtime,13)", hourTime, nowHour)

	num, _ := o.Raw(SQL).QueryRows(totalHourClick)
	fmt.Println(num)

	for _, v := range *totalHourClick {
		if v.ClickNum >= 2 && v.CampID != 0 {
			o.Insert(&v)
		}
		fmt.Println(v.HourTime, v.PubID, v.ClickNum, v.AffName, v.OfferID, v.CampID)
	}
}
