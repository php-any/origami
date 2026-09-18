<?php
/**
 * 复现仪表盘 widgets 展开与 Collection::all。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Collection;
use Illuminate\Database\Eloquent\Collection as EloquentCollection;

echo "step1 support collection spread\n";
try {
    $c = new Collection(['a' => 1, 'b' => 2]);
    $m = [...$c];
    echo "support_ok count=".count($m)."\n";
} catch (Throwable $e) {
    echo "support EX: ".$e->getMessage()."\n";
}

echo "step2 eloquent collection spread\n";
try {
    $c = new EloquentCollection([]);
    $m = [...$c];
    echo "eloquent_ok count=".count($m)."\n";
} catch (Throwable $e) {
    echo "eloquent EX: ".$e->getMessage()."\n";
}

echo "step3 eloquent all()\n";
try {
    $c = new EloquentCollection([]);
    $all = $c->all();
    echo "all_ok type=".gettype($all)."\n";
} catch (Throwable $e) {
    echo "all EX: ".$e->getMessage()."\n";
}

echo "step4 widget default properties\n";
try {
    $props = \Filament\Widgets\StatsOverviewWidget::getDefaultProperties();
    echo "props_type=".gettype($props)."\n";
    if (is_object($props)) echo "props_class=".get_class($props)."\n";
    $m = [...$props];
    echo "props_spread_ok count=".count($m)."\n";
} catch (Throwable $e) {
    echo "props EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

echo "DONE\n";
