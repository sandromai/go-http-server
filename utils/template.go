package utils

import "strings"

func UseTemplate(
	content string,
	variables map[string]string,
) string {
	for key, value := range variables {
		content = strings.ReplaceAll(content, "{"+key+"}", value)
	}

	return content
}
