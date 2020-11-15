package utils

import (
	"strings"
)

var messageFilter = []string{}

const MAX_MESSAGE_PRINT_LEN = 4096

func TruncateMessage(msg string) string {

	for _, f := range messageFilter {
		if strings.Contains(msg, f) {
			return "filtered..."
		}
	}

	if len(msg) > MAX_MESSAGE_PRINT_LEN {
		return string(msg[:MAX_MESSAGE_PRINT_LEN] + " [truncated]... ")
	}

	return FormatJson(msg)
}
