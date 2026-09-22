<?php

namespace tests\origami;

/**
 * Illuminate\Config\Repository Go 实现冒烟：点号路径 get/set/has/push/prepend。
 */

use Illuminate\Config\Repository;

if (!class_exists(Repository::class, false)) {
    \Log::fatal('Illuminate\\Config\\Repository 应已由 std/laravel 预注册');
}

$repo = new Repository([
    'app' => [
        'name' => 'Origami',
        'debug' => true,
        'providers' => ['A'],
    ],
]);

if (!$repo->has('app.name')) {
    \Log::fatal('has(app.name) 应为 true');
}
if ($repo->get('app.name') !== 'Origami') {
    \Log::fatal('get(app.name) 应为 Origami，实际: ' . var_export($repo->get('app.name'), true));
}
if ($repo->get('app.missing', 'x') !== 'x') {
    \Log::fatal('get 默认值失败');
}

$repo->set('app.timezone', 'UTC');
if ($repo->get('app.timezone') !== 'UTC') {
    \Log::fatal('set 点号路径失败');
}

$repo->push('app.providers', 'B');
$providers = $repo->get('app.providers');
if (!is_array($providers) || count($providers) !== 2 || $providers[1] !== 'B') {
    \Log::fatal('push 失败: ' . json_encode($providers));
}

$repo->prepend('app.providers', 'Z');
$providers = $repo->get('app.providers');
if (!is_array($providers) || $providers[0] !== 'Z') {
    \Log::fatal('prepend 失败: ' . json_encode($providers));
}

if ($repo->string('app.name') !== 'Origami') {
    \Log::fatal('string() 失败');
}
if ($repo->boolean('app.debug') !== true) {
    \Log::fatal('boolean() 失败');
}

\Log::info('illuminate_config_repository 测试通过');
