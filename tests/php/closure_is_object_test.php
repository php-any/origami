<?php

namespace tests\php;

/**
 * Closure 是 object：可传给 object 类型参数，且 spl_object_hash 可用。
 */

$fn = function () {
    return 1;
};

if (!($fn instanceof \Closure)) {
    \Log::fatal('闭包应 instanceof Closure');
}

$hash = spl_object_hash($fn);
if (!is_string($hash) || $hash === '') {
    \Log::fatal('spl_object_hash(Closure) 应返回非空字符串');
}

$id = spl_object_id($fn);
if (!is_int($id)) {
    \Log::fatal('spl_object_id(Closure) 应返回 int');
}

function ObjectTypeAcceptsClosure_take(object $o): bool
{
    return true;
}

if (!ObjectTypeAcceptsClosure_take($fn)) {
    \Log::fatal('object 类型参数应接受 Closure');
}

\Log::info('closure_is_object 测试通过');
