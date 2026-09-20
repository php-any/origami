<?php

namespace tests\php;

/**
 * ReflectionMethod / ReflectionClass 必须接受实例（含 $this），对齐 PHP object|string。
 * Filament InteractsWithSchemas：new ReflectionMethod($this, $methodName)。
 */

class ReflMethodThis_Target
{
    public function ping()
    {
        return 'pong';
    }

    public function viaThis()
    {
        return (new \ReflectionMethod($this, 'ping'))->getName();
    }
}

$o = new ReflMethodThis_Target();
$viaThis = $o->viaThis();
if ($viaThis !== 'ping') {
    Log::fatal('ReflectionMethod($this, ping) 失败: ' . json_encode($viaThis));
}

$m = new \ReflectionMethod($o, 'ping');
if ($m->getName() !== 'ping') {
    Log::fatal('ReflectionMethod($obj, ping) getName 失败');
}

$rc = new \ReflectionClass($o);
$rcName = $rc->getName();
if ($rcName !== ReflMethodThis_Target::class && !str_ends_with($rcName, 'ReflMethodThis_Target')) {
    Log::fatal('ReflectionClass($obj) 类名失败: ' . $rcName);
}

Log::info('ReflectionMethod/Class 接受对象测试通过');
