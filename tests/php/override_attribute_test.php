<?php

namespace tests\php;

/**
 * PHP 8.3 Override 属性类可加载，且可用于方法注解。
 */

if (!class_exists('Override', false)) {
    \Log::fatal('Override 属性类未注册');
}

class OverrideAttr_Base
{
    public function hello(): string
    {
        return 'base';
    }
}

class OverrideAttr_Child extends OverrideAttr_Base
{
    #[\Override]
    public function hello(): string
    {
        return 'child';
    }
}

$c = new OverrideAttr_Child();
if ($c->hello() !== 'child') {
    \Log::fatal('Override 注解方法调用失败');
}

\Log::info('override_attribute 测试通过');
