package router

import (
	"chat/api"
	"chat/conf"
	"chat/pkg"
	"chat/service"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	conf.Init()
	r := gin.Default()
	r.Use(gin.Recovery(), gin.Logger())
	v1 := r.Group("/user")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})
		v1.POST("/register", api.UserRegister)
		v1.GET("/login", api.UserLogin)
		v1.DELETE("/delete", pkg.AuthMiddleware(), api.UserDelete)
		v1.PUT("/update", pkg.AuthMiddleware(), api.UserUpdate)
		v1.POST("/emailCheck", pkg.AuthMiddleware(), api.EmailCheck)
		v1.GET("/find", pkg.AuthMiddleware(), api.FindUser)
		v1.POST("/upload", service.Upload)
	}
	v2 := r.Group("/friend")
	{
		v2.POST("/add", pkg.AuthMiddleware(), api.AddFriend)
		v2.PUT("/update", pkg.AuthMiddleware(), api.UpdateFriend)
		v2.POST("/delete", pkg.AuthMiddleware(), api.DeleteFriend)
		v2.GET("/list", pkg.AuthMiddleware(), api.ListFriend)
	}
	v3 := r.Group("/group")
	{
		v3.POST("/create", pkg.AuthMiddleware(), api.CreateGroup)
		v3.POST("/join", pkg.AuthMiddleware(), api.JoinGroup)
		v3.GET("/mine", pkg.AuthMiddleware(), api.MyGroup)
		v3.GET("/joined", pkg.AuthMiddleware(), api.GroupJoined)
		v3.DELETE("/delete", pkg.AuthMiddleware(), api.DeleteGroup)
	}
	v4 := r.Group("/chat")
	{
		v4.GET("/ws", service.WsHandler)
		v4.GET("/group_ws", service.GroupWsHandler)
	}
	return r
}
