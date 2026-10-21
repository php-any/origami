package php

import (
	"github.com/php-any/origami/data"
	"os"
	"reflect"
)

type StatFunction struct {
	IsFileFunction
	link bool
}

func (f *StatFunction) GetName() string {
	if f.link {
		return "lstat"
	}
	return "stat"
}
func (f *StatFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	var info os.FileInfo
	var err error
	if f.link {
		info, err = os.Lstat(v.AsString())
	} else {
		info, err = data.StatPHPFile(ctx.GetVM(), v.AsString())
	}
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	mode := int(info.Mode().Perm())
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		mode |= 0120000
	case info.IsDir():
		mode |= 0040000
	case info.Mode()&os.ModeNamedPipe != 0:
		mode |= 0010000
	case info.Mode()&os.ModeSocket != 0:
		mode |= 0140000
	default:
		mode |= 0100000
	}
	values := []int{0, 0, mode, 1, 0, 0, 0, int(info.Size()), int(info.ModTime().Unix()), int(info.ModTime().Unix()), int(info.ModTime().Unix()), -1, -1}
	// Unix exposes these fields through Stat_t; Windows supplies FILETIMEs.
	metadata := reflect.Indirect(reflect.ValueOf(info.Sys()))
	if metadata.IsValid() && metadata.Kind() == reflect.Struct {
		for index, name := range map[int]string{0: "Dev", 1: "Ino", 2: "Mode", 3: "Nlink", 4: "Uid", 5: "Gid", 6: "Rdev", 11: "Blksize", 12: "Blocks"} {
			field := metadata.FieldByName(name)
			if field.IsValid() {
				if field.CanInt() {
					values[index] = int(field.Int())
				} else if field.CanUint() {
					values[index] = int(field.Uint())
				}
			}
		}
		for index, names := range map[int][]string{8: {"Atim", "Atimespec", "LastAccessTime"}, 10: {"Ctim", "Ctimespec", "CreationTime"}} {
			for _, name := range names {
				field := metadata.FieldByName(name)
				if !field.IsValid() || field.Kind() != reflect.Struct {
					continue
				}
				sec := field.FieldByName("Sec")
				if sec.IsValid() && sec.CanInt() {
					values[index] = int(sec.Int())
					break
				}
				low, high := field.FieldByName("LowDateTime"), field.FieldByName("HighDateTime")
				if low.IsValid() && high.IsValid() {
					ticks := high.Uint()<<32 | low.Uint()
					values[index] = int(int64(ticks/10000000) - 11644473600)
					break
				}
			}
		}
	}
	names := []string{"dev", "ino", "mode", "nlink", "uid", "gid", "rdev", "size", "atime", "mtime", "ctime", "blksize", "blocks"}
	list := make([]data.Value, len(values))
	for i, v := range values {
		list[i] = data.NewIntValue(v)
	}
	out := data.NewArrayValue(list).(*data.ArrayValue)
	for i, name := range names {
		out.SetStringKey(name, list[i])
	}
	return out, nil
}
