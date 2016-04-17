package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenrestyName(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("ngx_openresty", openrestyName("1.9.7.2"))
	assert.Equal("openresty", openrestyName("1.9.7.3"))
	assert.Equal("openresty", openrestyName("1.9.7.4"))
}
