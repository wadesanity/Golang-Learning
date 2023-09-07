package res

import (
	"gateway/pkg/util"
	"net/http"
)

type Response struct {
	Status int    `json:"status"`
	Data   any    `json:"data,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Error  string `json:"error,omitempty"`
}

type DataList struct {
	Items any   `json:"items"`
	Total int64 `json:"total"`
}

func NewResList(dataList any, total int64, msg string) *Response {
	util.Logger.Debugf("dataList is %#v, %v", dataList, dataList == nil)
	if dataList == nil {
		util.Logger.Debugf("dataList is nil")
		dataList = []int{}
	}
	return &Response{
		Status: http.StatusOK,
		Data: DataList{
			Items: dataList,
			Total: total,
		},
		Msg: msg,
	}
}

func NewResError(status int, err error) *Response {
	return &Response{
		Status: status,
		Error:  err.Error(),
	}
}

func NewResOk(msg string, status int, data any) *Response {
	return &Response{
		Status: status,
		Msg:    msg,
		Data:   data,
	}
}
