<?php

namespace tests\php;

/**
 * isset($obj->prop) 必须走 __isset，不得调用 __get（避免 __get 内 Undefined array key）。
 */

class IssetMagic_NoGetProbe
{
    public int $gets = 0;
    protected $data = ['x' => 1];

    public function &__get($key)
    {
        $this->gets++;
        return $this->data[$key];
    }

    public function __isset($key)
    {
        return isset($this->data[$key]);
    }
}

$v = new IssetMagic_NoGetProbe();

if (isset($v->layoutConfig)) {
    Log::fatal('缺失属性 isset 应为 false');
}
if ($v->gets !== 0) {
    Log::fatal('isset 不应调用 __get, gets=' . $v->gets);
}

if (!isset($v->x)) {
    Log::fatal('存在键 x 时 isset 应为 true');
}
if ($v->gets !== 0) {
    Log::fatal('isset(x) 仍不应调用 __get, gets=' . $v->gets);
}

$got = $v->layoutConfig ?? 'DEFAULT';
if ($got !== 'DEFAULT' || $v->gets !== 0) {
    Log::fatal('?? 也不应调用 __get, got=' . var_export($got, true) . ' gets=' . $v->gets);
}

Log::info('isset/?? 不调用 __get 测试通过');
