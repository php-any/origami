<?php

/**
 * 命名空间内 bare instanceof ClassName 应优先匹配当前命名空间类，
 * 不能被同名的全局类抢占（PHP 类名无全局 fallback）。
 * 复现：tests/basic/inheritance_advanced.php 的全局 BaseClass
 * 与 tests/basic/static_return_type.php 的 namespaced BaseClass 冲突。
 */

// 全局同名类（无命名空间）
class InstanceofNsVsGlobal_Base
{
}

namespace tests\php;

class InstanceofNsVsGlobal_Base
{
}

$obj = new InstanceofNsVsGlobal_Base();

if (!($obj instanceof InstanceofNsVsGlobal_Base)) {
    Log::fatal('命名空间内 instanceof 被全局同名类抢占');
}
if ($obj instanceof \InstanceofNsVsGlobal_Base) {
    Log::fatal('对象不应 instanceof 全局同名类');
}
if (!($obj instanceof \tests\php\InstanceofNsVsGlobal_Base)) {
    Log::fatal('FQCN instanceof 失败');
}

Log::info('instanceof_ns_vs_global 测试通过');
