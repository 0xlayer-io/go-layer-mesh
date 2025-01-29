package utils

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
)

func PrintJson(v ...any) {
	s, err := prettyjson.Marshal(v)
	if err == nil {
		fmt.Println(string(s))
	}
}
