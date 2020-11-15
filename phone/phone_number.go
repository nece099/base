package phone

import "strings"

func NormalizePhoneNumber(phone string) string {
	if !strings.Contains(phone, "+") {
		return "+" + phone
	}

	return phone
}
