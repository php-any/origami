<?php

namespace tests\php;

/**
 * stream_context_create 函数测试：创建 stream-context 资源供 file_get_contents 使用。
 */

$context = stream_context_create([
    'http' => [
        'method' => 'GET',
        'header' => 'Accept: application/json',
        'ignore_errors' => true,
        'timeout' => 5,
    ],
]);

if (!is_resource($context)) {
    Log::fatal('stream_context_create 应返回 resource');
}

if (get_resource_type($context) !== 'stream-context') {
    Log::fatal('stream_context_create 资源类型应为 stream-context: ' . get_resource_type($context));
}

$empty = stream_context_create();
if (!is_resource($empty)) {
    Log::fatal('stream_context_create 无参数时应返回 resource');
}

Log::info('stream_context_create 函数测试通过');
