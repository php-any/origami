<?php

namespace tests\php;

/**
 * array_map 单关联数组必须保留字符串键（对齐 PHP；ComponentAttributeBag::merge 依赖）。
 */

$src = ['class' => 'x', 'id' => 'y'];
$mapped = array_map(function ($v) {
    return strtoupper($v);
}, $src);

if (!is_array($mapped) || !array_key_exists('class', $mapped) || !array_key_exists('id', $mapped)) {
    Log::fatal('array_map 未保留关联键: ' . json_encode($mapped));
}
if ($mapped['class'] !== 'X' || $mapped['id'] !== 'Y') {
    Log::fatal('array_map 关联值错误: ' . json_encode($mapped));
}

// 多数组仍使用顺序整数键
$a = ['a' => 1, 'b' => 2];
$b = ['a' => 10, 'b' => 20];
$mapped2 = array_map(function ($x, $y) {
    return $x + $y;
}, $a, $b);
if (!isset($mapped2[0]) || $mapped2[0] !== 11 || $mapped2[1] !== 22) {
    Log::fatal('array_map 多数组应重索引: ' . json_encode($mapped2));
}

// 单数组整数键值必须保留（对齐 PHP；槽位身份 CopyZValKeepName）
$sparse = [1 => 'a', 3 => 'b'];
$mapped3 = array_map('strtoupper', $sparse);
if ($mapped3[1] !== 'A' || $mapped3[3] !== 'B') {
    Log::fatal('array_map 整数键值错误: ' . json_encode($mapped3));
}

// ComponentAttributeBag::merge：array_map 后 array_merge 不得把 class 冲成 0=
$defaults = ['aria-label' => '打开', 'type' => 'button', 'class' => 'fi-icon-btn'];
$mappedDefaults = array_map(function ($v) {
    return $v;
}, $defaults);
$bag = array_merge($mappedDefaults, ['class' => '']);
$html = '';
foreach ($bag as $key => $value) {
    $html .= ' '.$key.'="'.$value.'"';
}
if (str_contains($html, ' 0="') || !str_contains($html, ' aria-label="打开"') || !str_contains($html, ' type="button"')) {
    Log::fatal('array_map+merge 属性键丢失: '.$html);
}

Log::info('array_map 保留关联键测试通过');
