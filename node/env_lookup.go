package node

import (
	"github.com/php-any/origami/data"
	"os"
	"sort"
)

type environmentHost interface {
	LookupPHPEnvironment(string) (string, bool)
	SetPHPEnvironment(string, *string) bool
	PHPEnvironment() map[string]string
}

// getenv/putenv operate independently of the PHP arrays $_ENV and $_SERVER.
func LookupEnvVar(ctx data.Context, name string) (string, bool) {
	if host, ok := ctx.GetVM().(environmentHost); ok {
		return host.LookupPHPEnvironment(name)
	}
	return os.LookupEnv(name)
}
func SetEnvVar(ctx data.Context, name string, value *string) bool {
	if host, ok := ctx.GetVM().(environmentHost); ok {
		return host.SetPHPEnvironment(name, value)
	}
	return false
}
func EnvironmentEntries(ctx data.Context) []string {
	if host, ok := ctx.GetVM().(environmentHost); ok {
		values := host.PHPEnvironment()
		entries := make([]string, 0, len(values))
		for name, value := range values {
			entries = append(entries, name+"="+value)
		}
		sort.Strings(entries)
		return entries
	}
	return os.Environ()
}
