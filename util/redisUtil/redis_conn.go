package redisUtil

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/go-redis/redis"
)

var redisRestartStatus bool // true 表示正在重启  false 表示redis正常运行

var redisClient *redis.Client

// CreateClient 创建 redis 客户端
func CreateClient() *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 5, // 十个连接池
	})

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)
	if err != nil {
		logs.Error("redis 连接失败")
	} else {
		logs.Info("redis 连接成功")
	}
	return redisClient
}

// 当检查到redis服务器挂了之后重新连接redis服务
func restartRedis() {
	if !redisRestartStatus {
		restartNum := 1
		redisRestartStatus = true
		for {
			for i := 1; i < 11; i++ {
				time.Sleep(3 * time.Second)
				restartNum++   // 重启次数
				CreateClient() // 重新创建redis连接
				_, pingErr := redisClient.Ping().Result()
				if pingErr != nil {
					logs.Error("第", restartNum, "次尝试重新连接redis 失败")
				} else {
					logs.Info("第", restartNum, "次尝试重新连接redis 成功")
					redisRestartStatus = false
					break
				}
			}
			// 表示重启成功
			if redisRestartStatus == false {
				break
			} else {
				BeegoEmail("NTH 服务器", " 连接 redis 失败", "")
			}
			time.Sleep(600 * time.Second)
		}
	}
}
