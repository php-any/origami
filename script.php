<?php

// 测试 DirectoryIterator 的 foreach 循环

echo "=== DirectoryIterator foreach 测试 ===\n\n";

// 1. 基本 foreach 循环（只有值）
echo "1. 基本 foreach 循环（只有值）:\n";
$iterator = new DirectoryIterator("tests/php");
$count = 0;
foreach ($iterator as $value) {
    $count++;
    if ($count > 5) {
        break;
    }
    echo "  " . $count . ": " . $value . "\n";
}

// 2. foreach 循环（键和值）
echo "\n2. foreach 循环（键和值）:\n";
$iterator2 = new DirectoryIterator("tests/php");
$count2 = 0;
foreach ($iterator2 as $key => $value) {
    $count2++;
    if ($count2 > 5) {
        break;
    }
    echo "  [" . $key . "] => " . $value . "\n";
}

// 3. foreach 循环中使用 DirectoryIterator 方法
echo "\n3. foreach 循环中使用 DirectoryIterator 方法:\n";
$iterator3 = new DirectoryIterator("tests/php");
$count3 = 0;
foreach ($iterator3 as $key => $value) {
    $count3++;
    if ($count3 > 5) {
        break;
    }
    $filename = $iterator3->getFilename();
    $isDir = $iterator3->isDir();
    $isFile = $iterator3->isFile();
    $isDot = $iterator3->isDot();
    $type = $isDir ? "目录" : ($isFile ? "文件" : "未知");
    $dot = $isDot ? " (.)" : "";
    echo "  [" . $key . "] " . $filename . " (" . $type . ")" . $dot . "\n";
}

// 4. 遍历所有文件（跳过 . 和 ..）
echo "\n4. 遍历所有文件（跳过 . 和 ..）:\n";
$iterator4 = new DirectoryIterator("tests/php");
foreach ($iterator4 as $key => $value) {
    if ($iterator4->isDot()) {
        continue;
    }
    $filename = $iterator4->getFilename();
    $pathname = $iterator4->getPathname();
    echo "  " . $filename . " -> " . $pathname . "\n";
}

// 5. 只遍历目录
echo "\n5. 只遍历目录:\n";
$iterator5 = new DirectoryIterator("tests");
foreach ($iterator5 as $key => $value) {
    if ($iterator5->isDot()) {
        continue;
    }
    if ($iterator5->isDir()) {
        echo "  目录: " . $iterator5->getFilename() . "\n";
    }
}

// 6. 只遍历文件
echo "\n6. 只遍历文件:\n";
$iterator6 = new DirectoryIterator("tests/php");
$fileCount = 0;
foreach ($iterator6 as $key => $value) {
    if ($iterator6->isDot()) {
        continue;
    }
    if ($iterator6->isFile()) {
        $fileCount++;
        if ($fileCount > 5) {
            break;
        }
        echo "  文件: " . $iterator6->getFilename() . "\n";
    }
}

echo "\n=== DirectoryIterator foreach 测试完成 ===\n";
