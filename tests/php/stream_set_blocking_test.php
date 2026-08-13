<?php

$stream = fopen('php://temp', 'w+');
if (!is_resource($stream)) {
    Log::fatal('无法创建测试流');
}
if (!stream_set_blocking($stream, false)) {
    Log::fatal('stream_set_blocking(false) 失败');
}
if (!stream_set_blocking($stream, true)) {
    Log::fatal('stream_set_blocking(true) 失败');
}
$read = [$stream];
$write = $except = [];
if (stream_select($read, $write, $except, 0, 0) !== 1) {
    Log::fatal('stream_select 未报告候选流');
}
if (fseek($stream, 0, SEEK_SET) !== 0) {
    Log::fatal('fseek 失败');
}
fclose($stream);

Log::info('stream_set_blocking_test OK');
