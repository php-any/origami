<?php
namespace tests\php;

class Hook {
    public static function render(string $name, $scopes = null) {
        return $name . ':' . json_encode($scopes);
    }
}

$got = Hook::render('x', scopes: ['a']);
if (!str_contains($got, 'a')) {
    Log::fatal('static named arg failed: '.$got);
}
Log::info('static named args ok: '.$got);
