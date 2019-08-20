package util

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/utils"
)

// BeegoEmail 发邮件方法
func BeegoEmail(serviceName, resonse, deail string, emails []string) {
	fmt.Println(serviceName, " start")
	config := `{"username":"604327242@qq.com","password":"awcnfdvicdeabbbe","host":"smtp.qq.com","port":587}`
	email := utils.NewEMail(config)
	email.To = []string{"18328504774@139.com", "wangangui@mobilecpx.com"}
	if len(emails) != 0 {
		email.To = emails
	}
	email.From = "604327242@qq.com"
	email.Subject = serviceName + "  " + resonse
	email.Text = resonse
	email.HTML = deail
	err := email.Send()
	if err != nil {
		logs.Error("发送邮件失败 error: ", err.Error())
		fmt.Println(err)
		return
	}
	fmt.Println("ok")
}
