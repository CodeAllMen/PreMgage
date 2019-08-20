package sp

import (
	"encoding/json"
	"fmt"
	"github.com/MobileCPX/PreBaseLib/common"
	"github.com/MobileCPX/PreMgage/models/sp"
	"io/ioutil"
)

type SetPostbackController struct {
	common.BaseController
}

func (c *SetPostbackController) Get() {
	affName := c.GetString("aff_name")
	promoter := c.GetString("promoter")
	postbackURL := c.GetString("postback_url")
	payout, _ := c.GetFloat("payout")
	postbackRate, _ := c.GetInt("rate")
	offerID, _ := c.GetInt("offer_id")
	campID, _ := c.GetInt("camp_id")
	fmt.Println(affName, promoter, postbackURL, offerID,campID)
	postback := new(sp.Postback)
	if offerID != 0 && campID != 0 {
		err := postback.CheckOfferID(offerID)
		if err == nil && postback.ID != 0 {
			c.Ctx.WriteString("ERROR,OfferID已经存在")
			c.StopRun()
		} else {
			postback.AffName = affName
			postback.PromoterName = promoter
			postback.OfferID = offerID
			postback.CampID = campID
			postback.PostbackRate = postbackRate
			postback.PostbackURL = postbackURL
			postback.Payout = float32(payout)
			if postbackRate == 0 {
				postbackRate = 50
			}
			err = postback.InsertPostback()
			if err != nil {
				c.Ctx.WriteString("ERROR,插入postbak失败")
				c.StopRun()
			}
		}
	} else {
		c.Ctx.WriteString("ERROR,offerID 是空")
		c.StopRun()
	}
	c.Ctx.WriteString("SUCCESS")
}




func (c *SetPostbackController) Post() {
	postback := new(sp.Postback)
	reqBody := c.Ctx.Request.Body
	reqByte, err := ioutil.ReadAll(reqBody)
	if err == nil {
		_ = json.Unmarshal(reqByte, postback)
		fmt.Println(postback)
	} else {
		c.StringResult("ERROR,json解析失败： " + err.Error())
	}

	if postback.OfferID != 0 && postback.CampID != 0 {
		oldPostback := new(sp.Postback)
		err = oldPostback.CheckOfferIDIsExist(postback.OfferID)
		// 如果offerId已经存在，则只需要更新
		if err == nil && oldPostback.ID != 0 {
			postback.ID = oldPostback.ID
			postback.CreateTime = oldPostback.CreateTime

			err = postback.Update()
			if err == nil {
				c.StringResult("SUCCESS")
			} else {
				c.StringResult("ERROR, 更新postback失败")
			}
		} else {
			// 新插入postback
			err = postback.Insert()
			if err == nil {
				c.StringResult("SUCCESS")
			} else {
				c.StringResult("ERROR, 存入postback失败")
			}
		}
	} else {
		c.Ctx.WriteString("ERROR,offerID 是空")
		c.StopRun()
	}
}
