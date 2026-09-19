<?php

namespace tests\php;

/**
 * 对齐 PHP array_merge_recursive：字符串键递归，整数键追加不覆盖。
 * Laravel RouteGroup 用它把 Filament 面板中间件与 authMiddleware 拼在一起。
 */

$merged = array_merge_recursive([], ['title' => ['required', 'min:1']]);
if (!isset($merged['title']) || $merged['title'][0] !== 'required') {
    Log::fatal('array_merge_recursive 丢失字符串键');
}

$mw = array_merge_recursive(
    ['middleware' => [
        'Illuminate\\Cookie\\Middleware\\EncryptCookies',
        'Illuminate\\Session\\Middleware\\StartSession',
    ]],
    ['middleware' => [
        'Filament\\Http\\Middleware\\Authenticate',
    ]]
);
if (($mw['middleware'][0] ?? null) !== 'Illuminate\\Cookie\\Middleware\\EncryptCookies') {
    Log::fatal('整数键 0 被覆盖，EncryptCookies 应仍在首位，实际=' . var_export($mw['middleware'][0] ?? null, true));
}
if (($mw['middleware'][1] ?? null) !== 'Illuminate\\Session\\Middleware\\StartSession') {
    Log::fatal('StartSession 应保留在整数键 1');
}
if (($mw['middleware'][2] ?? null) !== 'Filament\\Http\\Middleware\\Authenticate') {
    Log::fatal('Authenticate 应按整数键追加，实际=' . var_export($mw['middleware'][2] ?? null, true));
}
if (count($mw['middleware']) !== 3) {
    Log::fatal('middleware 应有 3 项，实际=' . count($mw['middleware']));
}

$str = array_merge_recursive(['a' => 'x'], ['a' => 'y']);
if (!is_array($str['a']) || $str['a'][0] !== 'x' || $str['a'][1] !== 'y') {
    Log::fatal('同名字符串键的标量应合并为数组 [x, y]');
}

$num = array_merge_recursive([0 => 'a', 1 => 'b'], [0 => 'c']);
if ($num[0] !== 'a' || $num[1] !== 'b' || $num[2] !== 'c') {
    Log::fatal('纯整数键应追加为 a,b,c，实际=' . var_export($num, true));
}

Log::info('array_merge_recursive 测试通过');
