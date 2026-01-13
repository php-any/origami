<?php

// 测试 basename 函数

// 1. 基本用法
echo "1. Basic usage:\n";
echo basename("/path/to/file.txt") . "\n"; // 输出: file.txt
echo basename("/path/to/file.txt", ".txt") . "\n"; // 输出: file
echo basename("/path/to/file.txt", "txt") . "\n"; // 输出: file.

// 2. Windows 路径
echo "\n2. Windows path:\n";
echo basename("C:\\path\\to\\file.txt") . "\n"; // 输出: file.txt
echo basename("C:\\path\\to\\file.txt", ".txt") . "\n"; // 输出: file

// 3. 相对路径
echo "\n3. Relative path:\n";
echo basename("file.txt") . "\n"; // 输出: file.txt
echo basename("./file.txt") . "\n"; // 输出: file.txt
echo basename("../file.txt") . "\n"; // 输出: file.txt

// 4. 目录路径
echo "\n4. Directory path:\n";
echo basename("/path/to/dir/") . "\n"; // 输出: dir
echo basename("/path/to/dir") . "\n"; // 输出: dir

// 5. 空字符串和特殊情况
echo "\n5. Edge cases:\n";
echo basename("") . "\n"; // 输出: (空字符串)
echo basename("/") . "\n"; // 输出: (空字符串)
echo basename("file.txt", ".txt") . "\n"; // 输出: file

// 6. 多个后缀匹配
echo "\n6. Suffix matching:\n";
echo basename("file.txt.txt", ".txt") . "\n"; // 输出: file.txt (只移除最后一个匹配的后缀)

echo "\nAll basename tests completed!\n";
