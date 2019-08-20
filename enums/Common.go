package enums

import "github.com/astaxie/beego"

type JsonResultCode int

type ErrorCode int

const (
	JRCodeSucc JsonResultCode = iota
	JRCodeFailed
	JRCode302 = 302 //跳转至地址
	JRCode401 = 401 //未授权访问
)

const (
	RedirectGoogle ErrorCode = iota
	Error502
)

const (
	Deleted = iota - 1
	Disabled
	Enabled
)

var DayLimitSub, _ = beego.AppConfig.Int("limitSubNum")
