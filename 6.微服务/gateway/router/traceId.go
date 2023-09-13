package router

import (
	"context"
	"gateway/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func traceIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := strconv.FormatInt(time.Now().Unix(), 10)
		rc := context.WithValue(c.Request.Context(), util.TraceIdKey, traceId)
		rc1, cancel := context.WithTimeout(rc, time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(rc1)
		util.Logger.Debugf("set traceId:%v done", c.Request.Context().Value(util.TraceIdKey))
		c.Next()
		if c.Writer.Status() == http.StatusGatewayTimeout {

		}
	}
}
