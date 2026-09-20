<?php

namespace tests\php;

/**
 * Filament cacheSchema：new ReflectionMethod($this, $methodName) 出现在带 int 形参/局部的方法里。
 * 槽 0 若被当成 $this 会得到 *IntValue。
 */

class ReflMethodIntLocals_Host
{
    public function content($schema = null)
    {
        return 'ok';
    }

    public function cacheLike(int $n, string $name = 'content')
    {
        $count = $n;
        $methodName = $name;
        if (!method_exists($this, $methodName)) {
            Log::fatal('method_exists($this, content) 失败, n=' . $count);
        }
        $m = new \ReflectionMethod($this, $methodName);
        return $m->getName();
    }
}

$o = new ReflMethodIntLocals_Host();
$got = $o->cacheLike(6);
if ($got !== 'content') {
    Log::fatal('ReflectionMethod($this) 在 int 局部变量下失败: ' . json_encode($got));
}

$m2 = new \ReflectionMethod($o, 'content');
if ($m2->getName() !== 'content') {
    Log::fatal('ReflectionMethod($obj, content) 失败');
}

Log::info('ReflectionMethod($this) 在 int 局部变量下测试通过');
