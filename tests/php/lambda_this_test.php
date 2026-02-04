<?php

namespace tests\php;

use Exception;

/**
 * 验证 LambdaExpression 在类方法中定义后：
 * 1) 作为普通闭包在全局调用时，$this 仍然指向原对象；
 * 2) 作为 set_exception_handler 的回调时，$this 也能正确指向原对象。
 */

class LambdaThisTest
{
    public function getHandler(): callable
    {
        // 在类方法中定义 lambda，并直接使用 $this
        return fn () => Log::info('lambda-handler this=' . static::class . ' actual=' . get_class($this));
    }

    public function registerExceptionHandlerAndThrow(): void
    {
        // 使用箭头函数作为 set_exception_handler 的回调，并访问 $this
        set_exception_handler(fn (\Throwable $e) => Log::info(
            'exception-handler this=' . static::class .
            ' actual=' . get_class($this) .
            ' msg=' . $e->getMessage()
        ));

        throw new Exception('lambda this test exception');
    }
}

$obj = new LambdaThisTest();

// 场景一：从方法中拿到 lambda，在全局执行，检查 $this 是否仍然指向 $obj
$handler = $obj->getHandler();
$handler();

// 场景二：在方法内部注册 set_exception_handler，抛出未捕获异常，看回调里 $this 是否正确
try {
    $obj->registerExceptionHandlerAndThrow();
} catch (Exception $e) {
    // 理论上 set_exception_handler 应该接管未捕获异常，这里只是兜底，避免测试环境意外终止
    Log::info('fallback-catch: ' . $e->getMessage());
}

