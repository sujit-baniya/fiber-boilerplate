package web

import (
	"errors"

	"github.com/thomasvvugt/fiber-boilerplate/app/models"
	"github.com/thomasvvugt/fiber-boilerplate/database"
)

// Return a single user as JSON
func FindUserByUsername(username string) (*models.User, error) {
	db := database.Instance()
	User := new(models.User)
	if res := db.Where("name = ?", username).First(&User); res.Error != nil {
		return User, res.Error
	}
	if User.ID == 0 {
		return User, errors.New("user not found")
	}
	// Match role to user
	if User.RoleID != 0 {
		Role := new(models.Role)
		if res := db.Find(&Role, User.RoleID); res.Error != nil {
			return User, errors.New("error when retrieving the role of the user")
		}
		if Role.ID != 0 {
			User.Role = *Role
		}
	}
	return User, nil
}

// Return a single user as JSON
func FindUserByID(id int64) (*models.User, error) {
	db := database.Instance()
	User := new(models.User)
	if res := db.Where("id = ?", id).First(&User); res.Error != nil {
		return User, res.Error
	}
	if User.ID == 0 {
		return User, errors.New("user not found")
	}
	// Match role to user
	if User.RoleID != 0 {
		Role := new(models.Role)
		if res := db.Find(&Role, User.RoleID); res.Error != nil {
			return User, errors.New("error when retrieving the role of the user")
		}
		if Role.ID != 0 {
			User.Role = *Role
		}
	}
	return User, nil
}
