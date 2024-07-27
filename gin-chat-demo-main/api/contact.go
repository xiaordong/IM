package api

import (
	"chat/model"
	"chat/pkg"
	"chat/pkg/e"
	"chat/serializer"
	"chat/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func AddFriend(c *gin.Context) {
	var f model.Friend
	f.OwnerID, _ = pkg.ParseSet(c)
	temp, _ := strconv.Atoi(c.PostForm("target_id"))
	f.TargetID = uint(temp)
	if f.TargetID == f.OwnerID {
		c.JSON(400, ErrorResponse(errors.New("不能添加自己为好友")))
		return
	}
	if f.Check() {
		c.JSON(400, ErrorResponse(errors.New("已添加用户")))
		return
	}
	if err := f.Add(); err != nil {
		log.Println(err)
		c.JSON(500, ErrorResponse(errors.New("添加好友失败")))
		return
	}
	c.JSON(200, serializer.Response{
		Status: e.SUCCESS,
		Msg:    e.GetMsg(e.SUCCESS),
	})
}
func UpdateFriend(c *gin.Context) {
	var f model.Friend
	f.OwnerID, _ = pkg.ParseSet(c)
	temp, _ := strconv.Atoi(c.PostForm("target_id"))
	f.TargetID = uint(temp)
	str := c.PostForm("type")
	f.Push(f.OwnerID, f.TargetID)
	f.Type, _ = strconv.Atoi(str)
	if err := f.Update(model.DB); err != nil {
		log.Println(err)
		c.JSON(400, ErrorResponse(err))
		return
	}
	c.JSON(200, serializer.Response{
		Status: e.SUCCESS,
		Msg:    e.GetMsg(e.SUCCESS),
	})
}
func DeleteFriend(c *gin.Context) {
	var f model.Friend
	ownerId, _ := pkg.ParseSet(c)
	temp, _ := strconv.Atoi(c.PostForm("target_id"))
	if err := f.Push(ownerId, uint(temp)); err != nil {
		log.Println(err)
		c.JSON(400, ErrorResponse(err))
		return
	}
	fmt.Println("arrive", f)
	if f.Delete(model.DB) != nil {
		log.Println("删除关系失败")
		c.JSON(500, ErrorResponse(errors.New("删除失败")))
		return
	}
	c.JSON(200, serializer.Response{
		Status: e.SUCCESS,
		Msg:    e.GetMsg(e.SUCCESS),
	})

}
func ListFriend(c *gin.Context) {
	ownerId, _ := pkg.ParseSet(c)
	var f model.Friend
	datta, err := f.Get(ownerId)
	if err != nil {
		log.Println(err)
		c.JSON(400, ErrorResponse(err))
		return
	}
	c.JSON(200, serializer.Response{
		Status: e.SUCCESS,
		Msg:    e.GetMsg(e.SUCCESS),
		Data:   datta,
	})

}
func CreateGroup(c *gin.Context) {
	ownerId, _ := pkg.ParseSet(c)
	groupName := c.PostForm("name")
	description := c.PostForm("description")
	if groupName == "" {
		c.JSON(400, ErrorResponse(errors.New("群名不能为空")))
		return
	}
	res := service.Create(ownerId, groupName, description)
	c.JSON(200, res)
}

func JoinGroup(c *gin.Context) {
	userId, _ := pkg.ParseSet(c)
	str := c.PostForm("group_id")
	if str == "" {
		c.JSON(400, ErrorResponse(errors.New("群ID为空")))
		log.Println("群ID为空")
		return
	}
	groupID, _ := strconv.Atoi(str)
	res := service.Join(uint(groupID), userId)
	c.JSON(200, res)
}

func DeleteGroup(c *gin.Context) {
	userId, _ := pkg.ParseSet(c)
	str := c.PostForm("group_id")
	if str == "" {
		c.JSON(400, ErrorResponse(errors.New("群ID为空")))
		return
	}
	groupID, _ := strconv.Atoi(str)
	res := service.Delete(uint(groupID), userId)
	c.JSON(200, res)
}
func MyGroup(c *gin.Context) {
	userId, _ := pkg.ParseSet(c)
	res := service.MyGroup(userId)
	c.JSON(200, res)
}
func GroupJoined(c *gin.Context) {
	userId, _ := pkg.ParseSet(c)
	res := service.GroupJoined(userId)
	c.JSON(200, res)
}
