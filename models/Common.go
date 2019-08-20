package models

import "github.com/MobileCPX/PreMgage/enums"

// JsonResult 用于返回请求的基类
type JsonResult struct {
	Code enums.JsonResultCode `json:"code"`
	Msg  string               `json:"msg"`
	Obj  interface{}          `json:"obj"`
}



