package model

import "time"

type Comment struct {
	ID      uint   `gorm:"column:id;autoIncrement;primaryKey" form:"id"` //id
	UserID  uint   `gorm:"column:userID" form:"userID"`                  //用户id
	VideoID uint   `gorm:"column:videoID" form:"videoID"`                //视频id
	Content string `gorm:"column:content" form:"content"`                //内容
	PID     uint   `gorm:"column:pid" form:"pid"`                        //父id

	CreatedAt time.Time `gorm:"column:createdTime"` // Set to current time if it is zero on creating
	UpdatedAt time.Time `gorm:"column:updatedTime"` // Set to current unix seconds on updating or if it is zero on creating

}

type TimeComment struct {
	ID      uint      `gorm:"column:id;autoIncrement;primaryKey" form:"id"` //id
	UserID  uint      `gorm:"column:userID" form:"userID"`                  //用户id
	VideoID uint      `gorm:"column:videoID" form:"videoID"`                //视频id
	Content string    `gorm:"column:content" form:"content"`                //内容
	VideoAt time.Time `gorm:"column:videoTime"`                             //弹幕时间

	CreatedAt time.Time `gorm:"column:createdTime"` // Set to current time if it is zero on creating
	UpdatedAt time.Time `gorm:"column:updatedTime"` // Set to current unix seconds on updating or if it is zero on creating
}
