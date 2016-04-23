package openssl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenSSLParallelBuildAvailable(t *testing.T) {
	assert := assert.New(t)

	// 0.9.x
	assert.Equal(false, ParallelBuildAvailable("0.9.6"))
	assert.Equal(false, ParallelBuildAvailable("0.9.7a"))
	assert.Equal(false, ParallelBuildAvailable("0.9.8b"))

	// 1.0.0
	assert.Equal(false, ParallelBuildAvailable("1.0.0"))
	assert.Equal(false, ParallelBuildAvailable("1.0.0a"))
	assert.Equal(false, ParallelBuildAvailable("1.0.0b"))
	assert.Equal(false, ParallelBuildAvailable("1.0.0c"))

	// 1.0.1
	assert.Equal(false, ParallelBuildAvailable("1.0.1"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1a"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1b"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1c"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1d"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1e"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1f"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1g"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1h"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1i"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1j"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1k"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1l"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1m"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1n"))
	assert.Equal(false, ParallelBuildAvailable("1.0.1o"))
	assert.Equal(true, ParallelBuildAvailable("1.0.1p"))
	assert.Equal(true, ParallelBuildAvailable("1.0.1za"))

	// 1.0.2
	assert.Equal(false, ParallelBuildAvailable("1.0.2"))
	assert.Equal(false, ParallelBuildAvailable("1.0.2a"))
	assert.Equal(false, ParallelBuildAvailable("1.0.2b"))
	assert.Equal(false, ParallelBuildAvailable("1.0.2c"))
	assert.Equal(true, ParallelBuildAvailable("1.0.2d"))
	assert.Equal(true, ParallelBuildAvailable("1.0.2za"))
}
