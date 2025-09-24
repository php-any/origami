package log

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewLogClass() data.ClassStmt {
	source := NewLog()
	return &LogClass{
		debug:  &LogDebugMethod{source},
		error:  &LogErrorMethod{source},
		fatal:  &LogFatalMethod{source},
		info:   &LogInfoMethod{source},
		notice: &LogNoticeMethod{source},
		trace:  &LogTraceMethod{source},
		warn:   &LogWarnMethod{source},
	}
}

type LogClass struct {
	node.Node
	debug  data.Method
	error  data.Method
	fatal  data.Method
	info   data.Method
	notice data.Method
	trace  data.Method
	warn   data.Method
}

func (s *LogClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s

	return &clone, nil
}

func (s *LogClass) GetName() string {
	return "Log"
}

func (s *LogClass) GetExtend() *string {
	return nil
}

func (s *LogClass) GetImplements() []string {
	return nil
}

func (s *LogClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (s *LogClass) GetProperties() map[string]data.Property {
	return nil
}

func (s *LogClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "debug":
		return s.debug, true
	case "error":
		return s.error, true
	case "fatal":
		return s.fatal, true
	case "info":
		return s.info, true
	case "notice":
		return s.notice, true
	case "trace":
		return s.trace, true
	case "warn":
		return s.warn, true
	}
	return nil, false
}

func (s *LogClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "debug":
		return s.debug, true
	case "error":
		return s.error, true
	case "fatal":
		return s.fatal, true
	case "info":
		return s.info, true
	case "notice":
		return s.notice, true
	case "trace":
		return s.trace, true
	case "warn":
		return s.warn, true
	}
	return nil, false
}

func (s *LogClass) GetMethods() []data.Method {
	return []data.Method{
		s.debug,
		s.error,
		s.fatal,
		s.info,
		s.notice,
		s.trace,
		s.warn,
	}
}

func (s *LogClass) GetConstruct() data.Method {
	return nil
}
