package service

import (
	"fmt"
	"gateway/pkg/e"
	"github.com/stretchr/testify/assert"
	"github.com/wadesanity/hystrix-go/hystrix"
	"net/http"

	"testing"
)

func Test_Common(t *testing.T) {
	err := ConvertGrpcError2http(nil)
	assert.Nil(t, err)

	err = ConvertGrpcError2http(fmt.Errorf("unknow test err"))
	var apiError *e.ApiError
	assert.ErrorAs(t, err, &apiError)
	assert.Equal(t, http.StatusInternalServerError, apiError.HttpStatus)

	err = ConvertGrpcError2http(hystrix.ErrTimeout)
	assert.ErrorAs(t, err, &apiError)
	assert.Equal(t, http.StatusGatewayTimeout, apiError.HttpStatus)
}
