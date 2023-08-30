package types

import (
	"time"
	"videoGo/case/wadesanity_4/repository/db/model"
)

type UserRegisterRes struct {
	Status int        `json:"status"`
	Msg    string     `json:"msg"`
	Data   model.User `json:"data"`
}

type UserLoginRes struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Token  string `json:"token"`
}

type UserChangePwdRes struct {
	Status   int       `json:"status"`
	Msg      string    `json:"msg"`
	UpdateAt time.Time `json:"update_at"`
}

type Data struct {
	Item  []map[string]interface{} `json:"item"`
	Total int                      `json:"total"`
}
