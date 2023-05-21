package utils

import (
	"regexp"

	"github.com/sandromai/go-http-server/types"
)

func FormatLineBreaks(
	content string,
) (string, *types.AppError) {
	lineBreaksRegExp, err := regexp.Compile(`\r\n|\r`)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to format line breaks.",
		}
	}

	return lineBreaksRegExp.ReplaceAllString(content, "\n"), nil
}
