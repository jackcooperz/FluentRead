// Package utils common 通用工具
package utils

import (
	"net/url"
	"os"
	"unicode"

	"FluentRead/misc/log"
)

// CheckStringInSlice 判断字符串是否在字符串数组中
func CheckStringInSlice(str string, strs []string) bool {
	for _, v := range strs {
		if v == str {
			return true
		}
	}
	return false
}

// GetHostByString 根据 string 获取链接的host
func GetHostByString(link string) string {
	parsedUrl, err := url.Parse(link)
	if err != nil {
		log.Errorf("链接解析失败: %v", err)
		return ""
	}
	return parsedUrl.Host
}

// IsNonChinese 检查字符串是非中文
func IsNonChinese(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			return false
		}
	}
	return true
}

// GetEnvDefault 获取环境变量，如果为空则返回备用值
func GetEnvDefault(env string, backup string) string {
	value := os.Getenv(env)
	if value == "" {
		return backup
	}
	return value
}
