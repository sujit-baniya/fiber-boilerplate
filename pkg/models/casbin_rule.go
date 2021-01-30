package models

import (
	"gorm.io/gorm"
	"time"
)

type RoleAndPermission struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	PType     string         `gorm:"p_type" json:"p_type"`
	V0        string         `gorm:"v0" json:"v0"`
	V1        string         `gorm:"v1" json:"v1"`
	V2        string         `gorm:"v2" json:"v2"`
	V3        string         `gorm:"v3" json:"v3"`
	V4        string         `gorm:"v4" json:"v4"`
	V5        string         `gorm:"v5" json:"v5"`
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
