package api

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

/*
{
  "status": 200,
  "data": {
    "item": [
      {
        "id": 1,
        "title": "更改好了！",
        "content": "好耶！",
        "view": 0,
        "status": 1,
        "created_at": 1638257438,
        "start_time": 1638257437,
        "end_time": 0
      }
    ],
    "total": 1
  },
  "msg": "ok",
  "error": ""
}
*/

type Data struct {
	Item  []map[string]interface{} `json:"item"`
	Total int                      `json:"total"`
}
