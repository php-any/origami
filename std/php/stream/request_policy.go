package stream

import (
	"fmt"
	"github.com/php-any/origami/data"
	"os"
)

func init() {
	data.RegisterNativeRequestPolicy[*StreamInfo](data.NativeNew, func(_ *data.RequestObjectScope, source *StreamInfo) *StreamInfo {
		file := source.openFile()
		if file == os.Stdin {
			panic(data.NewErrorThrowByName(nil, fmt.Errorf("Retained stdin cannot be shared across PHP requests; bind request input explicitly"), "Error"))
		}
		if file != os.Stdout && file != os.Stderr {
			panic(data.NewErrorThrowByName(nil, fmt.Errorf("Retained stream must be reopened for each request"), "Error"))
		}
		return NewStreamInfo(file, source.Mode)
	})
}
