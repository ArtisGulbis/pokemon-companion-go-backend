package utils

import (
	"strings"
)

type Utils struct{}

func (u *Utils) GetId(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-2]
}
