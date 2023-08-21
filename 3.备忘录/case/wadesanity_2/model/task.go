package model

import "time"

type Task struct {
	ID      uint   `gorm:"column:id;autoIncrement;primaryKey" form:"id"` //备忘录id
	UserID  uint   `gorm:"column:userID" form:"userID"`                      //所属用户id
	Title   string `gorm:"column:title" form:"title"`                       //备忘录题目
	Content string `gorm:"column:content" form:"content"`                     //备忘录内容
	View    uint `gorm:"column:view;default:0"`                        //备忘录查看次数
	Status  *uint   `gorm:"column:status;default:0" form:"status"`                      //备忘录状态 0未完成 1已完成

	CreatedAt time.Time `gorm:"column:createdTime"` // Set to current time if it is zero on creating
	UpdatedAt time.Time `gorm:"column:updatedTime"`  // Set to current unix seconds on updating or if it is zero on creating

}
