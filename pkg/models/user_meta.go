package models

import (
	"gorm.io/gorm"
)

type UserMeta struct {
	*gorm.Model
	UserID uint `gorm:"user_id" json:"user_id"`
	Key    uint `gorm:"key" json:"key"`
	Value  uint `gorm:"value" json:"value"`
}
