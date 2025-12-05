package util

import (
	"strings"
)

// Base62Chars chứa 62 ký tự (0-9, a-z, A-Z)
const Base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const Base = 62

// EncodeToBase62 chuyển đổi một số nguyên (ID) sang chuỗi Base62.
func EncodeToBase62(n uint) string {
	if n == 0 {
		return string(Base62Chars[0])
	}

	var result strings.Builder
	for n > 0 {
		remainder := n % Base
		result.WriteByte(Base62Chars[remainder])
		n /= Base
	}

	// chuỗi bị ngược cần đảo lại
	return reverseString(result.String())
}

// reverseString đảo ngược chuỗi
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
