<?php

namespace tests\php;

/**
 * AnonymousComponent::resolve + data()/withAttributes 后 attributes 不得为 null，
 * 否则 Livewire processComponentKey 会 $component->attributes->has() 炸在 null 上。
 */

class AttrPersist_Bag
{
    public $attrs = [];

    public function has($key)
    {
        return array_key_exists($key, $this->attrs);
    }

    public function setAttributes(array $attributes)
    {
        $this->attrs = $attributes;
    }

    public function getAttributes()
    {
        return $this->attrs;
    }
}

class AttrPersist_Comp
{
    public $attributes;
    public $view;
    public $data = [];

    public function __construct($view, $data)
    {
        $this->view = $view;
        $this->data = $data;
    }

    public static function resolve($data)
    {
        $parameters = ['view', 'data'];
        $dataKeys = array_keys($data);
        if (empty(array_diff($parameters, $dataKeys))) {
            return new static(...array_intersect_key($data, array_flip($parameters)));
        }

        return new static($data['view'] ?? '', $data['data'] ?? []);
    }

    public function data()
    {
        $this->attributes = $this->attributes ?: new AttrPersist_Bag();

        return array_merge(
            ($this->data['attributes'] ?? null)?->getAttributes() ?: [],
            $this->attributes->getAttributes(),
            $this->data,
            ['attributes' => $this->attributes]
        );
    }

    public function withAttributes(array $attributes)
    {
        $this->attributes = $this->attributes ?: new AttrPersist_Bag();
        $this->attributes->setAttributes($attributes);

        return $this;
    }

    public function resolveView()
    {
        return $this->view;
    }
}

$c = AttrPersist_Comp::resolve([
    'view' => 'filament-panels::components.page.index',
    'data' => [],
]);
if ($c->view !== 'filament-panels::components.page.index') {
    \Log::fatal('resolve 未绑定 view: '.var_export($c->view, true));
}

$view = $c->resolveView();
$data = $c->data();
if ($c->attributes === null) {
    \Log::fatal('data() 后 attributes 仍为 null');
}
$c->withAttributes([]);
if ($c->attributes === null) {
    \Log::fatal('withAttributes 后 attributes 仍为 null');
}
$has = $c->attributes->has('wire:key');
if ($has !== false) {
    \Log::fatal('has(wire:key) 应为 false, 实际: '.var_export($has, true));
}
if ($view !== 'filament-panels::components.page.index') {
    \Log::fatal('resolveView 错误: '.var_export($view, true));
}

// 缺构造参数时不应把后续 attributes 写丢
try {
    $c2 = new AttrPersist_Comp();
} catch (\Throwable $e) {
    $c2 = 'threw:'.$e->getMessage();
}
if (is_object($c2)) {
    $c2->withAttributes([]);
    if ($c2->attributes === null) {
        \Log::fatal('缺参 new 后 withAttributes 仍 null');
    }
}

\Log::info('component_attributes_persist 测试通过');
