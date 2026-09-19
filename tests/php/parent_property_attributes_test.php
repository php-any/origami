<?php

namespace tests\php;

/**
 * 属性声明在父类、data()/withAttributes 写 $this->attributes，子类实例上必须能读到。
 * 对齐 Illuminate\View\Component vs AnonymousComponent。
 */

class ParentProp_Bag
{
    public function has($key)
    {
        return false;
    }

    public function setAttributes(array $a)
    {
    }

    public function getAttributes()
    {
        return [];
    }
}

class ParentProp_Component
{
    public $attributes;
    public $componentName;

    public function withName($name)
    {
        $this->componentName = $name;

        return $this;
    }

    public function withAttributes(array $attributes)
    {
        $this->attributes = $this->attributes ?: new ParentProp_Bag();
        $this->attributes->setAttributes($attributes);

        return $this;
    }

    public function shouldRender()
    {
        return true;
    }

    public function newAttributeBag(array $attributes = [])
    {
        return new ParentProp_Bag();
    }
}

class ParentProp_Anonymous extends ParentProp_Component
{
    protected $view;
    protected $data = [];

    public function __construct($view, $data)
    {
        $this->view = $view;
        $this->data = $data;
    }

    public static function resolve($data)
    {
        $parameters = ['view', 'data'];
        if (empty(array_diff($parameters, array_keys($data)))) {
            return new static(...array_intersect_key($data, array_flip($parameters)));
        }

        return new static('', []);
    }

    public function render()
    {
        return $this->view;
    }

    public function resolveView()
    {
        return $this->render();
    }

    public function data()
    {
        $this->attributes = $this->attributes ?: $this->newAttributeBag();

        return array_merge($this->data, ['attributes' => $this->attributes]);
    }
}

$c = ParentProp_Anonymous::resolve(['view' => 'filament-panels::components.page.index', 'data' => []]);
$c->withName('filament-panels::page');
if (!$c->shouldRender()) {
    \Log::fatal('shouldRender 应为 true');
}
$c->data();
$c->withAttributes([]);
if ($c->attributes === null) {
    \Log::fatal('父类 public $attributes 写入后仍为 null');
}
$ok = $c->attributes->has('wire:key');
if ($ok !== false) {
    \Log::fatal('has 失败: '.var_export($ok, true));
}

\Log::info('parent_property_attributes 测试通过');
