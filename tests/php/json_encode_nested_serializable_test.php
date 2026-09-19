<?php

namespace tests\php;

/**
 * json_encode 数组里的 JsonSerializable 应编成 JSON 值，而不是 Object(Class) 字符串。
 * Telescope /telescope-api/requests 的 entries 是 Collection。
 */
class JsonEncNested_Entries implements \JsonSerializable
{
    public function jsonSerialize(): mixed
    {
        return [
            ['id' => 1],
        ];
    }
}

$json = json_encode([
    'entries' => new JsonEncNested_Entries(),
    'status' => 'enabled',
]);

if ($json === false) {
    Log::fatal('json_encode 失败: ' . json_last_error_msg());
}
if (str_contains($json, 'Object(')) {
    Log::fatal('json_encode 把对象编成 Object() 字符串: ' . $json);
}
if (!str_contains($json, '"entries"') || !str_contains($json, '"id"') || !str_contains($json, '"status":"enabled"')) {
    Log::fatal('json_encode 嵌套 JsonSerializable 失败: ' . $json);
}

Log::info('json_encode 嵌套 JsonSerializable 测试通过: ' . $json);
