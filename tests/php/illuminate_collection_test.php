<?php

namespace tests\php;

/**
 * Illuminate\Support\Collection 原生 Go 实现：常用变换与变更方法。
 */

if (!class_exists(\Illuminate\Support\Collection::class, false)) {
    Log::info('skip: Collection 原生类未注册');
    return;
}

$c = collect([1, 2, 3, 4, 5]);
$chunks = $c->chunk(2)->all();
if (count($chunks) !== 3) {
    Log::fatal('chunk 分组数错误: ' . count($chunks));
}
if ($chunks[0]->all() !== [1, 2]) {
    Log::fatal('chunk[0] 错误: ' . json_encode($chunks[0]->all()));
}

$slice = collect([10, 20, 30, 40])->slice(1, 2)->values()->all();
if ($slice !== [20, 30]) {
    Log::fatal('slice 错误: ' . json_encode($slice));
}

$take = collect(['a', 'b', 'c'])->take(2)->values()->all();
if ($take !== ['a', 'b']) {
    Log::fatal('take 错误: ' . json_encode($take));
}

$rev = collect([1, 2, 3])->reverse()->values()->all();
if ($rev !== [3, 2, 1]) {
    Log::fatal('reverse 错误: ' . json_encode($rev));
}

$flip = collect(['a' => 'x', 'b' => 'y'])->flip()->all();
if (($flip['x'] ?? null) !== 'a' || ($flip['y'] ?? null) !== 'b') {
    Log::fatal('flip 错误: ' . json_encode($flip));
}

$bag = collect(['keep' => 1, 'drop' => 2]);
$bag->forget('drop');
if ($bag->has('drop') || !$bag->has('keep')) {
    Log::fatal('forget 错误: ' . json_encode($bag->all()));
}

$prep = collect([2, 3])->prepend(1)->values()->all();
if ($prep !== [1, 2, 3]) {
    Log::fatal('prepend 错误: ' . json_encode($prep));
}

$pullBag = collect(['k' => 'v', 'other' => 9]);
$pulled = $pullBag->pull('k');
if ($pulled !== 'v' || $pullBag->has('k')) {
    Log::fatal('pull 错误 pulled=' . var_export($pulled, true) . ' all=' . json_encode($pullBag->all()));
}

$shiftBag = collect(['first', 'second']);
$head = $shiftBag->shift();
if ($head !== 'first' || $shiftBag->values()->all() !== ['second']) {
    Log::fatal('shift 错误 head=' . var_export($head, true));
}

$shuffled = collect([1, 2, 3, 4, 5])->shuffle()->sort()->values()->all();
if ($shuffled !== [1, 2, 3, 4, 5]) {
    Log::fatal('shuffle 内容丢失: ' . json_encode($shuffled));
}

$skip = collect([1, 2, 3, 4])->skip(2)->values()->all();
if ($skip !== [3, 4]) {
    Log::fatal('skip 错误: ' . json_encode($skip));
}

$desc = collect([3, 1, 2])->sortDesc()->values()->all();
if ($desc !== [3, 2, 1]) {
    Log::fatal('sortDesc 错误: ' . json_encode($desc));
}

$u = collect(['a' => 1])->union(['b' => 2, 'a' => 99])->all();
if (($u['a'] ?? null) !== 1 || ($u['b'] ?? null) !== 2) {
    Log::fatal('union 错误: ' . json_encode($u));
}

$zipped = collect([1, 2])->zip([3, 4], [5, 6])->all();
if (!isset($zipped[0]) || $zipped[0] !== [1, 3, 5]) {
    Log::fatal('zip 错误: ' . json_encode($zipped));
}

$added = collect(['x' => 1]);
$added->add('y', 2);
if (!$added->has('y') || $added->get('x') !== 1) {
    Log::fatal('add 错误: ' . json_encode($added->all()));
}

$collapsed = collect([[1, 2], [3]])->collapse()->values()->all();
if ($collapsed !== [1, 2, 3]) {
    Log::fatal('collapse 错误: ' . json_encode($collapsed));
}

if (!collect([1, '1', 2])->containsStrict(1) || !collect([1, '1', 2])->containsStrict('1')) {
    Log::fatal('containsStrict 应对 int/string 分别匹配');
}
if (collect([1, '1', 2])->containsStrict('2')) {
    Log::fatal('containsStrict 不应匹配 string 2');
}

if (collect([1, 2])->doesntContain(3) !== true || collect([1, 2])->doesntContain(1) !== false) {
    Log::fatal('doesntContain 错误');
}

$tx = collect([1, 2, 3]);
$tx->transform(fn ($v) => $v * 10);
if ($tx->all() !== [10, 20, 30]) {
    Log::fatal('transform 错误: ' . json_encode($tx->all()));
}

Log::info('illuminate_collection 测试通过');
