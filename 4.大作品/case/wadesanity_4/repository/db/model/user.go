package model

import (
	"crypto/md5"
	"fmt"
	"time"
	"videoGo/case/wadesanity_4/pkg/util"
)

type User struct {
	ID        uint   `gorm:"column:id;autoIncrement;primaryKey" json:"id"`
	Name      string `gorm:"column:name;unique" form:"name" json:"name"`
	PwdDigest string `gorm:"column:pwd_digest;not null" form:"pwd_digest" json:"pwd_digest"`
	Role      uint   `gorm:"role;default:0"`   //0普通用户 1审核员
	Status    uint   `gorm:"status;default:0"` //0正常状态 1拉黑 2封号

	Avatar   string `gorm:"column:avatar" form:"avatar" json:"avatar"`       //头像地址
	Up       string `gorm:"column:up" form:"up" json:"up"`                   //收藏id-json字符串
	Bookmark string `gorm:"column:bookmark" form:"bookmark" json:"bookmark"` //收藏id-json字符串

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"` // Set to current time if it is zero on creating
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"` // Set to current unix seconds on updating or if it is zero on creating
}

func (u *User) Md5sumPwd(userPwd string) {
	util.Logger.Debugln("userPwd before:", userPwd)
	u.PwdDigest = fmt.Sprintf("%x", md5.Sum([]byte(userPwd)))
	util.Logger.Debugln("userPwd after", u.PwdDigest)
}

func (u User) CheckPwd(userPwd string) bool {
	return u.PwdDigest == fmt.Sprintf("%x", md5.Sum([]byte(userPwd)))
}
