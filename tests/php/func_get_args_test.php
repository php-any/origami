<?php

namespace tests\php;

/**
 * \func_get_args() / \func_num_args()：词法器会把前导 \ 收进 IDENTIFIER，需还原为语言结构
 */

function FuncGetArgs_outer($a, $b, $c = null)
{
    return [
        'kw' => func_get_args(),
        'bs' => \func_get_args(),
        'num_kw' => func_num_args(),
        'num_bs' => \func_num_args(),
    ];
}

$r = FuncGetArgs_outer(1, 2, 3);
if ($r['kw'] !== [1, 2, 3] && !($r['kw'][0] === 1 && $r['kw'][1] === 2 && $r['kw'][2] === 3)) {
    // 允许稀疏/多余 null，但前三个必须正确
    if ($r['kw'][0] !== 1 || $r['kw'][1] !== 2 || $r['kw'][2] !== 3) {
        Log::fatal('func_get_args() 关键字失败: ' . json_encode($r['kw']));
    }
}
if ($r['bs'][0] !== 1 || $r['bs'][1] !== 2 || $r['bs'][2] !== 3) {
    Log::fatal('\\func_get_args() 失败: ' . json_encode($r['bs']));
}
if ($r['num_kw'] !== 3) {
    Log::fatal('func_num_args() 失败: ' . $r['num_kw']);
}
if ($r['num_bs'] !== 3) {
    Log::fatal('\\func_num_args() 失败: ' . $r['num_bs']);
}

// Carbon 风格：spread \func_get_args()
function FuncGetArgs_spread($a, $b)
{
    return [...\func_get_args()];
}
$s = FuncGetArgs_spread(10, 20);
if ($s[0] !== 10 || $s[1] !== 20) {
    Log::fatal('...\\func_get_args() 失败: ' . json_encode($s));
}

Log::info('func_get_args 反斜杠调用测试通过');
