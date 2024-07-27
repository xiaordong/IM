package model

import (
	"github.com/jinzhu/gorm"
)

type Friend struct {
	gorm.Model
	OwnerID  uint `json:"owner_id"`
	TargetID uint `json:"target_id"`
	Type     int  `gorm:"default:1"`
}

// Push 从数据中拉取 一段关系
func (f *Friend) Push(ownerID uint, targetID uint) error {
	return DB.Model(&Friend{}).Where("owner_id = ? and target_id = ?", ownerID, targetID).First(&f).Error
}

// Add 添加好友
func (f *Friend) Add() error {
	return DB.Create(f).Error
}

// Delete 删除好友:获取OwnerID,targetID->得到一个Friend数据，执行delete方法
func (f *Friend) Delete(db *gorm.DB) error {
	return db.Delete(f).Error
}

// Update 更改好友类型
func (f *Friend) Update(db *gorm.DB) error {
	return db.Model(&Friend{}).Save(&f).Error
}

// Check 检查是否存在这段关系
func (f *Friend) Check() bool {
	if res := DB.Model(&Friend{}).Where("owner_id=? and target_id = ?", f.OwnerID, f.TargetID).First(&f); res.RowsAffected != 0 {
		return true
	}
	return false //不存在，可添加
}

// Get 得到所有好友列表
func (f *Friend) Get(ownerId uint) ([]Friend, error) {
	var friends []Friend
	if err := DB.Model(&Friend{}).Where("owner_id=?", ownerId).Find(&friends).Error; err != nil {
		return nil, err
	}
	return friends, nil
}
