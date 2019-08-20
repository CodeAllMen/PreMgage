package sp

import (
	"github.com/MobileCPX/PreBaseLib/splib/click"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(Mo), new(Notification), new(AffTrack), new(Postback),new(click.HourClick))
}

func MoTBName() string {
	return "mo"
}

func NotificationTBName() string {
	return "notification"
}

func PostbackTBName() string {
	return "postback"
}

func WapResponseTBName()string{
	return "wap_response"
}
