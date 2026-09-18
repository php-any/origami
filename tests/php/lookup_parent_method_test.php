<?php

namespace tests\php;

/**
 * 子类实例调用定义在父类上的方法时不应因 VM 为空而 panic。
 */

class LookupParentMethod_Base
{
    public function greet()
    {
        return 'ok';
    }
}

class LookupParentMethod_Child extends LookupParentMethod_Base
{
}

$c = new LookupParentMethod_Child();
if ($c->greet() !== 'ok') {
    Log::fatal('继承方法调用失败: '.var_export($c->greet(), true));
}

Log::info('lookup_parent_method 测试通过');
