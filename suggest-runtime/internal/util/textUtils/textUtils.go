package textUtils

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	safeRe    = regexp.MustCompile(`(?i)[^a-zа-яё\d .,\-]+`)
	spaceRe   = regexp.MustCompile(`(?i)[ \t]+`)
	replaceYo = regexp.MustCompile(`ё`)
)

func replaceNumberCommas(s string) string {
	sb := strings.Builder{}
	sb.Grow(len(s))

	runes := []rune(s)
	for i, r := range runes {
		if r == '.' || r == ',' {
			if i == 0 || i >= len(runes)-1 {
				continue
			}
			if unicode.IsDigit(runes[i-1]) && unicode.IsDigit(runes[i+1]) {
				sb.WriteRune('.')
			}
			continue
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

func Sanitize(s string) string {
	charsToDrop := ";:\\\t\n\v\f\r\""
	s = safeRe.ReplaceAllLiteralString(s, "")
	s = replaceNumberCommas(s)
	s = spaceRe.ReplaceAllLiteralString(s, " ")
	s = replaceYo.ReplaceAllLiteralString(s, "е")
	s = strings.Trim(s, charsToDrop)
	return s
}
