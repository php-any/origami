<?php

namespace tests\php;

/**
 * json_encode Collection 列表应对齐 PHP：entries 是对象数组，不能多包一层 [[...]]。
 * Telescope /telescope-api/requests 前端读 entry.content.method。
 */

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

class TelEntries_Box implements \JsonSerializable
{
    public function __construct(private string $id, private string $method)
    {
    }

    public function jsonSerialize(): mixed
    {
        return [
            'id' => $this->id,
            'content' => ['method' => $this->method, 'uri' => '/'],
        ];
    }
}

$c = new \Illuminate\Support\Collection([
    new TelEntries_Box('a', 'GET'),
    new TelEntries_Box('b', 'POST'),
]);

$nested = new \Illuminate\Support\Collection($c);
$nestedAll = $nested->all();
if (!is_array($nestedAll) || count($nestedAll) !== 2) {
    Log::fatal('new Collection($collection) 应摊平: count=' . (is_countable($nestedAll) ? count($nestedAll) : -1) . ' json=' . json_encode($nestedAll));
}
if (isset($nestedAll[0]) && $nestedAll[0] instanceof \Illuminate\Support\Collection) {
    Log::fatal('new Collection($collection) 把内层 Collection 包成了一项');
}

$all = $c->all();
if (count($all) !== 2) {
    Log::fatal('Collection::all 数量错误: ' . json_encode($all));
}
if (is_array($all[0] ?? null) && array_is_list($all[0]) && isset($all[0][0])) {
    Log::fatal('Collection::all 多包一层: ' . json_encode($all));
}

$json = json_encode(['entries' => $c, 'status' => 'enabled']);
if ($json === false) {
    Log::fatal('json_encode 失败: ' . json_last_error_msg());
}

$decoded = json_decode($json, true);
if (!is_array($decoded['entries'] ?? null)) {
    Log::fatal('entries 不是数组: ' . $json);
}
if (isset($decoded['entries'][0]) && is_array($decoded['entries'][0]) && array_is_list($decoded['entries'][0])) {
    Log::fatal('entries 多包一层数组: ' . $json);
}
if (($decoded['entries'][0]['content']['method'] ?? null) !== 'GET') {
    Log::fatal('entries[0].content.method 缺失: ' . $json);
}
if (($decoded['entries'][1]['id'] ?? null) !== 'b') {
    Log::fatal('entries[1].id 错误: ' . $json);
}

$nestedJson = json_encode(['entries' => $nested, 'status' => 'enabled']);
$nestedDecoded = json_decode($nestedJson, true);
if (isset($nestedDecoded['entries'][0]) && is_array($nestedDecoded['entries'][0]) && array_is_list($nestedDecoded['entries'][0]) && isset($nestedDecoded['entries'][0][0]['content'])) {
    Log::fatal('json_encode(new Collection($collection)) 多包一层: ' . $nestedJson);
}
if (($nestedDecoded['entries'][0]['content']['method'] ?? null) !== 'GET') {
    Log::fatal('嵌套 Collection json_encode 缺少 method: ' . $nestedJson);
}

$keepInner = new \Illuminate\Support\Collection([[1, 2, 3]]);
$keepAll = $keepInner->all();
if (!is_array($keepAll[0] ?? null)) {
    Log::fatal('collect([[1,2,3]]) 不应摊平内层数组: ' . json_encode($keepAll));
}

Log::info('Collection json_encode 列表不多层测试通过: ' . $json);
