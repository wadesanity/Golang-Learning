package common

type ResponseError struct {
	Status int   `json:"status"`
	Error  string `json:"error"`
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

type ResponseOk struct {
	Status int    `json:"status"`
	Data   Data   `json:"data"`
	Msg    string `json:"msg"`
}

type Data struct {
	Item  []map[string]interface{} `json:"item"`
	Total int                      `json:"total"`
}
