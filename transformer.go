package goconfig

import (
	"strings"
	"unicode"
)

// VoidTransformer is the default field name transformer that returns the key unchanged.
// It can be used with WithKeyTransformer to maintain the original field names.
func VoidTransformer(key string) string {
	return key
}

// UperCaseTransformer transforms PascalCase field names into MACRO_CASE environment variable names.
// It handles various cases including:
//   - DBConnection => DB_CONNECTION
//   - AKey => A_KEY
//   - KeyA => KEY_A
//   - ThisISMyKey => THIS_IS_MY_KEY
//
// This transformer is useful when you want to maintain a consistent uppercase naming convention
// for environment variables while using PascalCase in your Go code.
func UperCaseTransformer(key string) string {
	var (
		sb         strings.Builder
		runes      = []rune(key)
		n          = len(runes)
		canToUpper = func(c rune) bool {
			return unicode.ToUpper(c) != c
		}
	)

	for i, c := range runes {
		if i > 0 && unicode.IsUpper(c) && unicode.IsLetter(c) && runes[i-1] != '_' &&
			((i < n-1 && canToUpper(runes[i+1])) ||
				(canToUpper(runes[i-1]))) {
			_, _ = sb.WriteRune('_')
			_, _ = sb.WriteRune(c)
		} else {
			_, _ = sb.WriteRune(unicode.ToUpper(c))
		}
	}

	return sb.String()
}
