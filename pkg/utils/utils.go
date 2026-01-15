package utils

import "strings"

// Contains 检查字符串切片是否包含某个元素
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

// ContainsInt 检查int切片是否包含某个元素
func ContainsInt(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

// ContainsUInt 检查uint切片是否包含某个元素
func ContainsUInt(slice []uint, item uint) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}
