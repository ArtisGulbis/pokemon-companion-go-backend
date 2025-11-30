package utils

import (
	"strings"
)

func GetId(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-2]
}
