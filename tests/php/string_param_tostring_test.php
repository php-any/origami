<?php

namespace tests\php;

/**
 * 弱类型下 string 参数应接受带 __toString 的对象（Filament Notification::make(Str::orderedUuid())）。
 */

class StringParam_UuidLike
{
    public function __toString(): string
    {
        return 'uuid-like-ok';
    }
}

function StringParam_accept(string $id): string
{
    return $id;
}

$got = StringParam_accept(new StringParam_UuidLike());
if ($got !== 'uuid-like-ok') {
    Log::fatal('string 参数未通过 __toString 转换: ' . var_export($got, true));
}

Log::info('string 参数 __toString 测试通过');
