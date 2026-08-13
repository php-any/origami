<?php
namespace tests\php;

function byref_set_key(&$array, $key, $value) {
    $array[$key] = $value;
}

function byref_nested_set(&$array, $key, $value) {
    if (!isset($array['app']) || !is_array($array['app'])) {
        $array['app'] = [];
    }
    $ref = &$array['app'];
    $ref[$key] = $value;
}

$a = ['x' => 1];
byref_set_key($a, 'y', 2);
if (($a['y'] ?? null) !== 2) {
    Log::fatal('byref new key 失败: '.var_export($a, true));
}

$b = ['app' => ['name' => 'L']];
byref_nested_set($b, 'providers', ['A']);
if (($b['app']['providers'] ?? null) !== ['A']) {
    Log::fatal('byref nested 失败: '.var_export($b, true));
}

class ByRef_Prop {
    public $items = ['x' => 1];
    public function setY() {
        byref_set_key($this->items, 'y', 2);
    }
}
$o = new ByRef_Prop();
$o->setY();
if (($o->items['y'] ?? null) !== 2) {
    Log::fatal('byref property new key 失败: '.var_export($o->items, true));
}

Log::info('byref_array_new_key 测试通过');
