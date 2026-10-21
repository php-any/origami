package data

import (
	"os"
	"strings"

	"github.com/php-any/origami/internal/pharfs"
)

type PharFileHost interface{ PharFS() *pharfs.FS }

func ReadPHPFile(vm VM, filename string) ([]byte, error) {
	if strings.HasPrefix(filename, "phar://") {
		if host, ok := vm.(PharFileHost); ok {
			return host.PharFS().ReadFile(filename)
		}
		return nil, os.ErrNotExist
	}
	return os.ReadFile(filename)
}
func StatPHPFile(vm VM, filename string) (os.FileInfo, error) {
	if strings.HasPrefix(filename, "phar://") {
		if host, ok := vm.(PharFileHost); ok {
			return host.PharFS().Stat(filename)
		}
		return nil, os.ErrNotExist
	}
	return os.Stat(filename)
}
