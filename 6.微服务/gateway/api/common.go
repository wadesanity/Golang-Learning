package api

import (
	"gateway/pkg/e"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func ConvertCodeAndMessage2http(s *status.Status) (code int, errString string) {
	switch s.Code() {
	case codes.InvalidArgument:
		code = http.StatusBadRequest
		errString = e.ReqParamsError.Error()
	case codes.Unauthenticated:
		code = http.StatusUnauthorized
		errString = e.AuthorizeError.Error()
	case codes.NotFound:
		code = http.StatusNotFound
		errString = e.DbQueryNotFound.Error()
	case codes.AlreadyExists:
		code = http.StatusBadRequest
		errString = e.DbQueryAlreadyFound.Error()
	default:
		code = http.StatusInternalServerError
		errString = e.UnknowError.Error()
	}
	return
}
