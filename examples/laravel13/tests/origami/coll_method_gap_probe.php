<?php
/**
 * 列出 Illuminate\Support\Collection 的方法清单，用于对照 origami 桩类是否缺方法。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$names = [];
try {
    // 只看公开方法：protected/private 的内部辅助（sortByMany、valueRetriever 等）
    // Go 桩本来就不建模，混在一起会掩盖真正的公开 API 缺口。
    foreach ((new ReflectionClass(Illuminate\Support\Collection::class))->getMethods() as $m) {
        if (! $m->isPublic()) {
            continue;
        }
        $names[] = ($m->isStatic() ? 'static:' : '').$m->getName();
    }
} catch (Throwable $e) {
    echo 'EX '.$e->getMessage()."\n";
    exit(1);
}
sort($names);
echo implode("\n", $names), "\n";
echo 'count=', count($names), "\n";
