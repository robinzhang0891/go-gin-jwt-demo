package router

import (
	"go-gin-jwt-demo/controller"
	"go-gin-jwt-demo/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// 创建gin实例，默认使用debug模式
	r := gin.Default()

	// Only trust forwarding headers from  local Nginx proxy.
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		panic(err)
	}
	r.RemoteIPHeaders = []string{"X-Real-IP"}

	// added for front end web
	r.GET("/", func(c *gin.Context) {
		c.File("./web/login.html")
	})

	r.GET("/login", func(c *gin.Context) {
		c.File("./web/login.html")
	})
	r.GET("/register", func(c *gin.Context) {
		c.File("./web/register.html")
	})

	r.GET("/profile", func(c *gin.Context) {
		c.File("./web/profile.html")
	})
	// end of added for front end web

	// 公开路由组，不需要登录鉴权
	public := r.Group("/api")
	{
		public.POST("/register", controller.Register)
		public.POST("/login", controller.Login)
	}

	// 需要JWT鉴权的私有路由组
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuth()) // 注册JWT中间件
	{
		auth.GET("/user/info", controller.Info)
	}

	return r
}
