package utils

// SliceContainStr 判断字符串 slice 里是否含有指定 value
func SliceContainStr(src []string, value string) bool {
	isContain := false
	for _, srcValue := range src {
		if srcValue == value {
			isContain = true
			break
		}
	}
	return isContain
}

// SliceContainInt 判断int数组 slice 里是否含有指定 value
func SliceContainInt(src []int, value int) bool {
	isContain := false
	for _, srcValue := range src {
		if srcValue == value {
			isContain = true
			break
		}
	}
	return isContain
}
