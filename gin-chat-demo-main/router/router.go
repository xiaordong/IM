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
	user := r.Group("/user")
	{
		user.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})
		user.POST("/register", api.UserRegister)
		user.GET("/login", api.UserLogin)
		user.DELETE("/delete", pkg.AuthMiddleware(), api.UserDelete)
		user.PUT("/update", pkg.AuthMiddleware(), api.UserUpdate)
		user.POST("/emailCheck", pkg.AuthMiddleware(), api.EmailCheck)
		user.GET("/find", pkg.AuthMiddleware(), api.FindUser)
		user.POST("/upload", service.Upload)
	}
	friend := r.Group("/friend")
	{
		friend.POST("/add", pkg.AuthMiddleware(), api.AddFriend)
		friend.PUT("/update", pkg.AuthMiddleware(), api.UpdateFriend)
		friend.POST("/delete", pkg.AuthMiddleware(), api.DeleteFriend)
		friend.GET("/list", pkg.AuthMiddleware(), api.ListFriend)
	}
	group := r.Group("/group")
	{
		group.POST("/create", pkg.AuthMiddleware(), api.CreateGroup)
		group.POST("/join", pkg.AuthMiddleware(), api.JoinGroup)
		group.GET("/mine", pkg.AuthMiddleware(), api.MyGroup)
		group.GET("/joined", pkg.AuthMiddleware(), api.GroupJoined)
		group.DELETE("/delete", pkg.AuthMiddleware(), api.DeleteGroup)
	}
	chat := r.Group("/chat")
	{
		chat.GET("/ws", service.WsHandler)
		chat.GET("/group_ws", service.GroupWsHandler)
	}
	return r
}
