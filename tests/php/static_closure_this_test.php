<?php

namespace tests\php;

/**
 * static 闭包不得绑定定义处的 $this，但 self:: 仍可用；普通闭包与 Closure::bind 仍可用 $this。
 */
class StaticClosureThis_FS
{
    public $marker = 'filesystem';

    private const LABEL = 'self-ok';

    public function viaStatic()
    {
        return (static function () {
            return isset($this);
        })();
    }

    public function viaStaticSelf()
    {
        return (static function () {
            return self::LABEL;
        })();
    }

    public function viaNormal()
    {
        return (function () {
            return $this->marker;
        })();
    }

    public function viaBind($obj)
    {
        return \Closure::bind(function () {
            return $this->content;
        }, $obj, $obj)();
    }
}

class StaticClosureThis_Comp
{
    public $content = 'from-comp';
}

$fs = new StaticClosureThis_FS();

if ($fs->viaStatic() !== false) {
    Log::fatal('static 闭包 isset($this) 应为 false，实际: ' . var_export($fs->viaStatic(), true));
}

if ($fs->viaStaticSelf() !== 'self-ok') {
    Log::fatal('static 闭包 self:: 失败: ' . var_export($fs->viaStaticSelf(), true));
}

if ($fs->viaNormal() !== 'filesystem') {
    Log::fatal('普通闭包应保留 $this，实际: ' . var_export($fs->viaNormal(), true));
}

$c = new StaticClosureThis_Comp();
if ($fs->viaBind($c) !== 'from-comp') {
    Log::fatal('Closure::bind 失败: ' . var_export($fs->viaBind($c), true));
}

Log::info('static_closure_this 测试通过');
