package model

import (
	"time"
)

type Video struct {
	ID            uint   `gorm:"column:id;autoIncrement;primaryKey" form:"id" json:"id"` //id
	UserID        uint   `gorm:"column:userID" form:"userID" json:"user_id"`                  //作者id
	Title         string `gorm:"column:title" form:"title" json:"title"`                    //题目
	StaticPath    string `gorm:"column:static_path" form:"static_path" json:"static_path"`        //内容
	Status        uint  `gorm:"column:status;default:0" form:"status" json:"status"`        //状态: 0未审核 1已审核
	ViewCount     uint   `gorm:"column:viewCount;default:0" json:"view_count"`                   //浏览次数
	BookmarkCount uint   `gorm:"column:bookmarkCount;default:0" json:"bookmark_count"`               //浏览次数
	UpCount       uint   `gorm:"column:upCount;default:0" json:"up_count"`
	ForwardCount  uint   `gorm:"column:forwardCount;default:0" json:"forward_count"`
	TimeCommentCount uint `json:"time_comment_count" gorm:"column:timeCommentCount;default:0"`

	CreatedAt time.Time `gorm:"column:createdTime" json:"created_at"` // Set to current time if it is zero on creating
	UpdatedAt time.Time `gorm:"column:updatedTime" json:"updated_at"` // Set to current unix seconds on updating or if it is zero on creating

}
