<?php

use Illuminate\Support\Arr;
use Illuminate\Support\Collection;
use Illuminate\Support\HtmlString;
use Illuminate\Support\Str;
use Illuminate\Support\Stringable;

$attributes = ['rel' => 'stylesheet', 'href' => '/build/theme.css', 'nonce' => false, 'data-navigate-track' => 'reload', 'disabled' => true];
$parsed = (new Collection($attributes))
    ->reject(fn ($value, $key) => in_array($value, [false, null], true))
    ->flatMap(fn ($value, $key) => $value === true ? [$key] : [$key => $value])
    ->map(fn ($value, $key) => is_int($key) ? $value : $key.'="'.$value.'"')
    ->values()->all();
if ($parsed !== ['rel="stylesheet"', 'href="/build/theme.css"', 'data-navigate-track="reload"', 'disabled']) {
    throw new Exception('Vite attribute pipeline must preserve string keys during collapse');
}
$collapsed = Arr::collapse([['x' => 1, 8 => 'a', '' => 'empty'], ['x' => 2, 5 => 'b'], collect(['y' => 3]), 'skip']);
if ($collapsed !== ['x' => 2, 0 => 'a', '' => 'empty', 1 => 'b', 'y' => 3]) {
    throw new Exception('collapse must follow array_merge key semantics');
}
$html = new HtmlString('<link rel="stylesheet" href="/build/theme.css" />');
if (!Str::of($html)->contains('<link') || (string) new Stringable($html) !== $html->toHtml()) {
    throw new Exception('Stringable must call HtmlString::__toString');
}
class BrokenStyleString {
    public function __toString(): string { throw new RuntimeException('string conversion failed'); }
}
try {
    Str::of(new BrokenStyleString);
    throw new Exception('String conversion exceptions must propagate');
} catch (RuntimeException $expected) {}
echo "vite style attributes OK\n";
