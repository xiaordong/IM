package service

import (
	"chat/model"
	"chat/pkg"
	"chat/pkg/e"
	"chat/serializer"
	logging "github.com/sirupsen/logrus"
	"strconv"
)

// UserRegisterService 管理用户注册服务
type UserRegisterService struct { //相当于要求这个服务需要的数据，用于json绑定
	UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=15"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=16"`
}

func (service UserRegisterService) Register() serializer.Response { //因为是注册用户，属于new一个实例，所以不需要传指针。
	var user model.User
	var count int
	code := e.SUCCESS
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).Count(&count)
	if count != 0 {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户名已存在",
		}
	}
	user = model.User{
		UserName: service.UserName,
		Status:   model.Active,
	}
	//加密密码
	if err := user.SetPassword(service.Password); err != nil {
		logging.Info(err)
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user.Avatar = "http://q1.qlogo.cn/g?b=qq&nk=294350394&s=640"
	//创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
func (service UserRegisterService) Login() serializer.Response {
	var user model.User
	code := e.SUCCESS
	if res := model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).Find(&user); res.Error != nil || res.RowsAffected == 0 || user.CheckPassword(service.Password) == false {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户名或密码错误",
		}
	}
	aToken, rToken, _ := pkg.GenToken(strconv.Itoa(int(user.ID)))
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   map[string]string{"aToken": aToken, "rToken": rToken},
	}
}
func (service UserRegisterService) Delete() serializer.Response {
	var user model.User
	code := e.SUCCESS
	if res := model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).Find(&user); res.Error != nil || res.RowsAffected == 0 || user.CheckPassword(service.Password) == false {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户名或密码错误",
		}
	}
	if err := model.DB.Unscoped().Delete(&user).Error; err != nil {
		logging.Info(err)
		code = e.DeleteError
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func MyGroup(userID uint) serializer.Response {
	code := e.SUCCESS
	data, err := model.MyGroup(userID)
	if err != nil {
		logging.Info(err)
		code = e.ErrorNoData
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   data,
	}
}
