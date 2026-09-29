<?php

namespace tests\php;

/**
 * Collection::groupBy() 必须执行回调；Filament 的导航构建依赖这一语义。
 */

if (!class_exists(\Illuminate\Support\Collection::class, false)) {
    Log::info('skip: Collection 原生类未注册');
    return;
}

$iterator = collect(['first' => 1, 'second' => 2])->getIterator();
$iteratorArray = iterator_to_array($iterator);
if (! $iterator instanceof \Traversable || count($iteratorArray) !== 2 || $iteratorArray['first'] !== 1 || $iteratorArray['second'] !== 2) {
    Log::fatal('Collection::getIterator() 应返回保留键的 Traversable，实际=' . get_debug_type($iterator) . ' data=' . json_encode($iteratorArray));
}

$items = collect([
    'first' => ['name' => 'A', 'group' => 'x'],
    'second' => ['name' => 'B', 'group' => 'x'],
    'third' => ['name' => 'C', 'group' => 'y'],
]);

$groups = $items->groupBy(fn (array $item, string $key): string => $item['group']);
if ($groups->keys()->all() !== ['x', 'y']) {
    Log::fatal('groupBy 回调分组键错误: ' . json_encode($groups->keys()->all()));
}
if ($groups->get('x')->pluck('name')->all() !== ['A', 'B']) {
    Log::fatal('groupBy 回调分组内容错误: ' . json_encode($groups->get('x')->all()));
}

$preserved = $items->groupBy(fn (array $item): string => $item['group'], true);
if ($preserved->get('x')->keys()->all() !== ['first', 'second']) {
    Log::fatal('groupBy preserveKeys 错误: ' . json_encode($preserved->get('x')->keys()->all()));
}

$multi = collect([1, 2])->groupBy(fn (int $value): array => ['all', $value % 2 ? 'odd' : 'even']);
if ($multi->get('all')->all() !== [1, 2] || $multi->get('odd')->all() !== [1] || $multi->get('even')->all() !== [2]) {
    Log::fatal('groupBy 多分组键错误');
}

$ungrouped = collect(['dashboard', 'orders'])->groupBy(fn (): string => '');
if ($ungrouped->get('missing') !== null || $ungrouped->get('')->all() !== ['dashboard', 'orders']) {
    Log::fatal('groupBy 空字符串键错误');
}

Log::info('illuminate_collection_group_by 测试通过');
