package api_test

import (
	"testing"

	"github.com/kentik/libkflow/api"
	"github.com/stretchr/testify/assert"
)

func TestErrorStatusCodeCheck(t *testing.T) {
	assert := assert.New(t)

	e401 := &api.Error{StatusCode: 401}
	e404 := &api.Error{StatusCode: 404}

	assert.False(api.IsErrorWithStatusCode(nil, 404))
	assert.False(api.IsErrorWithStatusCode(e401, 404))
	assert.True(api.IsErrorWithStatusCode(e404, 404))
}
