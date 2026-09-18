<?php

namespace tests\php;

/**
 * catch ($e) 必须绑定 throw 的异常实例，而不是内部控制流包装。
 * Laravel Handler 用 $e instanceof AuthenticationException + $e->redirectTo()。
 */

class CatchExObj_Auth extends \Exception
{
    public $redirectTo = '/admin/login';

    public function redirectTo()
    {
        return $this->redirectTo;
    }
}

try {
    throw new CatchExObj_Auth('Unauthenticated.');
} catch (\Throwable $e) {
    if (!($e instanceof CatchExObj_Auth)) {
        \Log::fatal('catch 变量 instanceof 子类失败, get_class=' . get_class($e));
    }
    if (!($e instanceof \Exception)) {
        \Log::fatal('catch 变量 instanceof Exception 失败');
    }
    if (!is_a($e, CatchExObj_Auth::class)) {
        \Log::fatal('is_a 子类失败');
    }
    if ($e->getMessage() !== 'Unauthenticated.') {
        \Log::fatal('getMessage 失败: ' . $e->getMessage());
    }
    if ($e->redirectTo() !== '/admin/login') {
        \Log::fatal('子类方法失败: ' . var_export($e->redirectTo(), true));
    }
    $kind = match (true) {
        $e instanceof CatchExObj_Auth => 'auth',
        default => 'other',
    };
    if ($kind !== 'auth') {
        \Log::fatal('match(true)+instanceof 失败: ' . $kind);
    }
}

\Log::info('catch_exception_object 测试通过');
