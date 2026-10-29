//go:build !windows

package client

import (
	"os"
	"strings"
)

// consoleIsUnicode reads the locale charset from LC_ALL, LC_CTYPE then LANG, and assumes UTF-8 when none names one.
func consoleIsUnicode() bool {
	for _, name := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		value := os.Getenv(name)
		if value == "" {
			continue
		}
		charset := value
		if at := strings.IndexByte(charset, '@'); at >= 0 {
			charset = charset[:at]
		}
		dot := strings.IndexByte(charset, '.')
		if dot < 0 {
			if isCLocale(charset) {
				return false
			}
			continue
		}
		charset = strings.ToLower(strings.NewReplacer("-", "", "_", "").Replace(charset[dot+1:]))
		return charset == "utf8"
	}
	return true
}

func isCLocale(value string) bool {
	return value == "C" || value == "POSIX"
}
