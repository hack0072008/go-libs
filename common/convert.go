package common

func GetString(origin interface{}, defaultValue ...string) string {
	if origin != nil {
		if value, ok := origin.(string); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func GetInt(origin interface{}, defaultValue ...int) int {
	if origin != nil {
		if value, ok := origin.(int); ok {
			return value
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}
