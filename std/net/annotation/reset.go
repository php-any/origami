package annotation

// ResetHTTPApplicationState 重置 HTTP 应用注解扫描状态，供开发模式热重载使用。
func ResetHTTPApplicationState() {
	applicationState.Lock()
	defer applicationState.Unlock()
	for k := range applicationState.scanningDirs {
		delete(applicationState.scanningDirs, k)
	}
	for k := range applicationState.registeredExitClasses {
		delete(applicationState.registeredExitClasses, k)
	}
}
