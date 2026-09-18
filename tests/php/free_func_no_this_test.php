<?php

namespace tests\php;

/**
 * 实例方法里调用的普通函数/文件级闭包不得继承 $this；方法内闭包仍绑定定义处 $this。
 */

class FreeFuncNoThis_Box
{
    public $n = 7;

    public function viaHelper()
    {
        return free_func_no_this_helper();
    }

    public function viaFileClosure()
    {
        $fn = free_func_no_this_file_closure();
        return $fn();
    }

    public function viaMethodClosure()
    {
        $fn = function () {
            return isset($this) ? $this->n : -1;
        };
        return $fn();
    }
}

function free_func_no_this_helper()
{
    return isset($this);
}

function free_func_no_this_file_closure()
{
    return function () {
        return isset($this);
    };
}

$o = new FreeFuncNoThis_Box();
if ($o->viaHelper() !== false) {
    \Log::fatal('实例方法调用的普通函数不应带 $this');
}
if ($o->viaFileClosure() !== false) {
    \Log::fatal('文件级闭包从方法里调用不应带 $this');
}
if ($o->viaMethodClosure() !== 7) {
    \Log::fatal('方法内闭包应绑定定义处 $this');
}

\Log::info('free_func_no_this 测试通过');
