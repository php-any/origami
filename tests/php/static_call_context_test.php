<?php

namespace tests\php;

/**
 * 静态调用上下文：后期静态绑定、默认参数、parent::、isset($this)、__callStatic。
 */

class StaticCall_Base
{
    public static $shared = 'base';
    const K = 'base-k';

    public static function who()
    {
        return 'Base';
    }

    public static function callWho()
    {
        return static::who();
    }

    public static function selfWho()
    {
        return self::who();
    }

    public static function withDefault($a, $b = 'def')
    {
        return "$a|$b";
    }

    public static function readConst()
    {
        return self::K;
    }

    public static function readStatic()
    {
        return static::$shared;
    }
}

class StaticCall_Child extends StaticCall_Base
{
    public static $shared = 'child';
    const K = 'child-k';

    public static function who()
    {
        return 'Child';
    }
}

if (StaticCall_Base::who() !== 'Base' || StaticCall_Child::who() !== 'Child') {
    \Log::fatal('静态方法解析错误');
}
if (StaticCall_Child::callWho() !== 'Child') {
    \Log::fatal('static:: 后期静态绑定应解析到 Child');
}
if (StaticCall_Child::selfWho() !== 'Base') {
    \Log::fatal('self:: 应词法绑定到 Base');
}
if (StaticCall_Base::withDefault('x') !== 'x|def' || StaticCall_Base::withDefault('x', 'y') !== 'x|y') {
    \Log::fatal('静态方法默认参数未补齐');
}
if (StaticCall_Child::readConst() !== 'base-k') {
    \Log::fatal('self::K 应为词法 Base 常量');
}
if (StaticCall_Child::readStatic() !== 'child') {
    \Log::fatal('static::$shared 应为 Child');
}

class StaticCall_P
{
    public static function f()
    {
        return 'P';
    }
}

class StaticCall_Q extends StaticCall_P
{
    public static function f()
    {
        return 'Q+' . parent::f();
    }
}

if (StaticCall_Q::f() !== 'Q+P') {
    \Log::fatal('parent:: 静态调用错误');
}

class StaticCall_R
{
    public $prop = 'v';

    public static function bad()
    {
        return isset($this);
    }
}

if (StaticCall_R::bad() !== false) {
    \Log::fatal('静态上下文 isset($this) 应为 false');
}

class StaticCall_M
{
    public static function __callStatic($name, $args)
    {
        return $name . ':' . implode(',', $args);
    }
}

if (StaticCall_M::anything('a', 'b') !== 'anything:a,b') {
    \Log::fatal('__callStatic 行为变化: ' . var_export(StaticCall_M::anything('a', 'b'), true));
}

if (call_user_func(['tests\\php\\StaticCall_Base', 'who']) !== 'Base') {
    \Log::fatal('call_user_func 数组回调静态方法失败');
}

$c = new StaticCall_Child();
if ($c::who() !== 'Child') {
    \Log::fatal('实例上调用静态方法失败');
}

\Log::info('static_call_context 测试通过');
