<?php

namespace tests\php;

/**
 * JsonSerializable 接口行为测试：
 * - 类实现 \JsonSerializable 且提供 jsonSerialize(): mixed
 * - json_encode() 应优先编码 jsonSerialize() 返回值
 */

class JsonSerializable_Test implements \JsonSerializable
{
    private string $value;

    public function __construct(string $value)
    {
        $this->value = $value;
    }

    public function jsonSerialize(): mixed
    {
        return [
            'wrapped' => $this->value,
        ];
    }
}

$obj = new JsonSerializable_Test('hello');
$json = json_encode($obj);

// 预期：{"wrapped":"hello"}（键顺序可能不同，但至少要包含此子串）
if (strpos($json, '"wrapped"') === false || strpos($json, '"hello"') === false) {
    Log::fatal('JsonSerializable + json_encode 测试失败: ' . $json);
}

Log::info('JsonSerializable + json_encode 测试通过');

