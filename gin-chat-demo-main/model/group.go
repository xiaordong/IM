package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Group struct {
	gorm.Model
	GroupName  string `json:"group_name"`
	GroupOwner uint   `json:"group_owner"`
	Desc       string `json:"desc" gorm:"size:100"`
	// 为 Size 设置默认值，如果数据库中没有这个值，将使用默认值 200
	Size int `json:"size" gorm:"default:200"`
}
type GroupMember struct {
	gorm.Model
	GroupID uint `json:"group_id"`
	UserID  uint `json:"user_id"`
	Role    int  `json:"role"`
}

// Create 创建群以及群组成员模型
func (g *Group) Create(groupOwner uint, groupName, desc string) error {
	g.GroupOwner = groupOwner
	g.GroupName = groupName
	g.Desc = desc
	if err := DB.Model(&Group{}).Create(&g).Error; err != nil {
		return err
	}
	member := GroupMember{
		GroupID: g.ID,
		UserID:  groupOwner,
		Role:    1,
	}
	return DB.Model(&Group{}).Create(&member).Error
}

// Info 拿到一个群信息
func (g *Group) Info(groupID uint) error {
	if err := DB.Model(&Group{}).Where("id = ?", groupID).First(&g).Error; err != nil {
		return err
	}
	return nil
}

// Check 检查是否以及在群里
func (g *Group) Check(userID uint) bool {
	var member GroupMember
	if DB.Model(&GroupMember{}).Where("group_id = ? and user_id = ?", g.ID, userID).First(&member).RowsAffected != 0 {
		fmt.Println("true", member)
		return true //已经在群里
	}
	return false //不在群里
}

// Join 加入群
func (g *Group) Join(UserID uint) error {
	member := GroupMember{
		GroupID: g.ID,
		UserID:  UserID,
		Role:    0,
	}
	return DB.Model(&GroupMember{}).Create(&member).Error
}

// Info 群成员的信息
func (gm *GroupMember) Info(GroupId uint) error {
	return DB.Model(&GroupMember{}).Where("group_id = ? and user_id = ?", GroupId, gm.UserID).First(&gm).Error
}

// Delete 删除群成员
func (gm *GroupMember) Delete() error {
	// 检查是否是群主
	if gm.Role == 1 {
		// 查找下一个最高角色的成员，如果没有其他成员，则解散群组
		var newOwner GroupMember
		if err := DB.Model(&GroupMember{}).Where("group_id = ? AND user_id != ?", gm.GroupID, gm.UserID).
			Order("role DESC, created_at ASC").Limit(1).Find(&newOwner).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 如果没有其他成员，解散群组
				return DisbandGroup(gm.GroupID)
			}
			return err
		}

		// 将新群主的 Role 更新为 1
		if err := DB.Model(&GroupMember{}).Where("id = ?", newOwner.ID).Update("role", 1).Error; err != nil {
			return err
		}

		// 更新群组表的 GroupOwner 为新群主的 UserID
		var group Group
		if err := DB.Where("id = ?", gm.GroupID).First(&group).Error; err != nil {
			return err
		}
		group.GroupOwner = newOwner.UserID
		if err := DB.Save(&group).Error; err != nil {
			return err
		}
	}

	// 删除群成员记录
	return DB.Delete(gm).Error
}

// DisbandGroup 解散群组
func DisbandGroup(groupID uint) error {
	// 开启事务
	tx := DB.Begin()

	// 先尝试删除群成员记录
	if err := tx.Where("group_id = ?", groupID).Delete(GroupMember{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 然后删除群组记录
	if err := tx.Where("id = ?", groupID).Delete(Group{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// MyGroup 用户创建的群列表
func MyGroup(ownerId uint) ([]Group, error) {
	var groups []Group
	if err := DB.Model(&Group{}).Where("group_owner = ?", ownerId).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// GroupJoined 加入的群
func GroupJoined(userID uint) ([]Group, error) {
	var groups []GroupMember
	if err := DB.Model(&GroupMember{}).Where("user_id = ?", userID).Find(&groups).Error; err != nil { //一组groupID
		return nil, err
	}
	idList := make([]uint, len(groups))
	for i, g := range groups {
		idList[i] = g.GroupID
	}
	var data []Group
	if err := DB.Model(&Group{}).Find(&data, idList).Error; err != nil {
		return nil, err
	}
	return data, nil
}
