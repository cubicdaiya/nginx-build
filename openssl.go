package main

import "regexp"

var (
	OpenSSLVersionRe    *regexp.Regexp
	OpenSSLVersion09xRe *regexp.Regexp
	OpenSSLVersion100Re *regexp.Regexp
	OpenSSLVersion101Re *regexp.Regexp
	OpenSSLVersion102Re *regexp.Regexp
)

func init() {
	OpenSSLVersionRe = regexp.MustCompile("(\\d+\\.\\d+\\.\\d+)([a-z]+)?")
	OpenSSLVersion09xRe = regexp.MustCompile("^0\\.9\\.")
	OpenSSLVersion100Re = regexp.MustCompile("^1\\.0\\.0")
	OpenSSLVersion101Re = regexp.MustCompile("^1\\.0\\.1")
	OpenSSLVersion102Re = regexp.MustCompile("^1\\.0\\.2")
}

func opensslParallelBuildAvailable(version string) bool {
	group := OpenSSLVersionRe.FindSubmatch([]byte(version))
	vn := group[1]
	sym := group[2]

	if OpenSSLVersion09xRe.Match(vn) {
		return false
	}

	if OpenSSLVersion100Re.Match(vn) {
		return false
	}

	if OpenSSLVersion101Re.Match(vn) {
		if len(sym) == 0 {
			return false
		}
		if len(sym) > 1 {
			return true
		}
		// len(sym) == 1
		symn := group[2][0]
		if symn > 111 { // ord('o') => 111
			return true
		}
		return false
	}

	if OpenSSLVersion102Re.Match(vn) {
		if len(sym) == 0 {
			return false
		}
		if len(sym) > 1 {
			return true
		}
		// len(sym) == 1
		symn := group[2][0]
		if symn > 99 { // ord('99') => c
			return true
		}

		return false
	}

	return true
}
