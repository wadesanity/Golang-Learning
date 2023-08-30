package types

import (
	"time"
)

type VideoListRes struct {
	ID               uint   `gorm:"column:id;autoIncrement;primaryKey" form:"id" json:"id"`   //id
	UserID           uint   `gorm:"column:userID" form:"userID" json:"user_id"`               //作者id
	Title            string `gorm:"column:title" form:"title" json:"title"`                   //题目
	StaticPath       string `gorm:"column:static_path" form:"static_path" json:"static_path"` //内容
	ViewCount        uint   `gorm:"column:viewCount;default:0" json:"view_count"`             //浏览次数
	TimeCommentCount uint   `json:"time_comment_count" gorm:"column:timeCommentCount;default:0"`

	CreatedAt time.Time `gorm:"column:createdTime" json:"created_at"` // Set to current time if it is zero on creating
	UpdatedAt time.Time `gorm:"column:updatedTime" json:"updated_at"` // Set to current unix seconds on updating or if it is zero on creating
}

type VideoOneRes struct {
	ID            uint `json:"id"`
	ViewCount     uint `json:"viewCount"`
	UpCount       uint `json:"upCount"`
	BookmarkCount uint `json:"bookmarkCount"`
	ForwardCount  uint `json:"forwardCount"`
	IsUp          bool `json:"isUp"`
	IsBookMark    bool `json:"isBookMark"`
	IsForward     bool `json:"isForward"`
}
