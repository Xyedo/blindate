package mod

import (
	"strings"
)

func TrimWhiteSpace(s *string) {
	*s = strings.Join(strings.Fields(*s), " ")
}

func Trim(s *string) {
	*s = strings.TrimSpace(*s)
}
