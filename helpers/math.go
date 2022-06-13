package helpers

func Int64(data interface{}) (int64, bool) {
	if i, ok := data.(int64); ok {
		return i, true
	}
	if i, ok := data.(int32); ok {
		return int64(i), true
	}
	if i, ok := data.(int16); ok {
		return int64(i), true
	}
	if i, ok := data.(int8); ok {
		return int64(i), true
	}
	if i, ok := data.(int); ok {
		return int64(i), true
	}
	return 0, false
}
