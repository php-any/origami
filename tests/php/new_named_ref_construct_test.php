<?php

namespace tests\php;

/**
 * new Class(name: $var) 应对 by-ref 构造参数解开 NamedArgument，写回调用方变量。
 * 未解开时 Origami 会抛「引用参数只能传入变量」，Livewire/Filament 页面变成 RootTagMissing。
 */
class NewNamedRef_Box
{
    public int $n;

    public function __construct(&$items)
    {
        $items[] = 'x';
        $this->n = count($items);
    }
}

$a = [1];
$o = new NewNamedRef_Box(items: $a);
if ($o->n !== 2) {
    \Log::fatal('命名 by-ref 构造参数未生效 n='.var_export($o->n, true));
}
if ($a !== [1, 'x']) {
    \Log::fatal('命名 by-ref 构造参数未写回: '.var_export($a, true));
}

$b = [7];
$o2 = new NewNamedRef_Box($b);
if ($b !== [7, 'x']) {
    \Log::fatal('位置 by-ref 构造参数未写回: '.var_export($b, true));
}

class NewNamedRef_Prop
{
    public array $data = ['k' => 1];
}

$p = new NewNamedRef_Prop();
new NewNamedRef_Box(items: $p->data);
if (!in_array('x', $p->data, true)) {
    \Log::fatal('属性作命名 by-ref 构造实参未写回: '.var_export($p->data, true));
}

\Log::info('new 命名引用构造参数测试通过');
