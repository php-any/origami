package spl

import (
	"github.com/php-any/origami/data"
)

// splRegisteredClasses 列出 Origami 已实现的 SPL 类（�?load.go 注册保持一致）
var splRegisteredClasses = []string{
	"AppendIterator",
	"ArrayIterator",
	"ArrayObject",
	"CallbackFilterIterator",
	"CachingIterator",
	"DirectoryIterator",
	"EmptyIterator",
	"FilterIterator",
	"FilesystemIterator",
	"GlobIterator",
	"InfiniteIterator",
	"IteratorIterator",
	"LimitIterator",
	"MultipleIterator",
	"NoRewindIterator",
	"ParentIterator",
	"RecursiveCachingIterator",
	"RecursiveCallbackFilterIterator",
	"RecursiveDirectoryIterator",
	"RecursiveFilterIterator",
	"RecursiveIteratorIterator",
	"RecursiveArrayIterator",
	"RecursiveRegexIterator",
	"RecursiveTreeIterator",
	"RegexIterator",
	"SplDoublyLinkedList",
	"SplFileInfo",
	"SplFileObject",
	"SplFixedArray",
	"SplHeap",
	"SplMaxHeap",
	"SplMinHeap",
	"SplObjectStorage",
	"SplPriorityQueue",
	"SplQueue",
	"SplStack",
	"SplTempFileObject",
}

// SplClassesFunction 实现 spl_classes
type SplClassesFunction struct{}

func NewSplClassesFunction() data.FuncStmt {
	return &SplClassesFunction{}
}

func (f *SplClassesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	list := make([]*data.ZVal, len(splRegisteredClasses))
	for i, name := range splRegisteredClasses {
		list[i] = data.NewNamedZVal(name, data.NewStringValue(name))
	}
	return &data.ArrayValue{List: list}, nil
}

func (f *SplClassesFunction) GetName() string { return "spl_classes" }

func (f *SplClassesFunction) GetParams() []data.GetValue { return nil }

func (f *SplClassesFunction) GetVariables() []data.Variable { return nil }
