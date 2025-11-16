package utils

func GetStringArg(args map[string]interface{}, key, defaultValue string) string {
	if val, ok := args[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func GetIntArg(args map[string]interface{}, key string, defaultValue int) int {
	if val, ok := args[key]; ok && val != nil {
		if num, ok := val.(float64); ok {
			return int(num)
		}
		if num, ok := val.(int); ok {
			return num
		}
	}
	return defaultValue
}

func GetBoolArg(args map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := args[key]; ok && val != nil {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}

func InterfaceSliceToStringSlice(slice interface{}) []string {
	if s, ok := slice.([]interface{}); ok {
		result := make([]string, len(s))
		for i, v := range s {
			if str, ok := v.(string); ok {
				result[i] = str
			}
		}
		return result
	}
	return []string{}
}
