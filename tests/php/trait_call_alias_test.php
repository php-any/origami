<?php

namespace tests\php;

/**
 * 回归：use Trait { __call as macroCall; } 时别名必须指向 trait 方法，
 * 不能错误绑定到类自己覆盖的 __call（否则会无限递归，如 Illuminate\View\View）。
 */

trait TraitCallAliasMacroable
{
    public function __call($method, $parameters)
    {
        return 'trait:' . $method;
    }
}

class TraitCallAliasView
{
    use TraitCallAliasMacroable {
        __call as macroCall;
    }

    public function __call($method, $parameters)
    {
        return 'class:' . $this->macroCall($method, $parameters);
    }
}

$v = new TraitCallAliasView();
$got = $v->foo();
if ($got !== 'class:trait:foo') {
    Log::fatal('trait __call 别名错误: expected class:trait:foo, got ' . var_export($got, true));
}

Log::info('trait __call as macroCall 别名测试通过');
