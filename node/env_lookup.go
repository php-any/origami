package node

import (
	"os"
	"strings"

	"github.com/php-any/origami/data"
)

func ensureEnvValue() *data.ObjectValue {
	if envValue == nil {
		envValue = data.NewObjectValue()
		mergeOSEnviron(envValue)
	}
	return envValue
}

func mergeOSEnviron(target *data.ObjectValue) {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			target.SetProperty(parts[0], data.NewStringValue(parts[1]))
		}
	}
}

func lookupInObject(o *data.ObjectValue, name string) (string, bool) {
	if o == nil {
		return "", false
	}
	val, ok := o.GetProperties()[name]
	if !ok || val == nil {
		return "", false
	}
	return val.AsString(), true
}

// LookupEnvVar 按 PHP/Symfony 惯例查找环境变量：先 $_ENV，再 $_SERVER（非 HTTP_），最后 OS 环境。
func LookupEnvVar(name string) (string, bool) {
	if value, ok := lookupInObject(envValue, name); ok {
		return value, true
	}
	if !strings.HasPrefix(name, "HTTP_") {
		if value, ok := lookupInObject(serverValue, name); ok {
			return value, true
		}
	}
	return os.LookupEnv(name)
}

// SetEnvVar 设置环境变量，并同步 $_ENV / $_SERVER。
func SetEnvVar(name, value string) error {
	if err := os.Setenv(name, value); err != nil {
		return err
	}
	ensureEnvValue().SetProperty(name, data.NewStringValue(value))
	if !strings.HasPrefix(name, "HTTP_") {
		if serverValue == nil {
			serverValue = data.NewObjectValue()
		}
		serverValue.SetProperty(name, data.NewStringValue(value))
	}
	return nil
}
