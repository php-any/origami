package node

// ResetSuperglobals 清空所有超全局变量缓存，确保每次 HTTP 请求都能获取最新数据
func ResetSuperglobals() {
	getValue = nil
	postValue = nil
	serverValue = nil
	requestValue = nil
	cookieValue = nil
	sessionValue = nil
	filesValue = nil
	envValue = nil
	globalsValue = nil
}
