<?php

namespace tests\php;

/**
 * 对齐 Illuminate Collection：(new Collection($params))->map->getName()->all()
 * 依赖 __get 返回 HigherOrder 代理，再 __call 调到 map。
 */

class HOC_Param
{
    public function __construct(private string $name)
    {
    }

    public function getName()
    {
        return $this->name;
    }
}

class HOC_Proxy
{
    public function __construct(protected $collection, protected $method)
    {
    }

    public function __call($method, $parameters)
    {
        return $this->collection->{$this->method}(function ($value) use ($method, $parameters) {
            return is_string($value)
                ? $value::{$method}(...$parameters)
                : $value->{$method}(...$parameters);
        });
    }
}

class HOC_Col
{
    public function __construct(protected array $items)
    {
    }

    public function map($callback)
    {
        $out = [];
        foreach ($this->items as $k => $v) {
            $out[$k] = $callback($v);
        }

        return new self($out);
    }

    public function all()
    {
        return $this->items;
    }

    public function __get($key)
    {
        return new HOC_Proxy($this, $key);
    }
}

$c = new HOC_Col([new HOC_Param('view'), new HOC_Param('data')]);
$names = $c->map->getName()->all();
if ($names !== ['view', 'data']) {
    \Log::fatal('map->getName()->all() 失败: '.var_export($names, true));
}

$data = ['view' => 'filament-tables::components.search-field', 'data' => ['debounce' => 500]];
$parameters = $names;
$intersect = array_intersect_key($data, array_flip($parameters));
if (!isset($intersect['view']) || $intersect['view'] !== $data['view']) {
    \Log::fatal('array_intersect_key+flip 失败: '.var_export($intersect, true));
}

\Log::info('higher_order_collection_proxy 测试通过');
