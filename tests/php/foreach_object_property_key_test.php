<?php

namespace tests\php;

/**
 * foreach ($arr as $this->k => $this->v) 必须写入对象属性，不能因 GetIndex panic。
 */

class ForeachObjectPropertyKey_Box
{
    public $k = null;
    public $v = null;

    public function run(array $arr)
    {
        foreach ($arr as $this->k => $this->v) {
        }
        return [$this->k, $this->v];
    }
}

$box = new ForeachObjectPropertyKey_Box();
[$k, $v] = $box->run(['a' => 1, 'b' => 2]);
if ($k !== 'b' || $v !== 2) {
    \Log::fatal('foreach 写入 $this->k/$this->v 失败 k='.var_export($k, true).' v='.var_export($v, true));
}

$obj = new ForeachObjectPropertyKey_Box();
foreach (['x' => 10] as $obj->k => $obj->v) {
}
if ($obj->k !== 'x' || $obj->v !== 10) {
    \Log::fatal('foreach 写入 $obj->k/$obj->v 失败');
}

\Log::info('foreach 对象属性键值测试通过');
