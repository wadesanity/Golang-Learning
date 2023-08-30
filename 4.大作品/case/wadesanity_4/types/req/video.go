package types

import (
	"time"
)

type VideoCreateReq struct {
	UserID     uint   //作者id
	Title      string `form:"title" binding:"required" json:"title"` //题目
	StaticPath string
}

type VideoListReq struct {
	UserID uint //作者id
	Offset *int `form:"offset" binding:"required" json:"offset"`
	Limit  int  `json:"limit" form:"limit" binding:"required"`

	CreatedAt time.Time `gorm:"column:createdTime"` // Set to current time if it is zero on creating
	UpdatedAt time.Time `gorm:"column:updatedTime"` // Set to current unix seconds on updating or if it is zero on creating
}

type VideoActionReq struct {
	UserID uint   `json:"userID"`
	ID     uint   `json:"id" form:"id" binding:"required"`
	Value  int64  `json:"value" form:"value" binding:"required,oneof=-1 1"`
	Action string `json:"action" form:"action" binding:"required,oneof=up bookmark forward"`
}

type VideoCommentReq struct {
	UserID  uint   `json:"userID" form:"userID"`                      //用户id
	VideoID uint   `json:"videoID" form:"videoID" binding:"required"` //视频id
	Content string `json:"content" form:"content" binding:"required"` //内容
	PID     *uint  `json:"pid" form:"pid" binding:"required"`         //父id
}

type VideoCommentListReq struct {
	VideoID uint  `json:"videoID" form:"videoID" binding:"required"`
	Offset  *int  `json:"offset" form:"offset" binding:"required"`
	Limit   int   `json:"limit" form:"limit" binding:"required"`
	PID     *uint `json:"pid" form:"pid" binding:"required"`
}

type VideoTimeCommentCreateReq struct {
	UserID  uint   `json:"userID" form:"userID"`                      //用户id
	VideoID uint   `json:"videoID" form:"videoID" binding:"required"` //视频id
	Content string `json:"content" form:"content" binding:"required"` //内容
	VideoAt uint   `json:"videoAt" form:"videoAt" binding:"required"` //弹幕时间
}

type VideoTimeCommentList struct {
	VideoID    uint  `json:"videoID" form:"videoID" binding:"required"` //视频id
	VideoStart int64 `json:"videoStart" form:"videoStart" binding:"required"`
	VideoEnd   int64 `json:"videoEnd" form:"videoEnd" binding:"required"`
}
