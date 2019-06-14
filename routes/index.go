package routes

import (
	"github.com/gin-gonic/gin"
	"tinyUrl/common/middleware"
	"tinyUrl/routes/business"
)

func InitRoute(router *gin.Engine) error {

	v1 := router.Group("v1")
	err := InitApiRoute(v1)
	return err
}

//func InitHttpRoute(router *gin.RouterGroup) error {
//	//业务逻辑路由
//	routerAction := router.Group("/bus")
//	// 登陆
//	router.Use(middleware.VisitLog())
//	business.InitBusinessRoute(routerAction)
//
//	return nil
//}

func InitApiRoute(router *gin.RouterGroup) error {
	router.Use(middleware.VisitLog())
	//业务逻辑路由
	routerAction := router.Group("api")
	// 登陆
	business.InitBusinessRoute(routerAction)

	return nil
}
