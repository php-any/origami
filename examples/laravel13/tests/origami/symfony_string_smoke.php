<?php

namespace tests\origami;

/**
 * vendoraccel：Symfony String 具体类必须在未 include vendor 类文件时可用。
 */

if (!class_exists(\Symfony\Component\String\UnicodeString::class, false)) {
    \Log::fatal('UnicodeString 应已由 std/symfony/string 预注册');
}
if (!class_exists(\Symfony\Component\String\ByteString::class, false)) {
    \Log::fatal('ByteString 应已由 std/symfony/string 预注册');
}
if (!class_exists(\Symfony\Component\String\AbstractString::class, false)) {
    \Log::fatal('AbstractString 应已由 std/symfony/string 预注册');
}
if (!class_exists(\Symfony\Component\String\CodePointString::class, false)) {
    \Log::fatal('CodePointString 应已由 std/symfony/string 预注册');
}
if (!class_exists(\Symfony\Component\String\AbstractUnicodeString::class, false)) {
    \Log::fatal('AbstractUnicodeString 应已由 std/symfony/string 预注册');
}

$u = new \Symfony\Component\String\UnicodeString('AbC');
$lower = $u->lower();
if ((string) $u !== 'AbC') {
    \Log::fatal('UnicodeString 不可变语义失败');
}
if ((string) $lower !== 'abc') {
    \Log::fatal('UnicodeString::lower 失败: ' . (string) $lower);
}
if ($u->length() !== 3) {
    \Log::fatal('UnicodeString::length 失败');
}

$b = new \Symfony\Component\String\ByteString('你好');
if ($b->length() !== 6) {
    \Log::fatal('ByteString::length 应为 6，实际 ' . $b->length());
}

$rand = \Symfony\Component\String\ByteString::fromRandom(12);
if ($rand->length() !== 12) {
    \Log::fatal('ByteString::fromRandom 失败');
}

\Log::info('symfony_string smoke 通过');
