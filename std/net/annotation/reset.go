package annotation

// ResetHTTPApplicationState 重置 HTTP 应用注解扫描状态，供开发模式热重载使用。
func ResetHTTPApplicationState() {
	for k := range scanningDirs {
		delete(scanningDirs, k)
	}
	for k := range registeredExitClasses {
		delete(registeredExitClasses, k)
	}
}
