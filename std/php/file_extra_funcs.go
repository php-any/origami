package php

import (
	"os"
	"time"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsExecutableFunction 实现 is_executable 函数
// is_executable(string $filename): bool
type IsExecutableFunction struct{}

func NewIsExecutableFunction() data.FuncStmt { return &IsExecutableFunction{} }

func (f *IsExecutableFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(v.AsString())
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(info.Mode()&0111 != 0), nil
}

func (f *IsExecutableFunction) GetName() string { return "is_executable" }
var isExecutableFunctionGetParams = []data.GetValue{node.NewParameter(nil, "filename", 0, nil, nil)}

func (f *IsExecutableFunction) GetParams() []data.GetValue {
	return isExecutableFunctionGetParams
}
var isExecutableFunctionGetVariables = []data.Variable{node.NewVariable(nil, "filename", 0, nil)}

func (f *IsExecutableFunction) GetVariables() []data.Variable {
	return isExecutableFunctionGetVariables
}

// FilectimeFunction 实现 filectime 函数
// filectime(string $filename): int|false
type FilectimeFunction struct{}

func NewFilectimeFunction() data.FuncStmt { return &FilectimeFunction{} }

func (f *FilectimeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(v.AsString())
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(int(info.ModTime().Unix())), nil
}

func (f *FilectimeFunction) GetName() string { return "filectime" }
var filectimeFunctionGetParams = []data.GetValue{node.NewParameter(nil, "filename", 0, nil, nil)}

func (f *FilectimeFunction) GetParams() []data.GetValue {
	return filectimeFunctionGetParams
}
var filectimeFunctionGetVariables = []data.Variable{node.NewVariable(nil, "filename", 0, nil)}

func (f *FilectimeFunction) GetVariables() []data.Variable {
	return filectimeFunctionGetVariables
}

// FileatimeFunction 实现 fileatime 函数
// fileatime(string $filename): int|false
type FileatimeFunction struct{}

func NewFileatimeFunction() data.FuncStmt { return &FileatimeFunction{} }

func (f *FileatimeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(v.AsString())
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(int(info.ModTime().Unix())), nil
}

func (f *FileatimeFunction) GetName() string { return "fileatime" }
var fileatimeFunctionGetParams = []data.GetValue{node.NewParameter(nil, "filename", 0, nil, nil)}

func (f *FileatimeFunction) GetParams() []data.GetValue {
	return fileatimeFunctionGetParams
}
var fileatimeFunctionGetVariables = []data.Variable{node.NewVariable(nil, "filename", 0, nil)}

func (f *FileatimeFunction) GetVariables() []data.Variable {
	return fileatimeFunctionGetVariables
}

// FilepermsFunction 实现 fileperms 函数
// fileperms(string $filename): int|false
type FilepermsFunction struct{}

func NewFilepermsFunction() data.FuncStmt { return &FilepermsFunction{} }

func (f *FilepermsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(v.AsString())
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewIntValue(int(info.Mode().Perm())), nil
}

func (f *FilepermsFunction) GetName() string { return "fileperms" }
var filepermsFunctionGetParams = []data.GetValue{node.NewParameter(nil, "filename", 0, nil, nil)}

func (f *FilepermsFunction) GetParams() []data.GetValue {
	return filepermsFunctionGetParams
}
var filepermsFunctionGetVariables = []data.Variable{node.NewVariable(nil, "filename", 0, nil)}

func (f *FilepermsFunction) GetVariables() []data.Variable {
	return filepermsFunctionGetVariables
}

// FileownerFunction 实现 fileowner 函数
// fileowner(string $filename): int|false
type FileownerFunction struct{}

func NewFileownerFunction() data.FuncStmt { return &FileownerFunction{} }

func (f *FileownerFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(v.AsString())
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	if uid, ok := fileOwnerID(info); ok {
		return data.NewIntValue(uid), nil
	}
	return data.NewBoolValue(false), nil
}

func (f *FileownerFunction) GetName() string { return "fileowner" }
var fileownerFunctionGetParams = []data.GetValue{node.NewParameter(nil, "filename", 0, nil, nil)}

func (f *FileownerFunction) GetParams() []data.GetValue {
	return fileownerFunctionGetParams
}
var fileownerFunctionGetVariables = []data.Variable{node.NewVariable(nil, "filename", 0, nil)}

func (f *FileownerFunction) GetVariables() []data.Variable {
	return fileownerFunctionGetVariables
}

// FilegroupFunction 实现 filegroup 函数
// filegroup(string $filename): int|false
type FilegroupFunction struct{}

func NewFilegroupFunction() data.FuncStmt { return &FilegroupFunction{} }

func (f *FilegroupFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	info, err := os.Stat(v.AsString())
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	if gid, ok := fileGroupID(info); ok {
		return data.NewIntValue(gid), nil
	}
	return data.NewBoolValue(false), nil
}

func (f *FilegroupFunction) GetName() string { return "filegroup" }
var filegroupFunctionGetParams = []data.GetValue{node.NewParameter(nil, "filename", 0, nil, nil)}

func (f *FilegroupFunction) GetParams() []data.GetValue {
	return filegroupFunctionGetParams
}
var filegroupFunctionGetVariables = []data.Variable{node.NewVariable(nil, "filename", 0, nil)}

func (f *FilegroupFunction) GetVariables() []data.Variable {
	return filegroupFunctionGetVariables
}

// TouchFunction 实现 touch 函数
// touch(string $filename, ?int $mtime = null, ?int $atime = null): bool
type TouchFunction struct{}

func NewTouchFunction() data.FuncStmt { return &TouchFunction{} }

func (f *TouchFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	mtimeV, _ := ctx.GetIndexValue(1)
	atimeV, _ := ctx.GetIndexValue(2)

	if v == nil {
		return data.NewBoolValue(false), nil
	}
	filename := v.AsString()

	var mtime, atime time.Time
	if mtimeV != nil {
		if _, isNull := mtimeV.(*data.NullValue); !isNull {
			if iv, ok := mtimeV.(*data.IntValue); ok {
				mtime = time.Unix(int64(iv.Value), 0)
			}
		}
	}
	if mtime.IsZero() {
		mtime = time.Now()
	}
	if atimeV != nil {
		if _, isNull := atimeV.(*data.NullValue); !isNull {
			if iv, ok := atimeV.(*data.IntValue); ok {
				atime = time.Unix(int64(iv.Value), 0)
			}
		}
	}
	if atime.IsZero() {
		atime = mtime
	}

	err := os.Chtimes(filename, atime, mtime)
	if err != nil {
		// 文件不存在，创建它
		file, cerr := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
		if cerr != nil {
			return data.NewBoolValue(false), nil
		}
		file.Close()
		err = os.Chtimes(filename, atime, mtime)
		if err != nil {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func (f *TouchFunction) GetName() string { return "touch" }
var touchFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "filename", 0, nil, nil),
	node.NewParameter(nil, "mtime", 1, node.NewNullLiteral(nil), nil),
	node.NewParameter(nil, "atime", 2, node.NewNullLiteral(nil), nil),
}

func (f *TouchFunction) GetParams() []data.GetValue {
	return touchFunctionGetParams
}
var touchFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "filename", 0, nil),
	node.NewVariable(nil, "mtime", 1, nil),
	node.NewVariable(nil, "atime", 2, nil),
}

func (f *TouchFunction) GetVariables() []data.Variable {
	return touchFunctionGetVariables
}

// SysGetTempDirFunction 实现 sys_get_temp_dir 函数
// sys_get_temp_dir(): string
type SysGetTempDirFunction struct{}

func NewSysGetTempDirFunction() data.FuncStmt { return &SysGetTempDirFunction{} }

func (f *SysGetTempDirFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(os.TempDir()), nil
}

func (f *SysGetTempDirFunction) GetName() string { return "sys_get_temp_dir" }
var sysGetTempDirFunctionGetParams = []data.GetValue{}

func (f *SysGetTempDirFunction) GetParams() []data.GetValue {
	return sysGetTempDirFunctionGetParams
}
var sysGetTempDirFunctionGetVariables = []data.Variable{}

func (f *SysGetTempDirFunction) GetVariables() []data.Variable {
	return sysGetTempDirFunctionGetVariables
}

// UsleepFunction 实现 usleep 函数
// usleep(int $microseconds): void
type UsleepFunction struct{}

func NewUsleepFunction() data.FuncStmt { return &UsleepFunction{} }

func (f *UsleepFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return nil, nil
	}
	if iv, ok := v.(*data.IntValue); ok {
		time.Sleep(time.Duration(iv.Value) * time.Microsecond)
	}
	return nil, nil
}

func (f *UsleepFunction) GetName() string { return "usleep" }
var usleepFunctionGetParams = []data.GetValue{node.NewParameter(nil, "microseconds", 0, nil, data.NewBaseType("int"))}

func (f *UsleepFunction) GetParams() []data.GetValue {
	return usleepFunctionGetParams
}
var usleepFunctionGetVariables = []data.Variable{node.NewVariable(nil, "microseconds", 0, data.NewBaseType("int"))}

func (f *UsleepFunction) GetVariables() []data.Variable {
	return usleepFunctionGetVariables
}
