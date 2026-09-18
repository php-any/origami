<?php
/**
 * 变量类名静态调用 + widgets blade 同等展开。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$widget = \App\Filament\Widgets\StatsOverviewWidget::class;
echo "widget=$widget\n";

try {
    $props = $widget::getDefaultProperties();
    echo "props=".json_encode($props)." type=".gettype($props)."\n";
} catch (Throwable $e) {
    echo "static EX: ".$e->getMessage()."\n";
}

$data = [];
try {
    $merged = [...(($widget instanceof \Filament\Widgets\WidgetConfiguration) ? [...$widget->widget::getDefaultProperties(), ...$widget->getProperties()] : $widget::getDefaultProperties()), ...$data];
    echo "merged=".json_encode($merged)."\n";
} catch (Throwable $e) {
    echo "merge EX: ".$e->getMessage()."\n";
}

// data 为 null 时
$data = null;
try {
    $merged = [...$widget::getDefaultProperties(), ...$data];
    echo "null_data_merged=".json_encode($merged)."\n";
} catch (Throwable $e) {
    echo "null_data EX: ".$e->getMessage()."\n";
}

echo "DONE\n";
