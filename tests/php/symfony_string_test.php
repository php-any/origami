<?php

namespace tests\php;

/**
 * Symfony String 热路径：UnicodeString / ByteString 不可变操作、长度单位、ascii/width。
 * laravel13 vendoraccel 原生类与 vendor PHP 实现共用本脚本。
 */

if (!class_exists(\Symfony\Component\String\UnicodeString::class, false)) {
    \Log::info('skip: Symfony String 原生类未注册（请用 examples/laravel13 运行本脚本）');
    return;
}

use Symfony\Component\String\ByteString;
use Symfony\Component\String\CodePointString;
use Symfony\Component\String\UnicodeString;

$src = ' Hello Symfony ';
$u = new UnicodeString($src);
if ((string) $u !== $src) {
    \Log::fatal('UnicodeString __toString 失败: [' . (string) $u . ']');
}
$lowered = $u->lower();
if ((string) $u !== $src) {
    \Log::fatal('UnicodeString 应不可变，lower 后原对象变成: [' . (string) $u . ']');
}
if ((string) $lowered !== ' hello symfony ') {
    \Log::fatal('UnicodeString::lower 失败: [' . (string) $lowered . ']');
}

$han = new UnicodeString('你好');
if ($han->length() !== 2) {
    \Log::fatal('UnicodeString::length 应为字素数 2，实际 ' . $han->length());
}
$bytes = new ByteString('你好');
if ($bytes->length() !== 6) {
    \Log::fatal('ByteString::length 应为字节数 6，实际 ' . $bytes->length());
}
$cp = new CodePointString('你好');
if ($cp->length() !== 2) {
    \Log::fatal('CodePointString::length 应为码点数 2，实际 ' . $cp->length());
}

$trim = (new UnicodeString('  ok  '))->trim();
if ((string) $trim !== 'ok') {
    \Log::fatal('trim 失败: [' . (string) $trim . ']');
}

$joined = (new UnicodeString('-'))->join(['a', 'b', 'c']);
if ((string) $joined !== 'a-b-c') {
    \Log::fatal('join 失败: [' . (string) $joined . ']');
}

$sliced = (new UnicodeString('Symfony'))->slice(3, 3);
if ((string) $sliced !== 'fon') {
    \Log::fatal('slice 失败: [' . (string) $sliced . ']');
}

$repl = (new UnicodeString('foo-bar-foo'))->replace('foo', 'baz');
if ((string) $repl !== 'baz-bar-baz') {
    \Log::fatal('replace 失败: [' . (string) $repl . ']');
}

$sw = new UnicodeString('Hello World');
if (!$sw->startsWith('Hello') || !$sw->endsWith('World') || !$sw->containsAny('lo W')) {
    \Log::fatal('startsWith/endsWith/containsAny 失败');
}
if ($sw->indexOf('World') !== 6) {
    \Log::fatal('indexOf 失败: ' . var_export($sw->indexOf('World'), true));
}

$camel = (new UnicodeString('hello-world_test'))->camel();
if ((string) $camel !== 'helloWorldTest') {
    \Log::fatal('camel 失败: [' . (string) $camel . ']');
}
$snake = (new UnicodeString('helloWorld'))->snake();
if ((string) $snake !== 'hello_world') {
    \Log::fatal('snake 失败: [' . (string) $snake . ']');
}

$ascii = (new UnicodeString('Crème Brûlée'))->ascii();
if (!str_contains((string) $ascii, 'Creme') || !str_contains((string) $ascii, 'Brulee')) {
    \Log::fatal('ascii 失败: [' . (string) $ascii . ']');
}

$wrapped = (new UnicodeString('hello world foo'))->wordwrap(5);
if (!str_contains((string) $wrapped, "\n")) {
    \Log::fatal('wordwrap 应插入换行，实际: [' . (string) $wrapped . ']');
}

$w = (new UnicodeString('你好'))->width();
if ($w < 4) {
    \Log::fatal('width 过小: ' . $w);
}

$after = (new UnicodeString('foo/bar/baz'))->after('/');
if ((string) $after !== 'bar/baz') {
    \Log::fatal('after 失败: [' . (string) $after . ']');
}
$beforeLast = (new UnicodeString('foo/bar/baz'))->beforeLast('/');
if ((string) $beforeLast !== 'foo/bar') {
    \Log::fatal('beforeLast 失败: [' . (string) $beforeLast . ']');
}

$ic = (new UnicodeString('Hello'))->ignoreCase()->startsWith('hello');
if ($ic !== true) {
    \Log::fatal('ignoreCase()->startsWith 失败');
}

$rand = ByteString::fromRandom(8);
if ($rand->length() !== 8) {
    \Log::fatal('fromRandom 长度失败: ' . $rand->length());
}

$json = json_encode(new UnicodeString('hi'));
if ($json !== '"hi"') {
    \Log::fatal('jsonSerialize 失败: ' . $json);
}

if (!($han instanceof \Symfony\Component\String\AbstractString)) {
    \Log::fatal('UnicodeString 应 instanceof AbstractString，实际 ' . get_class($han));
}
if (!($han instanceof \Symfony\Component\String\AbstractUnicodeString)) {
    \Log::fatal('UnicodeString 应 instanceof AbstractUnicodeString，实际 ' . get_class($han));
}

$fromCp = UnicodeString::fromCodePoints(0x4F60, 0x597D);
if ((string) $fromCp !== '你好') {
    \Log::fatal('fromCodePoints 失败: [' . (string) $fromCp . ']');
}

$bag = UnicodeString::wrap(['a' => 'hello', 'n' => ['x' => 'y']]);
if (!($bag['a'] instanceof UnicodeString) || (string) $bag['a'] !== 'hello') {
    \Log::fatal('wrap 失败');
}
$plain = UnicodeString::unwrap($bag);
if ($plain['a'] !== 'hello' || $plain['n']['x'] !== 'y') {
    \Log::fatal('unwrap 失败: ' . var_export($plain, true));
}

if (UnicodeString::NFC !== 4) {
    \Log::fatal('UnicodeString::NFC 应为 4');
}

\Log::info('symfony_string 测试通过');
