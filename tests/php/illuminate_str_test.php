<?php

namespace tests\php;

/**
 * Illuminate\Support\Str / HtmlString 原生 Go 实现冒烟。
 */

use Illuminate\Support\Str;
use Illuminate\Support\HtmlString;

if (!class_exists(Str::class, false)) {
    Log::info('skip: Str 原生类未注册');
    return;
}

if (Str::camel('foo_bar') !== 'fooBar') {
    Log::fatal('camel: ' . Str::camel('foo_bar'));
}
if (Str::snake('fooBar') !== 'foo_bar') {
    Log::fatal('snake: ' . Str::snake('fooBar'));
}
if (Str::studly('foo-bar') !== 'FooBar') {
    Log::fatal('studly: ' . Str::studly('foo-bar'));
}
if (!Str::contains('hello world', 'world')) {
    Log::fatal('contains failed');
}
if (!Str::is('foo*', 'foobar')) {
    Log::fatal('is wildcard failed');
}
if (Str::kebab('FooBar') !== 'foo-bar') {
    Log::fatal('kebab: ' . Str::kebab('FooBar'));
}
if (Str::after('foo-bar-baz', '-bar-') !== 'baz') {
    Log::fatal('after: ' . Str::after('foo-bar-baz', '-bar-'));
}
if (Str::before('foo-bar-baz', '-bar-') !== 'foo') {
    Log::fatal('before: ' . Str::before('foo-bar-baz', '-bar-'));
}
if (Str::limit('abcdef', 3) !== 'abc...') {
    Log::fatal('limit: ' . Str::limit('abcdef', 3));
}
if (Str::slug('Hello World!') !== 'hello-world') {
    Log::fatal('slug: ' . Str::slug('Hello World!'));
}

$uuid = (string) Str::uuid();
if (!preg_match('/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i', $uuid)) {
    Log::fatal('uuid format: ' . $uuid);
}

$rand = Str::random(32);
if (strlen($rand) !== 32) {
    Log::fatal('random length: ' . strlen($rand));
}

$html = new HtmlString('<em>hi</em>');
if ($html->toHtml() !== '<em>hi</em>') {
    Log::fatal('HtmlString toHtml: ' . $html->toHtml());
}
if ((string) $html !== '<em>hi</em>') {
    Log::fatal('HtmlString __toString');
}

Log::info('illuminate_str_test ok');
