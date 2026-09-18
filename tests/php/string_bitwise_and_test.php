<?php

namespace tests\php;

/**
 * 两端都是字符串时，& | ^ 按字节运算并返回字符串。
 * Symfony polyfill-mbstring 的 mb_encode_numericentity 依赖 $s[$i] & "\xF0"。
 */

$and = "\xE4" & "\xF0";
if ($and !== "\xE0") {
    \Log::fatal('字符串按位与错误: ' . bin2hex($and));
}

$ch = chr(228);
if ($ch !== "\xE4" || strlen($ch) !== 1) {
    \Log::fatal('chr(228) 应为单字节 E4, 实际: ' . bin2hex($ch));
}

$s = "A\xE4\xB8\xAD";
if ($s[1] !== "\xE4" || strlen($s[1]) !== 1) {
    \Log::fatal('$s[1] 应为单字节 E4, 实际: ' . bin2hex($s[1]));
}

$or = "\x0F" | "\xF0";
if ($or !== "\xFF") {
    \Log::fatal('字符串按位或错误: ' . bin2hex($or));
}

$xor = "\xFF" ^ "\x0F";
if ($xor !== "\xF0") {
    \Log::fatal('字符串按位异或错误: ' . bin2hex($xor));
}

if ((5 & 3) !== 1) {
    \Log::fatal('整数按位与错误');
}

$ulenMask = ["\xC0" => 2, "\xD0" => 2, "\xE0" => 3, "\xF0" => 4];
$s = "A\xE4\xB8\xAD";
$i = 0;
$len = strlen($s);
$got = [];
$steps = 0;
while ($i < $len) {
    $ulen = $s[$i] < "\x80" ? 1 : $ulenMask[$s[$i] & "\xF0"];
    if ($ulen < 1) {
        \Log::fatal('ulen 为 0，将死循环: i=' . $i . ' byte=' . bin2hex($s[$i]));
    }
    $got[] = $ulen;
    $i += $ulen;
    $steps++;
    if ($steps > 20) {
        \Log::fatal('循环步数过多');
    }
}
if ($got !== [1, 3]) {
    \Log::fatal('UTF-8 宽度错误: ' . var_export($got, true));
}

\Log::info('string_bitwise_and 测试通过');
