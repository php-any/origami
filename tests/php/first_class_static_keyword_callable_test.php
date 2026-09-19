<?php

namespace tests\php;

/**
 * PHP 8.1 static::method(...) 一等可调用：不得立刻执行，须绑定 late-static 类。
 * Filament Resource::configureTable 的 recordTitle(static::getRecordTitle(...))
 * 以及 EvaluatesClosures 按参数名注入 record 依赖此语义。
 */

class FirstClassStaticKw_Model
{
    public $name;

    public function __construct($name)
    {
        $this->name = $name;
    }
}

class FirstClassStaticKw_Resource
{
    public static function getRecordTitle($record): ?string
    {
        if ($record === null) {
            return static::label();
        }

        return $record->name;
    }

    public static function label(): string
    {
        return 'ResourceLabel';
    }

    public static function captureTitle()
    {
        return static::getRecordTitle(...);
    }
}

class FirstClassStaticKw_Child extends FirstClassStaticKw_Resource
{
    public static function label(): string
    {
        return 'ChildLabel';
    }
}

$fn = FirstClassStaticKw_Child::captureTitle();
if (!($fn instanceof \Closure)) {
    Log::fatal('static::getRecordTitle(...) 应返回 Closure，未立刻调用');
}

$params = (new \ReflectionFunction($fn))->getParameters();
if (count($params) !== 1) {
    Log::fatal('ReflectionFunction 参数个数应为 1，实际 '.count($params));
}
if ($params[0]->getName() !== 'record') {
    Log::fatal('参数名应为 record，实际 '.$params[0]->getName());
}

$named = ['record' => new FirstClassStaticKw_Model('Order#1')];
$deps = [];
foreach ($params as $parameter) {
    $n = $parameter->getName();
    if (array_key_exists($n, $named)) {
        $deps[] = $named[$n];
    } elseif ($parameter->isDefaultValueAvailable()) {
        $deps[] = $parameter->getDefaultValue();
    } elseif ($parameter->isOptional() || $parameter->allowsNull()) {
        $deps[] = null;
    } else {
        Log::fatal('无法解析参数: '.$n);
    }
}

$got = $fn(...$deps);
if ($got !== 'Order#1') {
    Log::fatal('evaluate named record 失败: '.var_export($got, true));
}

$got2 = $fn(null);
if ($got2 !== 'ChildLabel') {
    Log::fatal('late static binding 失败: '.var_export($got2, true));
}

Log::info('first_class_static_keyword_callable 测试通过');
