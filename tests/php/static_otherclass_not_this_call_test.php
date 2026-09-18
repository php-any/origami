<?php

namespace tests\php;

/**
 * 实例方法里调用「另一个类」上不存在的静态方法时，不得落到当前对象的 __call。
 * Laravel Collection::mapWithKeys 内部 Arr::mapWithKeys 依赖此语义。
 */

class StaticOtherNotThisCall_Other
{
}

class StaticOtherNotThisCall_Host
{
    public function __call($method, $args)
    {
        return 'magic:'.$method;
    }

    public function run()
    {
        return StaticOtherNotThisCall_Other::noSuchMethod();
    }
}

$h = new StaticOtherNotThisCall_Host();
$msg = '';
try {
    $got = $h->run();
    if (is_string($got) && strpos($got, 'magic:') === 0) {
        Log::fatal('OtherClass::noSuchMethod 错误地走了 Host::__call: '.$got);
    }
    Log::fatal('OtherClass::noSuchMethod 应抛错，实际: '.var_export($got, true));
} catch (\Throwable $e) {
    $msg = $e->getMessage();
}

if ($msg === '' || strpos($msg, 'noSuchMethod') === false) {
    Log::fatal('应报告 Other 上缺少 noSuchMethod，实际: '.$msg);
}
if (strpos($msg, 'magic:') !== false) {
    Log::fatal('错误地走了 Host::__call: '.$msg);
}

Log::info('static_otherclass_not_this_call 测试通过');
