<?php

namespace tests\php;

/**
 * json_encode 不得就地改写数组/关联数组：元素必须仍是原对象。
 * 就地替换会在 Laravel/Livewire 多次请求后把共享结构胀成巨型字符串。
 */
class JsonEncodeNoMut_Item implements \JsonSerializable
{
    public int $n = 1;

    public function jsonSerialize(): mixed
    {
        return ['n' => $this->n];
    }
}

$item = new JsonEncodeNoMut_Item();
$list = [$item];
$assoc = ['item' => $item];

$jsonList = json_encode($list);
$jsonAssoc = json_encode($assoc);

if (strpos((string) $jsonList, '"n"') === false) {
    Log::fatal('列表 json_encode 失败: ' . var_export($jsonList, true));
}
if (strpos((string) $jsonAssoc, '"n"') === false) {
    Log::fatal('关联 json_encode 失败: ' . var_export($jsonAssoc, true));
}
if (!($list[0] instanceof JsonEncodeNoMut_Item)) {
    Log::fatal('json_encode 改写了列表元素');
}
if (!($assoc['item'] instanceof JsonEncodeNoMut_Item)) {
    Log::fatal('json_encode 改写了关联元素');
}

$again = json_encode($list);
if ($again !== $jsonList) {
    Log::fatal('再次 json_encode 应变稳定，第一次=' . $jsonList . ' 第二次=' . $again);
}

Log::info('json_encode 不改写入参测试通过');
