package openssl

import "regexp"

var (
	openSSLVersionRe    *regexp.Regexp
	openSSLVersion09xRe *regexp.Regexp
	openSSLVersion100Re *regexp.Regexp
	openSSLVersion101Re *regexp.Regexp
	openSSLVersion102Re *regexp.Regexp
)

func init() {
	openSSLVersionRe = regexp.MustCompile("(\\d+\\.\\d+\\.\\d+)([a-z]+)?")
	openSSLVersion09xRe = regexp.MustCompile("^0\\.9\\.")
	openSSLVersion100Re = regexp.MustCompile("^1\\.0\\.0")
	openSSLVersion101Re = regexp.MustCompile("^1\\.0\\.1")
	openSSLVersion102Re = regexp.MustCompile("^1\\.0\\.2")
}

func ParallelBuildAvailable(version string) bool {
	group := openSSLVersionRe.FindSubmatch([]byte(version))
	vn := group[1]
	sym := group[2]

	if openSSLVersion09xRe.Match(vn) {
		return false
	}

	if openSSLVersion100Re.Match(vn) {
		return false
	}

	if openSSLVersion101Re.Match(vn) {
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

	if openSSLVersion102Re.Match(vn) {
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
