<?php

echo "=== basename() 函数测试 ===\n";

// 1. 基本用法测试
if(basename("/path/to/file.txt") == "file.txt") {
    Log::info("基本用法测试通过");
} else {
    Log::fatal("基本用法测试失败，结果: " . basename("/path/to/file.txt"));
}

// 2. 带后缀参数测试
if(basename("/path/to/file.txt", ".txt") == "file") {
    Log::info("带后缀参数测试通过");
} else {
    Log::fatal("带后缀参数测试失败，结果: " . basename("/path/to/file.txt", ".txt"));
}

// 3. 后缀参数不匹配测试
if(basename("/path/to/file.txt", ".php") == "file.txt") {
    Log::info("后缀参数不匹配测试通过");
} else {
    Log::fatal("后缀参数不匹配测试失败，结果: " . basename("/path/to/file.txt", ".php"));
}

// 4. Windows 路径测试
if(basename("C:\\path\\to\\file.txt") == "file.txt") {
    Log::info("Windows 路径测试通过");
} else {
    Log::fatal("Windows 路径测试失败，结果: " . basename("C:\\path\\to\\file.txt"));
}

// 5. Windows 路径带后缀测试
if(basename("C:\\path\\to\\file.txt", ".txt") == "file") {
    Log::info("Windows 路径带后缀测试通过");
} else {
    Log::fatal("Windows 路径带后缀测试失败，结果: " . basename("C:\\path\\to\\file.txt", ".txt"));
}

// 6. 相对路径测试
if(basename("file.txt") == "file.txt") {
    Log::info("相对路径测试通过");
} else {
    Log::fatal("相对路径测试失败，结果: " . basename("file.txt"));
}

// 7. 带相对路径前缀测试
if(basename("./file.txt") == "file.txt") {
    Log::info("带相对路径前缀测试通过");
} else {
    Log::fatal("带相对路径前缀测试失败，结果: " . basename("./file.txt"));
}

// 8. 目录路径测试（以斜杠结尾）
if(basename("/path/to/dir/") == "dir") {
    Log::info("目录路径测试（以斜杠结尾）通过");
} else {
    Log::fatal("目录路径测试（以斜杠结尾）失败，结果: " . basename("/path/to/dir/"));
}

// 9. 目录路径测试（不以斜杠结尾）
if(basename("/path/to/dir") == "dir") {
    Log::info("目录路径测试（不以斜杠结尾）通过");
} else {
    Log::fatal("目录路径测试（不以斜杠结尾）失败，结果: " . basename("/path/to/dir"));
}

// 10. 空字符串测试
if(basename("") == "") {
    Log::info("空字符串测试通过");
} else {
    Log::fatal("空字符串测试失败，结果: '" . basename("") . "'");
}

// 11. 只有斜杠测试
if(basename("/") == "") {
    Log::info("只有斜杠测试通过");
} else {
    Log::fatal("只有斜杠测试失败，结果: '" . basename("/") . "'");
}

// 12. 多个斜杠测试
if(basename("/path//to///file.txt") == "file.txt") {
    Log::info("多个斜杠测试通过");
} else {
    Log::fatal("多个斜杠测试失败，结果: " . basename("/path//to///file.txt"));
}

// 13. 后缀参数部分匹配测试
if(basename("/path/to/file.txt", "txt") == "file.") {
    Log::info("后缀参数部分匹配测试通过");
} else {
    Log::fatal("后缀参数部分匹配测试失败，结果: " . basename("/path/to/file.txt", "txt"));
}

// 14. 多个后缀匹配测试（只移除最后一个）
if(basename("/path/to/file.txt.txt", ".txt") == "file.txt") {
    Log::info("多个后缀匹配测试通过");
} else {
    Log::fatal("多个后缀匹配测试失败，结果: " . basename("/path/to/file.txt.txt", ".txt"));
}

// 15. 无扩展名文件测试
if(basename("/path/to/file") == "file") {
    Log::info("无扩展名文件测试通过");
} else {
    Log::fatal("无扩展名文件测试失败，结果: " . basename("/path/to/file"));
}

// 16. 无扩展名文件带后缀测试
if(basename("/path/to/file", ".txt") == "file") {
    Log::info("无扩展名文件带后缀测试通过");
} else {
    Log::fatal("无扩展名文件带后缀测试失败，结果: " . basename("/path/to/file", ".txt"));
}

// 17. 只有文件名测试
if(basename("file.txt") == "file.txt") {
    Log::info("只有文件名测试通过");
} else {
    Log::fatal("只有文件名测试失败，结果: " . basename("file.txt"));
}

// 18. 根目录文件测试
if(basename("/file.txt") == "file.txt") {
    Log::info("根目录文件测试通过");
} else {
    Log::fatal("根目录文件测试失败，结果: " . basename("/file.txt"));
}

echo "=== basename() 测试完成 ===\n";

