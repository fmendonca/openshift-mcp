package utils

func GetStringArg(args map[string]any, key, def string) string {
	if v, ok := args[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

func GetIntArg(args map[string]any, key string, def int) int {
	if v, ok := args[key]; ok && v != nil {
		switch t := v.(type) {
		case int:
			return t
		case int32:
			return int(t)
		case int64:
			return int(t)
		case float64:
			return int(t)
		}
	}
	return def
}

func GetBoolArg(args map[string]any, key string, def bool) bool {
	if v, ok := args[key]; ok && v != nil {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return def
}

func InterfaceSliceToStringSlice(val any) []string {
	s, ok := val.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(s))
	for _, v := range s {
		if str, ok := v.(string); ok {
			out = append(out, str)
		}
	}
	return out
}
