package test

import (
	_ "github.com/MobileCPX/PreMgage/routers"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}


// TestBeego is a sample to run an endpoint test
func TestBeego(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestBeego", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	                So(w.Code, ShouldEqual, 200)
	        })
	        Convey("The Result Should Not Be Empty", func() {
	                So(w.Body.Len(), ShouldBeGreaterThan, 0)
	        })
	})
}


//func TestUnsub(t *testing.T){
//
//	resp := new(unsubResp)
//	msisdn := ""
//
//	mo := new(sp.Mo)
//	err := mo.GetMoByMsisdnAndService(msisdn, serviceID)
//
//
//	unsubURL := "https://api.migpay.com/subscriptions/" + mo.SubscriptionID + "/unsubscribe.json"
//	fmt.Println(unsubURL)
//	req := httplib.Post(unsubURL)
//	req.Header("Host", "api.migpay.com")
//	req.Header("MigPay-API-Key", "ec10abf99d5ef8e97c6731d74225c17ab0f94fd7")
//	req.Header("Accept", "application/json")
//	body, err := req.Bytes()
//	if err != nil {
//		c.Ctx.WriteString("Unsub Failed")
//		c.StopRun()
//	}
//	err = json.Unmarshal(body, resp)
//	if err != nil {
//		c.Ctx.WriteString("Unsub Failed")
//		c.StopRun()
//	}
//	logs.Info("SubIDUnsub", string(body))
//	if resp.Response.MigSid != "" {
//		mo := new(sp.Mo)
//		err = mo.GetMoByUserID(resp.Response.MigSid)
//		_, err = mo.UnsubUpdateMo()
//		if err == nil {
//			//c.Ctx.WriteString("Unsub SUCCESS")
//			c.redirect(serviceConfig.UnsubSuccessURL)
//		} else {
//			//c.Ctx.WriteString("Unsub ERROR")
//			c.redirect(serviceConfig.UnsubFailedURL)
//		}
//	} else {
//		c.redirect(serviceConfig.UnsubFailedURL)
//	}
//
//}