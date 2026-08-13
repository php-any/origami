package os

import (
	"github.com/php-any/origami/data"
	"os"
	"path/filepath"
	"runtime"
)

func newOs() *OS {
	eol := "\n"
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}

	return &OS{
		EOL: eol,
	}
}

type OS struct {
	EOL string
}

// Exit 退出程序
func (o *OS) Exit(code int) {
	os.Exit(code)
}

// Hostname 获取主机名
func (o *OS) Hostname() (string, error) {
	return os.Hostname()
}

func (o *OS) Path(paths data.ArrayValue) string {
	var rets []string
	valueList := paths.ToValueList()
	for _, value := range valueList {
		rets = append(rets, value.AsString())
	}
	return filepath.Join(rets...)
}
