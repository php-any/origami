<?php

namespace tests\php;

/**
 * 子类 $this->parentMethod() 在声明类解析时不得因 VM 为空而 panic。
 */

class FindDeclaringClass_Parent
{
    public function ping()
    {
        return 'p';
    }
}

class FindDeclaringClass_Child extends FindDeclaringClass_Parent
{
    public function run()
    {
        return $this->ping();
    }
}

$c = new FindDeclaringClass_Child();
$got = $c->run();
if ($got !== 'p') {
    Log::fatal('继承方法调用失败: '.$got);
}

Log::info('find_declaring_class_parent_method 测试通过');
