<?php

namespace tests\php;

/**
 * file_put_contents 应支持 FILE_APPEND，否则多次写入会互相覆盖。
 */
$path = sys_get_temp_dir().DIRECTORY_SEPARATOR.'origami_file_append_test.txt';
@unlink($path);
file_put_contents($path, "a\n");
file_put_contents($path, "b\n", FILE_APPEND);
file_put_contents($path, "c\n", FILE_APPEND);
$got = file_get_contents($path);
@unlink($path);
if ($got !== "a\nb\nc\n") {
    Log::fatal('FILE_APPEND 失败 got='.var_export($got, true).' FILE_APPEND='.var_export(FILE_APPEND, true));
}
Log::info('file_append 测试通过');
