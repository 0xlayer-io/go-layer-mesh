package utils

import (
	"strconv"
	"strings"
)

func StringOr(v string, d string) string {
	if v == "" {
		return d
	}
	return v
}

func StringToInt(v string, d int) int {
	r, e := strconv.Atoi(v)
	if e != nil {
		return d
	}
	return r
}

func StringToStrs(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
