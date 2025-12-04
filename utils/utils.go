package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func GetId(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-2]
}

func ExtractIDFromURL(url string) (int, error) {
	// Remove trailing slash and split by "/"
	url = strings.TrimSuffix(url, "/")
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid URL format: %s", url)
	}
	// Last part should be the ID
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse ID from URL %s: %w", url, err)
	}
	return id, nil
}
