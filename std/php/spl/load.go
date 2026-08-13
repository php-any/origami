package spl

import (
	"github.com/php-any/origami/data"
)

// Load 注册 SPL 扩展：函数、接口、类与相关常量
func Load(vm data.VM) {
	for _, fun := range []data.FuncStmt{
		NewSplAutoloadRegisterFunction(),
		NewSplAutoloadUnregisterFunction(),
		NewSplAutoloadFunctionsFunction(),
		NewSplClassesFunction(),
		NewSplAutoloadCallFunction(),
		NewSplObjectHashFunction(),
		NewSplObjectIdFunction(),
	} {
		vm.AddFunc(fun)
	}

	vm.AddInterface(NewTraversableInterface())
	vm.AddInterface(NewIteratorAggregateInterface())
	vm.AddInterface(NewIteratorInterface())
	vm.AddInterface(NewRecursiveIteratorInterface())
	vm.AddInterface(NewOuterIteratorInterface())
	vm.AddInterface(NewSplObserverInterface())
	vm.AddInterface(NewSplSubjectInterface())
	vm.AddInterface(NewSeekableIteratorInterface())

	vm.AddClass(NewRecursiveDirectoryIteratorClass())
	vm.AddClass(NewRecursiveIteratorIteratorClass())
	vm.AddClass(NewSplFileInfoClass())
	vm.AddClass(NewSplFileObjectClass())
	vm.AddClass(NewSplTempFileObjectClass())
	vm.AddClass(NewGlobIteratorClass())
	vm.AddClass(&DirectoryIteratorClass{})
	vm.AddClass(NewFilesystemIteratorClass())
	vm.AddClass(NewArrayIteratorClass())
	vm.AddClass(NewRecursiveArrayIteratorClass())
	vm.AddClass(NewFilterIteratorClass())
	vm.AddClass(NewIteratorIteratorClass())
	vm.AddClass(NewRecursiveFilterIteratorClass())
	vm.AddClass(NewCallbackFilterIteratorClass())
	vm.AddClass(NewRecursiveCallbackFilterIteratorClass())
	vm.AddClass(NewRegexIteratorClass())
	vm.AddClass(NewRecursiveRegexIteratorClass())
	vm.AddClass(NewLimitIteratorClass())
	vm.AddClass(NewCachingIteratorClass())
	vm.AddClass(NewRecursiveCachingIteratorClass())
	vm.AddClass(NewNoRewindIteratorClass())
	vm.AddClass(NewInfiniteIteratorClass())
	vm.AddClass(NewAppendIteratorClass())
	vm.AddClass(NewMultipleIteratorClass())
	vm.AddClass(NewParentIteratorClass())
	vm.AddClass(NewRecursiveTreeIteratorClass())
	vm.AddClass(NewArrayObjectClass())
	vm.AddClass(NewEmptyIteratorClass())
	vm.AddClass(NewSplDoublyLinkedListClass())
	vm.AddClass(NewSplStackClass())
	vm.AddClass(NewSplQueueClass())
	vm.AddClass(NewSplFixedArrayClass())
	vm.AddClass(NewSplHeapClass())
	vm.AddClass(NewSplMinHeapClass())
	vm.AddClass(NewSplMaxHeapClass())
	vm.AddClass(NewSplPriorityQueueClass())
	vm.AddClass(NewSplObjectStorageClass())

	initConstants(vm)
}

func initConstants(vm data.VM) {
	vm.SetConstant("ArrayObject::STD_PROP_LIST", data.NewIntValue(1))
	vm.SetConstant("ArrayObject::ARRAY_AS_PROPS", data.NewIntValue(2))

	vm.SetConstant("ArrayIterator::STD_PROP_LIST", data.NewIntValue(1))
	vm.SetConstant("ArrayIterator::ARRAY_AS_PROPS", data.NewIntValue(2))

	vm.SetConstant("SplPriorityQueue::EXTR_DATA", data.NewIntValue(SpqExtrData))
	vm.SetConstant("SplPriorityQueue::EXTR_PRIORITY", data.NewIntValue(SpqExtrPriority))
	vm.SetConstant("SplPriorityQueue::EXTR_BOTH", data.NewIntValue(SpqExtrBoth))

	vm.SetConstant("RecursiveTreeIterator::PREORDER", data.NewIntValue(0))
	vm.SetConstant("RecursiveTreeIterator::POSTORDER", data.NewIntValue(1))

	vm.SetConstant("CachingIterator::CALL_TOSTRING", data.NewIntValue(1))
	vm.SetConstant("CachingIterator::CATCH_GET_CHILD", data.NewIntValue(2))
	vm.SetConstant("CachingIterator::TOSTRING_USE_KEY", data.NewIntValue(4))
	vm.SetConstant("CachingIterator::TOSTRING_USE_CURRENT", data.NewIntValue(8))
	vm.SetConstant("CachingIterator::TOSTRING_USE_INNER", data.NewIntValue(16))
	vm.SetConstant("CachingIterator::FULL_CACHE", data.NewIntValue(256))

	vm.SetConstant("RecursiveArrayIterator::CHILD_ARRAYS_ONLY", data.NewIntValue(4))

	vm.SetConstant("SplDoublyLinkedList::IT_MODE_FIFO", data.NewIntValue(SplITModeFIFO))
	vm.SetConstant("SplDoublyLinkedList::IT_MODE_LIFO", data.NewIntValue(SplITModeLIFO))
	vm.SetConstant("SplDoublyLinkedList::IT_MODE_KEEP", data.NewIntValue(SplITModeKeep))
	vm.SetConstant("SplDoublyLinkedList::IT_MODE_DELETE", data.NewIntValue(SplITModeDelete))
	vm.SetConstant("SplQueue::IT_MODE_FIFO", data.NewIntValue(SplITModeFIFO))
	vm.SetConstant("SplQueue::IT_MODE_DELETE", data.NewIntValue(SplITModeDelete))

	vm.SetConstant("FilesystemIterator::CURRENT_AS_PATHNAME", data.NewIntValue(FSI_CURRENT_AS_PATHNAME))
	vm.SetConstant("FilesystemIterator::CURRENT_AS_FILEINFO", data.NewIntValue(FSI_CURRENT_AS_FILEINFO))
	vm.SetConstant("FilesystemIterator::CURRENT_AS_SELF", data.NewIntValue(FSI_CURRENT_AS_SELF))
	vm.SetConstant("FilesystemIterator::KEY_AS_PATHNAME", data.NewIntValue(FSI_KEY_AS_PATHNAME))
	vm.SetConstant("FilesystemIterator::KEY_AS_FILENAME", data.NewIntValue(FSI_KEY_AS_FILENAME))
	vm.SetConstant("FilesystemIterator::FOLLOW_SYMLINKS", data.NewIntValue(FSI_FOLLOW_SYMLINKS))
	vm.SetConstant("FilesystemIterator::SKIP_DOTS", data.NewIntValue(FSI_SKIP_DOTS))
	vm.SetConstant("FilesystemIterator::UNIX_PATHS", data.NewIntValue(FSI_UNIX_PATHS))
	vm.SetConstant("FilesystemIterator::NEW_CURRENT_AND_KEY", data.NewIntValue(FSI_NEW_CURRENT_AND_KEY))

	vm.SetConstant("SplFileObject::DROP_NEW_LINE", data.NewIntValue(SFO_DROP_NEW_LINE))
	vm.SetConstant("SplFileObject::READ_AHEAD", data.NewIntValue(SFO_READ_AHEAD))
	vm.SetConstant("SplFileObject::SKIP_EMPTY", data.NewIntValue(SFO_SKIP_EMPTY))
	vm.SetConstant("SplFileObject::READ_CSV", data.NewIntValue(SFO_READ_CSV))
}
