package models

import (
	"gorm.io/gorm"
	"time"
)

type RoleAndPermission struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Ptype     string         `gorm:"size:100;uniqueIndex:unique_index"`
	V0        string         `gorm:"size:100;uniqueIndex:unique_index"`
	V1        string         `gorm:"size:100;uniqueIndex:unique_index"`
	V2        string         `gorm:"size:100;uniqueIndex:unique_index"`
	V3        string         `gorm:"size:100;uniqueIndex:unique_index"`
	V4        string         `gorm:"size:100;uniqueIndex:unique_index"`
	V5        string         `gorm:"size:100;uniqueIndex:unique_index"`
	Category  string         `gorm:"category" json:"category"`
}

type RoleRequest struct {
	UserID  uint   `json:"user_id"`
	Role    string `json:"role"`
	OldRole string `json:"old_role"`
}

type PermissionRequest struct {
	Role   string `json:"role"`
	Module string `json:"module"`
	Action string `json:"action"`
	Route  string `json:"route"`
	Method string `json:"method"`
}

func (RoleAndPermission) TableName() string {
	return "casbin_rule"
}
