<?php
/**
 * Probe: $view->missingKey ?? default must not warn; and who is $this for notifications.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

// layoutConfig ?? pattern
$view = Illuminate\Support\Facades\View::make('origami-capture-test', ['x' => 1]);
$before = error_get_last();
$layoutConfig = $view->layoutConfig ?? 'DEFAULT';
$after = error_get_last();
echo "coalesce_result=".var_export($layoutConfig, true)."\n";
if ($after !== $before && $after && str_contains($after['message'] ?? '', 'layoutConfig')) {
    echo "WARN_ON_COALESCE=yes msg=".$after['message']."\n";
} else {
    echo "WARN_ON_COALESCE=no\n";
}

// Direct __get should warn in PHP 8
$before2 = error_get_last();
try {
    $_ = $view->layoutConfig;
    echo "direct_get=".var_export($_, true)."\n";
} catch (Throwable $e) {
    echo "direct_EX: ".$e->getMessage()."\n";
}
$after2 = error_get_last();
if ($after2 && str_contains($after2['message'] ?? '', 'layoutConfig')) {
    echo "WARN_ON_DIRECT=yes\n";
} else {
    echo "WARN_ON_DIRECT=no\n";
}

echo "DONE\n";
