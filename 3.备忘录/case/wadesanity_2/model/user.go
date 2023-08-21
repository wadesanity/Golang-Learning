package model

import "time"

type User struct {
	ID       uint   `gorm:"column:id;autoIncrement;primaryKey"`
	UserName string `gorm:"column:userName;unique"`
	UserPwd  string `gorm:"column:userPwd;not null"`

	CreatedAt time.Time `gorm:"column:createdTime;"` // Set to current time if it is zero on creating
	UpdatedAt time.Time `gorm:"column:updatedTime"`  // Set to current unix seconds on updating or if it is zero on creating
}
