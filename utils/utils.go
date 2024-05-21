package utils

// StringInSlice 检查一个字符串是否在slice
func StrInSlice(str string, s []string) bool {
	if len(s) == 0 {
		return false
	}
	for _, v := range s {
		if str == v {
			return true
		}
	}

	return false
}

// If 三目运算符  warning:panic 不会截断 注意value里的panic
func If(condition bool, trueValue, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}
