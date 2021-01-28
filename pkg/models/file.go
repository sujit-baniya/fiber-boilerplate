package models

import (
	"time"

	"github.com/sujit-baniya/verify-rest/app"
	"gorm.io/gorm"
)

type File struct {
	*gorm.Model
	ModifiedAt time.Time
	FileName   string `json:"file_name" gorm:"file_name"` //nolint:gofmt
	Title      string `json:"title" gorm:"title"`         //nolint:gofmt
	MimeType   string `json:"mime_type" gorm:"mime_type"` //nolint:gofmt
	Size       string `json:"size" gorm:"size"`           //nolint:gofmt
	Extension  string `json:"extension" gorm:"extension"` //nolint:gofmt
	RowCount   int64  `json:"row_count" gorm:"row_count"` //nolint:gofmt
}

type UserFile struct {
	FileID   uint `json:"file_id" gorm:"file_id"` //nolint:gofmt
	UserID   uint `json:"user_id" gorm:"user_id"` //nolint:gofmt
	IsActive bool `json:"is_active" gorm:"is_active"`
}

func (UserFile) TableName() string {
	return "user_files"
}

func GetFileByName(name string) (*File, error) {
	var file File
	if err := app.Http.Database.DB.Where(&File{FileName: name}).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil //nolint:wsl
}
