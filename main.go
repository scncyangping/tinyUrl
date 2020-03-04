package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"tinyUrl/common/constants"
	"tinyUrl/config"
	"tinyUrl/config/db/mongo"
	"tinyUrl/config/db/redis"
	"tinyUrl/config/global"
	"tinyUrl/config/log"
	"tinyUrl/routes"
)

// 初始化
func init() {
	config.Init()
	redis.Init()
	mongo.Init()
	global.Init()
}

func main() {
	//	for gin log file
	mode := os.Getenv(constants.StartRunMode)

	if mode == constants.PRO { //生产环境
		f, _ := os.Create(config.Base.Log.Dir + "/http.log")
		gin.DefaultWriter = io.MultiWriter(f)
	}
	r := gin.Default()
	err := routes.InitRoute(r)
	if err != nil {
		log.GetLogger().Error("init http router error", err)
		return
	}

	_ = r.Run(":" + config.Base.Server.Port)
}
