<?php

namespace tests\php;

/**
 * WeakMap 对象键与 $wm[$obj][] = $x 写回，Livewire ComponentHookRegistry 依赖。
 */
class WeakMapHook_Comp
{
    public $id;
    public function __construct($id) { $this->id = $id; }
}

class WeakMapHook_Hook
{
    public $name;
    public function __construct($name) { $this->name = $name; }
}

$map = new \WeakMap();
$a = new WeakMapHook_Comp('a');
$b = new WeakMapHook_Comp('b');

if (isset($map[$a])) {
    Log::fatal('空 WeakMap isset 应为 false');
}

$map[$a] = [];
if (!isset($map[$a])) {
    Log::fatal('赋值后 isset 应为 true');
}

$map[$a][] = new WeakMapHook_Hook('h1');
$map[$a][] = new WeakMapHook_Hook('h2');

$hooks = $map[$a];
if (!is_array($hooks) || count($hooks) !== 2) {
    Log::fatal('WeakMap[$obj][] 追加失败: '.var_export($hooks, true));
}
if (($hooks[0]->name ?? null) !== 'h1' || ($hooks[1]->name ?? null) !== 'h2') {
    Log::fatal('WeakMap hook 内容错误');
}

// 同一对象再次查找
if (!isset($map[$a])) {
    Log::fatal('同一对象键丢失');
}
if (isset($map[$b])) {
    Log::fatal('不同对象不应命中');
}

$map[$b] = [new WeakMapHook_Hook('hb')];
if (count($map[$b]) !== 1 || $map[$b][0]->name !== 'hb') {
    Log::fatal('第二对象键失败');
}

Log::info('weakmap_hooks 测试通过');
