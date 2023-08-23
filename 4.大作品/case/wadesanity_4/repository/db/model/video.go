package model

import "time"

type Video struct {
	ID      uint   `gorm:"column:id;autoIncrement;primaryKey" form:"id"` //id
	UserID  uint   `gorm:"column:userID" form:"userID"`                  //作者id
	Title   string `gorm:"column:title" form:"title"`                    //题目
	Content string `gorm:"column:content" form:"content"`                //内容
	Status  *uint  `gorm:"column:status;default:0" form:"status"`        //状态: 0未审核 1已审核
	ViewCount    uint   `gorm:"column:viewCount;default:0"`                        //浏览次数
	BookmarkCount    uint   `gorm:"column:bookmarkCount;default:0"`                        //浏览次数
	UpCount uint	`gorm:"column:upCount;default:0"`
	ForwardCount uint	`gorm:"column:forwardCount;default:0"`

	CreatedAt time.Time `gorm:"column:createdTime"` // Set to current time if it is zero on creating
	UpdatedAt time.Time `gorm:"column:updatedTime"` // Set to current unix seconds on updating or if it is zero on creating

}
