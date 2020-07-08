package util

func StringP(value string) *string {
	return &value
}

// PString returns a string value from a pointer
func PString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func BoolP(b bool) *bool {
	return &b
}
