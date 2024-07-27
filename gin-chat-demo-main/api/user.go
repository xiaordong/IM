package api

import (
	"chat/cache"
	"chat/conf"
	"chat/model"
	"chat/pkg"
	"chat/pkg/e"
	"chat/serializer"
	"chat/service"
	"errors"
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
	"time"
)

func UserRegister(c *gin.Context) {
	var userRegisterService service.UserRegisterService //相当于创建了一个UserRegisterService对象，调用这个对象中的Register方法。
	if err := c.ShouldBind(&userRegisterService); err == nil {
		res := userRegisterService.Register()
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		logging.Info(err)
	}
}
func UserLogin(c *gin.Context) {
	var UserLoginService service.UserRegisterService
	if err := c.ShouldBind(&UserLoginService); err == nil {
		res := UserLoginService.Login()
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		logging.Info(err)
	}
}
func UserDelete(c *gin.Context) {
	var UserLoginService service.UserRegisterService
	if err := c.ShouldBind(&UserLoginService); err == nil {
		res := UserLoginService.Delete()
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		logging.Info(err)
	}
}
func UserUpdate(c *gin.Context) {
	var u model.User
	var data string
	u.ID, _ = pkg.ParseSet(c)
	u.Email = c.PostForm("email")
	redisKey := "verification_code" + u.Email
	code := c.PostForm("verification_code")
	stiredCode, err := cache.RedisClient.Get(redisKey).Result()
	if err != nil {
		c.JSON(400, ErrorResponse(errors.New("验证码过期或不存在")))
		logging.Info(err)
		return
	}
	if stiredCode != code {
		c.JSON(400, ErrorResponse(errors.New("验证码错误")))
		logging.Info(err)
		return
	}
	u.Phone = c.PostForm("phone")
	srcFile, head, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(400, ErrorResponse(err))
		logging.Info(err)
	}
	data, err = pkg.SaveImg(conf.PersonalPath, ".png", srcFile, head)
	if err != nil {
		c.JSON(400, ErrorResponse(err))
		logging.Info(err)
	}
	u.Avatar = data
	if err = model.DB.Model(&model.User{}).Updates(u).Error; err != nil {
		c.JSON(400, ErrorResponse(err))
		logging.Info(err)
	}
	c.JSON(200, serializer.Response{
		Status: e.SUCCESS,
		Msg:    e.GetMsg(e.SUCCESS),
		Data:   data,
	})
}
func EmailCheck(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(400, ErrorResponse(errors.New("无效的邮箱")))
		logging.Info("无效邮箱")
		return
	} else {
		code, err := pkg.SendCheckCode(email)
		if err != nil {
			c.JSON(400, ErrorResponse(err))
			logging.Info(err)
			return
		}
		//插入redis
		redisKey := "verification_code" + email
		expiration := 120 * time.Second
		err = cache.RedisClient.Set(redisKey, code, expiration).Err()
		if err != nil {
			c.JSON(500, ErrorResponse(err))
			logging.Error("Failed too store verification code in Redis", err)
			return
		}
		c.JSON(200, serializer.Response{
			Status: e.SUCCESS,
			Msg:    e.GetMsg(e.SUCCESS),
		})
	}

}
func FindUser(c *gin.Context) {
	info := c.Query("info")
	list, err := service.FindUserByInfo(info)
	if err != nil {
		c.JSON(400, ErrorResponse(err))
		logging.Info(err)
		return
	}
	c.JSON(200, serializer.Response{
		Status: e.SUCCESS,
		Msg:    e.GetMsg(e.SUCCESS),
		Data:   list,
	})
}
