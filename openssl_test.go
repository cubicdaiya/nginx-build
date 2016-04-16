package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenSSLParallelBuildAvailable(t *testing.T) {
	assert := assert.New(t)

	// 0.9.x
	assert.Equal(false, opensslParallelBuildAvailable("0.9.6"))
	assert.Equal(false, opensslParallelBuildAvailable("0.9.7a"))
	assert.Equal(false, opensslParallelBuildAvailable("0.9.8b"))

	// 1.0.0
	assert.Equal(false, opensslParallelBuildAvailable("1.0.0"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.0a"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.0b"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.0c"))

	// 1.0.1
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1a"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1b"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1c"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1d"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1e"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1f"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1g"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1h"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1i"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1j"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1k"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1l"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1m"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1n"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.1o"))
	assert.Equal(true, opensslParallelBuildAvailable("1.0.1p"))
	assert.Equal(true, opensslParallelBuildAvailable("1.0.1za"))

	// 1.0.2
	assert.Equal(false, opensslParallelBuildAvailable("1.0.2"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.2a"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.2b"))
	assert.Equal(false, opensslParallelBuildAvailable("1.0.2c"))
	assert.Equal(true, opensslParallelBuildAvailable("1.0.2d"))
	assert.Equal(true, opensslParallelBuildAvailable("1.0.2za"))
}
