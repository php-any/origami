<?php

namespace tests\php;

/**
 * Component::extractConstructorParameters：ReflectionParameter::getName
 * 经 Collection map->getName 必须得到 ['view','data']。
 */

class CtorParams_Anon
{
    public function __construct($view, $data)
    {
    }
}

$ref = new \ReflectionClass(CtorParams_Anon::class);
$ctor = $ref->getConstructor();
$params = $ctor->getParameters();
$names = [];
foreach ($params as $p) {
    $names[] = $p->getName();
}
if ($names !== ['view', 'data']) {
    \Log::fatal('foreach getName 失败: '.var_export($names, true));
}

class CtorParams_Col
{
    public function __construct(protected array $items)
    {
    }

    public function map($cb)
    {
        $out = [];
        foreach ($this->items as $k => $v) {
            $out[$k] = $cb($v);
        }

        return new self($out);
    }

    public function all()
    {
        return $this->items;
    }

    public function __get($key)
    {
        return new CtorParams_Proxy($this, $key);
    }
}

class CtorParams_Proxy
{
    public function __construct(protected $collection, protected $method)
    {
    }

    public function __call($method, $parameters)
    {
        return $this->collection->{$this->method}(function ($value) use ($method, $parameters) {
            return $value->{$method}(...$parameters);
        });
    }
}

$mapped = (new CtorParams_Col($params))->map->getName()->all();
if ($mapped !== ['view', 'data']) {
    \Log::fatal('map->getName()->all 失败: '.var_export($mapped, true));
}

\Log::info('reflection_ctor_param_names 测试通过');
