package redisUtil

import (
	"time"

	"github.com/astaxie/beego/logs"
)

// 同一个ip五分钟限制访问十次
const (
	ipLimtNum  = 10 //  ip访问限制个数
	expireTime = 300 * time.Second
)

// RedisCheckIPNum 检查五分钟之内ip访问次数
func RedisCheckIPNum(ip, affName, pubID string) (status bool) {
	nowTime, last5Mintime := getFormatTime()
	if redisClient != nil {
		ipNum, err := redisClient.LLen(ip).Result()
		if err != nil {
			// 出错就检查是否是redis是否挂了
			_, pingErr := redisClient.Ping().Result()
			if pingErr != nil {
				go restartRedis() // 重启redis
			}
			logs.Error("redis 获取一段时间内ip访问次数失败，ip:", ip, " error: ", err.Error())
			status = true // 为了防止redis挂了后服务器崩溃，返回true让程序继续运行
			return
		}
		if ipNum > ipLimtNum {
			ipTimeList, err := redisClient.LRange(ip, 0, ipNum).Result()
			if err != nil {
				status = true // 为了防止redis挂了后服务器崩溃，返回true让程序继续运行
				logs.Error("redis 获取一段时间内ip 列表失败，ip:", ip, " error: ", err.Error())
				return
			}
			for _, oneTime := range ipTimeList {
				if oneTime < last5Mintime {
					redisClient.LRem(ip, 1, oneTime)
				}
			}
			newipNum, _ := redisClient.LLen(ip).Result()
			if newipNum < ipLimtNum {
				status = true
			} else {
				logs.Info("五分钟内用户ip出现次数超过了10次，ip:", ip, " affName: ", affName, "pubId: ", pubID)
			}
		} else {
			status = true
		}
		redisClient.LPush(ip, nowTime)
		redisClient.Expire(ip, expireTime)
	} else {
		logs.Error("redis服务器已经挂了")
		CreateClient()
		status = true
	}

	return
}

// nowDatetim 现在的时间    last5Datetime  过去五分钟的时间
func getFormatTime() (nowDatetime, last5Mintime string) {
	time.LoadLocation("UTC")
	m, _ := time.ParseDuration("1m")
	nowDatetime = time.Now().UTC().Format("2006-01-02 15:04:05")
	last5Mintime = time.Now().Add(-5 * m).UTC().Format("2006-01-02 15:04:05")
	return
}
