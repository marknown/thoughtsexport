package utils

import "os"

// EnableProxy 代理到数据分析工具
func EnableProxy(proxy string) {
	os.Setenv("HTTP_PROXY", proxy)
}

// DisableProxy 清除代理
func DisableProxy() {
	os.Unsetenv("HTTP_PROXY")
}
