package utils

import "regexp"

func GetDeviceInfo(
	userAgent string,
) (
	platform,
	browser string,
) {
	if matched, _ := regexp.MatchString(`(?i)Android`, userAgent); matched {
		platform = "Android"
	} else if matched, _ := regexp.MatchString(`(?i)Linux`, userAgent); matched {
		platform = "Linux"
	} else if matched, _ := regexp.MatchString(`(?i)iPhone`, userAgent); matched {
		platform = "iPhone"
	} else if matched, _ := regexp.MatchString(`(?i)iPad`, userAgent); matched {
		platform = "iPad"
	} else if matched, _ := regexp.MatchString(`(?i)iPod`, userAgent); matched {
		platform = "iPod"
	} else if matched, _ := regexp.MatchString(`(?i)Macintosh|Mac OS X`, userAgent); matched {
		platform = "Mac"
	} else if matched, _ := regexp.MatchString(`(?i)Windows|Win32|Win64`, userAgent); matched {
		platform = "Windows"
	}

	if matched, _ := regexp.MatchString(`(?i)OPR`, userAgent); matched {
		browser = "Opera"
	} else if matched, _ := regexp.MatchString(`(?i)Edg`, userAgent); matched {
		browser = "Microsoft Edge"
	} else if matched, _ := regexp.MatchString(`(?i)FxiOS|Firefox`, userAgent); matched {
		browser = "Mozilla Firefox"
	} else if matched, _ := regexp.MatchString(`(?i)CriOS|Chrome`, userAgent); matched {
		browser = "Google Chrome"
	} else if matched, _ := regexp.MatchString(`(?i)Version`, userAgent); matched {
		if matched, _ := regexp.MatchString(`(?i)Safari`, userAgent); matched {
			browser = "Safari"
		}
	}

	return
}
