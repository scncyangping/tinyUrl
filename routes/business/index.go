package business

import (
	"github.com/gin-gonic/gin"
	"tinyUrl/handler/business"
	"tinyUrl/routes/business/tinyUrl"
)

/*
 * date : 2019/4/30
 * author : yangping
 * desc : 所有业务模块路由汇总，每个模块具体方法再分文件写
 */
func InitBusinessRoute(router *gin.RouterGroup) {

	// 鉴权中间件
	// router.Use(middleware.TokenAuthMiddleware())
	tinyGroup := router.Group("/tiny")
	tinyUrl.InitTinyUrlRoute(tinyGroup)

	router.GET("/go", business.Redirect4TinyUrl)
	router.POST("/login", business.Login)
	router.POST("/register", business.Register)

}
