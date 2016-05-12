package main

import (
	"strings"
)

// StringFlag implements the methods for flag.Value.
//
// nginx-build allows multiple flags in the specified options
// such as `--add-module` and `--add-dynamic-module`.
type StringFlag []string

func (s *StringFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *StringFlag) Set(v string) error {
	*s = append(*s, v)
	return nil
}
