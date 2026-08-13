<?php

namespace tests\php;

/**
 * sprintf 混用 %s / %d 时的弱类型兼容性测试。
 *
 * 覆盖 serve 命令等场景：sprintf('http://%s:%d', $host, $port)
 */

$host = '0.0.0.0';
$port = (int) '8080';

$result = sprintf('Server running on [http://%s:%d].', $host, $port);
$expected = 'Server running on [http://0.0.0.0:8080].';

if ($result !== $expected) {
    Log::fatal('sprintf 混用 %s/%d 测试失败: 期望 '.$expected.' 实际 '.$result);
}

// 字符串端口也应被 %d 正确格式化
$result2 = sprintf('http://%s:%d', $host, '8080');
if ($result2 !== 'http://0.0.0.0:8080') {
    Log::fatal('sprintf %d 字符串参数测试失败: 实际 '.$result2);
}

Log::info('sprintf 混用 %s/%d 测试通过: '.$result);
