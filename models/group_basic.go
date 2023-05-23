package models

import (
	"gorm.io/gorm"
)

// 群
type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Type    string
	Desc    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
