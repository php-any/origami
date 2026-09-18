<?php

namespace tests\php;

/**
 * filter_var(FILTER_VALIDATE_URL)：相对路径必须失败，带 scheme 的 URL 通过。
 * Laravel UrlGenerator::asset / isValidUrl 依赖此行为拼接 asset 根路径。
 */

$rel = filter_var('js/filament/filament/app.js', FILTER_VALIDATE_URL);
if ($rel !== false) {
    Log::fatal('相对路径不应通过 FILTER_VALIDATE_URL, got: ' . var_export($rel, true));
}

$slashRel = filter_var('/js/x', FILTER_VALIDATE_URL);
if ($slashRel !== false) {
    Log::fatal('/js/x 不应通过 FILTER_VALIDATE_URL');
}

$hostOnly = filter_var('example.com', FILTER_VALIDATE_URL);
if ($hostOnly !== false) {
    Log::fatal('无 scheme 的 example.com 不应通过');
}

$emptyHttp = filter_var('http://', FILTER_VALIDATE_URL);
if ($emptyHttp !== false) {
    Log::fatal('http:// 无 host 不应通过');
}

$ok = filter_var('http://127.0.0.1:8000/js/x', FILTER_VALIDATE_URL);
if ($ok !== 'http://127.0.0.1:8000/js/x') {
    Log::fatal('完整 URL 应原样返回, got: ' . var_export($ok, true));
}

$short = filter_var('http://a', FILTER_VALIDATE_URL);
if ($short !== 'http://a') {
    Log::fatal('http://a 应通过, got: ' . var_export($short, true));
}

if (FILTER_VALIDATE_URL !== 273) {
    Log::fatal('FILTER_VALIDATE_URL 常量应为 273, got: ' . FILTER_VALIDATE_URL);
}

Log::info('filter_var FILTER_VALIDATE_URL 测试通过');
