package openresty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenrestyName(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("ngx_openresty", Name("1.9.7.2"))
	assert.Equal("openresty", Name("1.9.7.3"))
	assert.Equal("openresty", Name("1.9.7.4"))
	assert.Equal("openresty", Name("1.15.8.1rc1"))
}
