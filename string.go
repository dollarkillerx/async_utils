package async_utils

import "strings"

func Combinations(str []string, sep string) string {
	var build strings.Builder
	for k, v := range str {
		if k < len(str)-1 {
			build.WriteString(v)
			build.WriteString(sep)
		} else {
			build.WriteString(v)
		}
	}
	return build.String()
}

