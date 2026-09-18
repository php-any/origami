<?php

namespace tests\php;

/**
 * new Class(...['a'=>$x,'b'=>$y]) 命名参数展开须按形参名绑定。
 */
class NamedUnpackCtor_Box
{
    public $view;
    public $data;

    public function __construct($view, $data)
    {
        $this->view = $view;
        $this->data = $data;
    }
}

$payload = [
    'view' => 'filament-panels::components.page.simple',
    'data' => [],
];

$params = ['view', 'data'];
$args = array_intersect_key($payload, array_flip($params));
$obj = new NamedUnpackCtor_Box(...$args);

if ($obj->view !== 'filament-panels::components.page.simple') {
    Log::fatal('view 绑定失败: ' . var_export($obj->view, true));
}
if (!is_array($obj->data)) {
    Log::fatal('data 应为数组: ' . var_export($obj->data, true));
}

// 乱序命名键也应按名绑定
$obj2 = new NamedUnpackCtor_Box(...['data' => ['x' => 1], 'view' => 'v2']);
if ($obj2->view !== 'v2' || ($obj2->data['x'] ?? null) !== 1) {
    Log::fatal('乱序命名展开失败: view=' . var_export($obj2->view, true) . ' data=' . json_encode($obj2->data));
}

Log::info('named_unpack_ctor 测试通过');
