<?php

namespace tests\php;

use Exception;

/**
 * set_exception_handler 对可变参数箭头函数（Laravel HandleExceptions::forwardsTo）应传入 [$e]。
 */

class ExcHandler_VariadicForward
{
    public function handleException($e)
    {
        if ($e->getMessage() !== 'variadic-handler-ok') {
            Log::fatal('可变参数异常处理器消息错误: ' . $e->getMessage());
        }
        Log::info('set_exception_handler_variadic 测试通过');
    }

    protected function forwardsTo($method)
    {
        return fn (...$arguments) => $this->{$method}(...$arguments);
    }

    public function install()
    {
        set_exception_handler($this->forwardsTo('handleException'));
    }
}

$h = new ExcHandler_VariadicForward();
$h->install();

throw new Exception('variadic-handler-ok');
