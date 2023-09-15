package service

import (
	"errors"
	"gateway/pkg/e"
	"gateway/pkg/util"
	"github.com/wadesanity/hystrix-go/hystrix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

//var re = regexp.MustCompile(`hystrix: (.*)'\.`)

func ConvertGrpcError2http(err error) error {
	if err == nil {
		return nil
	}
	var httpStatus int = http.StatusInternalServerError
	var errString string = e.UnknowError.Error()
	if errors.Is(err, hystrix.ErrTimeout) {
		httpStatus = http.StatusGatewayTimeout
		errString = e.GrpcResTimeoutError.Error()
	} else if errors.Is(err, hystrix.ErrMaxConcurrency) {
		httpStatus = http.StatusTooManyRequests
		errString = e.GrpcReqToManyError.Error()
	} else if errors.Is(err, hystrix.ErrCircuitOpen) {
		httpStatus = http.StatusInternalServerError
		errString = e.GrpcCircuitOpenError.Error()
	} else {
		util.Logger.Infof("not hystrix err:%v", err)
		err1 := errors.Unwrap(err)
		if err1 != nil {
			err = err1
		}
		util.Logger.Infof("after unwrap err:%v", err)
		s, ok := status.FromError(err)
		if ok {
			switch s.Code() {
			case codes.InvalidArgument:
				httpStatus = http.StatusBadRequest
				errString = e.ReqParamsError.Error()
			case codes.Unauthenticated:
				httpStatus = http.StatusUnauthorized
				errString = e.AuthorizeError.Error()
			case codes.NotFound:
				httpStatus = http.StatusNotFound
				errString = e.DbQueryNotFound.Error()
			case codes.AlreadyExists:
				httpStatus = http.StatusBadRequest
				errString = e.DbQueryAlreadyFound.Error()
			case codes.DeadlineExceeded:
				httpStatus = http.StatusGatewayTimeout
				errString = e.GrpcResTimeoutError.Error()
			case codes.Canceled:
				httpStatus = http.StatusBadRequest
				errString = e.GrpcReqCancelError.Error()
			case codes.Unavailable:
				httpStatus = http.StatusInternalServerError
				errString = e.GrpcDialError.Error()
			default:
				httpStatus = http.StatusInternalServerError
				errString = e.UnknowError.Error()
			}
		}
	}
	return e.NewApiError(httpStatus, errString)
}
