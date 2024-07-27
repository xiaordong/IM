package service

import (
	"chat/model"
	"chat/pkg/e"
	"chat/serializer"
)

func Create(ownerID uint, groupName, desc string) serializer.Response {
	var group model.Group
	code := e.SUCCESS
	if err := group.Create(ownerID, groupName, desc); err != nil {
		code = e.ErrorCreateData
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   group,
	}
}

func Join(groupID uint, userID uint) serializer.Response {
	var group model.Group
	code := e.SUCCESS
	err := group.Info(groupID)
	if err != nil {
		code = e.ErrorNoData
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if group.Check(userID) {
		code = e.ExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if err = group.Join(userID); err != nil {
		code = e.ErrorCreateData
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}

	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func Delete(groupID uint, userID uint) serializer.Response {
	var gm model.GroupMember
	code := e.SUCCESS
	gm.UserID = userID
	err := gm.Info(groupID)
	if err != nil {
		code = e.ErrorNoData
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err = gm.Delete(); err != nil {
		code = e.DeleteError
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	var group model.Group
	group.ID = groupID
	group.Info(groupID)

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}

func GroupJoined(userID uint) serializer.Response {
	code := e.SUCCESS
	data, err := model.GroupJoined(userID)
	if err != nil {
		code = e.ErrorNoData
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   data,
	}
}
