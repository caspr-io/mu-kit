package watermill

import (
	"reflect"
	"strings"
	"unicode"
)

func DefaultTopicName(m interface{}) string {
	t := reflect.TypeOf(m)

	name := ""
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	var sb strings.Builder

	for i, c := range name {
		if unicode.IsUpper(c) {
			if i != 0 {
				sb.WriteRune('.')
			}

			sb.WriteRune(unicode.ToLower(c))
		} else {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}
