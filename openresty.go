package main

import (
	"math"
	"strconv"
	"strings"
)

func openrestyName(version string) string {
	numbers := strings.Split(version, ".")
	size := len(numbers)
	sum := 0
	for i := 0; i < size; i++ {
		n, err := strconv.Atoi(numbers[i])
		if err != nil {
			return "ngx_openresty"
		}
		sum += int(math.Pow10(size-i-1)) * n
	}

	// the source distribution name of openresty is renamed from 1.9.7.3
	if sum > 1972 {
		return "openresty"
	}

	return "ngx_openresty"
}
