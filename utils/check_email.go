package utils

import (
	"regexp"
	"strings"
)

func CheckEmail(email string) bool {
	if matched, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, strings.TrimSpace(email)); !matched || err != nil {
		return false
	}

	return true
}
