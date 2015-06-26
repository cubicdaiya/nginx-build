package main

import (
	"strings"
)

type StringFlag []string

func (s *StringFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *StringFlag) Set(v string) error {
	*s = append(*s, v)
	return nil
}
