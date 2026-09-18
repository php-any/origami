<?php

namespace tests\fileinfo;

/**
 * fileinfo：finfo_open / finfo_buffer。
 */

$f = finfo_open(FILEINFO_MIME_TYPE);
if ($f === false) {
    Log::fatal('finfo_open 失败');
}
$mime = finfo_buffer($f, '<?php echo 1;');
if (!is_string($mime) || $mime === '') {
    Log::fatal('finfo_buffer 应返回 mime 字符串');
}
finfo_close($f);

Log::info('fileinfo 测试通过');
