package api

import (
	"errors"
	"gateway/pkg/e"
	"gateway/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/wadesanity/hystrix-go/hystrix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

//var re = regexp.MustCompile(`hystrix: (.*)'\.`)

func ConvertGrpcError2http(err error, c *gin.Context) {
	var code int
	var errString string
	if errors.Is(err, hystrix.ErrTimeout) {
		code = http.StatusGatewayTimeout
		errString = e.GrpcResTimeoutError.Error()
	} else if errors.Is(err, hystrix.ErrMaxConcurrency) {
		code = http.StatusTooManyRequests
		errString = e.GrpcReqToManyError.Error()
	} else if errors.Is(err, hystrix.ErrCircuitOpen) {
		code = http.StatusInternalServerError
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
			case codes.DeadlineExceeded:
				code = http.StatusGatewayTimeout
				errString = e.GrpcResTimeoutError.Error()
			case codes.Canceled:
				code = http.StatusBadRequest
				errString = e.GrpcReqCancelError.Error()
			case codes.Unavailable:
				code = http.StatusInternalServerError
				errString = e.GrpcDialError.Error()
			default:
				code = http.StatusInternalServerError
				errString = e.UnknowError.Error()
			}
		} else {
			code = http.StatusInternalServerError
			errString = e.UnknowError.Error()
		}
	}

	//} else {
	//messString := err.Error()
	//if rs := re.FindStringSubmatch(messString); rs != nil {
	//	switch rs[1] {
	//	case "timeout":
	//		code = http.StatusGatewayTimeout
	//		errString = e.GrpcResTimeoutError.Error()
	//	case "circuit open":
	//		code = http.StatusInternalServerError
	//		errString = e.GrpcCircuitOpenError.Error()
	//	case "max concurrency":
	//		code = http.StatusTooManyRequests
	//		errString = e.GrpcReqToManyError.Error()
	//	default:
	//		code = http.StatusInternalServerError
	//		errString = e.UnknowError.Error()
	//	}
	//} else {
	//	code = http.StatusInternalServerError
	//	errString = e.UnknowError.Error()
	//}
	//}

	c.JSON(code, gin.H{
		"status": code,
		"error":  errString,
	})
	return
}
