<?php

namespace tests\php;

/**
 * SPL Phase 4 测试：SplFileObject、SplTempFileObject、GlobIterator
 */

$testFile = tempnam('', 'origami_sfo_test_');
file_put_contents($testFile, "line0\nline1\nline2\n");
$tmpDir = dirname($testFile);

// ---- SplFileObject 常量 ----
if (\SplFileObject::DROP_NEW_LINE !== 1) {
    Log::fatal('spl_file_object_test: DROP_NEW_LINE 应为 1');
}
if (\SplFileObject::READ_AHEAD !== 2) {
    Log::fatal('spl_file_object_test: READ_AHEAD 应为 2');
}
Log::info('SplFileObject 常量测试通过');

// ---- SplFileObject 迭代 ----
$file = new \SplFileObject($testFile, 'r', \SplFileObject::DROP_NEW_LINE);
$lines = [];
foreach ($file as $line) {
    $lines[] = $line;
}
if ($lines !== ['line0', 'line1', 'line2']) {
    Log::fatal('spl_file_object_test: foreach 行内容错误: ' . json_encode($lines));
}
Log::info('SplFileObject foreach 测试通过');

// ---- fgets / eof ----
$file2 = new \SplFileObject($testFile, 'r', \SplFileObject::DROP_NEW_LINE);
$file2->rewind();
if ($file2->fgets() !== 'line0') {
    Log::fatal('spl_file_object_test: fgets 首行错误');
}
if ($file2->eof()) {
    Log::fatal('spl_file_object_test: 读取首行后 eof 应为 false');
}
Log::info('SplFileObject fgets/eof 测试通过');

// ---- seek / key ----
$file3 = new \SplFileObject($testFile, 'r', \SplFileObject::DROP_NEW_LINE | \SplFileObject::READ_AHEAD);
$file3->seek(1);
if ($file3->key() !== 1 || $file3->current() !== 'line1') {
    Log::fatal('spl_file_object_test: seek(1) 失败 key=' . $file3->key() . ' current=' . $file3->current());
}
Log::info('SplFileObject seek 测试通过');

// ---- fgetcsv / fputcsv ----
$csvFile = tempnam('', 'origami_sfo_csv_');
$csv = new \SplFileObject($csvFile, 'w+');
$written = $csv->fputcsv(['a', 'b', 'c']);
if ($written === false) {
    Log::fatal('spl_file_object_test: fputcsv 失败');
}
$csv->rewind();
$row = $csv->fgetcsv();
if (!is_array($row) || $row[0] !== 'a' || $row[1] !== 'b' || $row[2] !== 'c') {
    Log::fatal('spl_file_object_test: fgetcsv 结果错误: ' . json_encode($row));
}
Log::info('SplFileObject fgetcsv/fputcsv 测试通过');

// ---- hasChildren ----
$file4 = new \SplFileObject($testFile);
if ($file4->hasChildren()) {
    Log::fatal('spl_file_object_test: SplFileObject hasChildren 应为 false');
}
Log::info('SplFileObject hasChildren 测试通过');

// ---- SplTempFileObject ----
$temp = new \SplTempFileObject('php://temp', 'w+b', \SplFileObject::DROP_NEW_LINE);
$temp->fwrite("temp-data\n");
$temp->rewind();
if ($temp->fgets() !== 'temp-data') {
    Log::fatal('spl_file_object_test: SplTempFileObject 读写失败');
}
if ($temp->getPathname() === '') {
    Log::fatal('spl_file_object_test: SplTempFileObject getPathname 不应为空');
}
Log::info('SplTempFileObject 测试通过');

// ---- GlobIterator ----
$globPattern = $tmpDir . DIRECTORY_SEPARATOR . 'origami_sfo_test_*';
// 至少匹配 $testFile
$gi = new \GlobIterator($globPattern);
$found = [];
foreach ($gi as $entry) {
    $found[] = $entry->getPathname();
}
if (!in_array($testFile, $found, true)) {
    Log::fatal('spl_file_object_test: GlobIterator 未找到测试文件');
}
Log::info('GlobIterator 测试通过');

// ---- spl_classes 包含新类 ----
$classes = spl_classes();
foreach (['SplFileObject', 'SplTempFileObject', 'GlobIterator'] as $cls) {
    if (!isset($classes[$cls])) {
        Log::fatal('spl_file_object_test: spl_classes 缺少 ' . $cls);
    }
}
Log::info('spl_classes 新类注册测试通过');

@unlink($testFile);
@unlink($csvFile);
@unlink($temp->getPathname());

Log::info('SPL Phase 4 测试全部通过');
