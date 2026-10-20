<?php

namespace tests\php;

/**
 * 没有 __isset 时，isset/empty/?? 的属性下标必须调用 __get。
 * 裸 isset($obj->prop) 仍然不能调用 __get。
 */

class IssetMagicDimension_Bag
{
    public int $gets = 0;

    public function __get($key)
    {
        $this->gets++;
        if ($key === 'missing') {
            return null;
        }
        return [
            'view' => ['paths' => ['x']],
            'empty' => [],
            'nullish' => null,
        ];
    }
}

$bag = new IssetMagicDimension_Bag();

if (isset($bag->config)) {
    Log::fatal('裸 isset 应为 false');
}
if ($bag->gets !== 0) {
    Log::fatal('裸 isset 不应调用 __get');
}

if (!isset($bag->config['view']['paths'])) {
    Log::fatal('下标 isset 应为 true');
}
if ($bag->gets !== 1) {
    Log::fatal('下标 isset 应只调用一次 __get, gets=' . $bag->gets);
}

if (isset($bag->config['nullish'])) {
    Log::fatal('null 键 isset 应为 false');
}
if (isset($bag->missing['view'])) {
    Log::fatal('__get 返回 null 时下标 isset 应为 false');
}

if (empty($bag->config['view']['paths'])) {
    Log::fatal('非空下标 empty 应为 false');
}
if (!empty($bag->config['empty'])) {
    Log::fatal('空数组下标 empty 应为 true');
}
if (!empty($bag->config['nullish'])) {
    Log::fatal('null 键 empty 应为 true');
}

$paths = $bag->config['view']['paths'] ?? 'fallback';
if ($paths !== ['x']) {
    Log::fatal('?? 应返回 __get 里的下标值');
}
if (($bag->missing['view'] ?? 'fallback') !== 'fallback') {
    Log::fatal('__get 返回 null 时 ?? 应走右侧');
}
if (($bag->config['nullish'] ?? 'fallback') !== 'fallback') {
    Log::fatal('null 键 ?? 应走右侧');
}

class IssetMagicDimension_WithIsset
{
    public int $gets = 0;

    public function __isset($key)
    {
        return false;
    }

    public function __get($key)
    {
        $this->gets++;
        return ['view' => ['paths' => ['x']]];
    }
}

$blocked = new IssetMagicDimension_WithIsset();
if (isset($blocked->config['view']['paths']) || $blocked->gets !== 0) {
    Log::fatal('__isset 返回 false 时下标 isset 不能调用 __get');
}

Log::info('isset 属性下标 __get 测试通过');
