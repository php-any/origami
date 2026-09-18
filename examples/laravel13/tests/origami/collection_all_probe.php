<?php
/**
 * collect()->map()->all() 与 getWidgetData 行为。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

$page = app(\App\Filament\Pages\Dashboard::class);
echo "widgetData=";
try {
    $wd = $page->getWidgetData();
    echo gettype($wd)." ".json_encode($wd)."\n";
} catch (Throwable $e) {
    echo "EX ".$e->getMessage()."\n";
}

$widgets = ['App\\Filament\\Widgets\\StatsOverviewWidget'];
echo "collect_all\n";
try {
    $result = collect($widgets)
        ->values()
        ->map(fn (string $w, int $k) => $w)
        ->all();
    echo "result_type=".gettype($result)." countable=".(is_countable($result)?'y':'n')."\n";
    echo "result=".json_encode($result)."\n";
} catch (Throwable $e) {
    echo "collect EX: ".$e->getMessage()."\n";
}

echo "eloquent_all\n";
try {
    $c = new Illuminate\Database\Eloquent\Collection([]);
    $a = $c->all();
    echo "eloquent_all_type=".gettype($a)."\n";
} catch (Throwable $e) {
    echo "eloquent EX: ".$e->getMessage()."\n";
}

echo "support_all\n";
try {
    $c = collect([1,2,3]);
    echo "class=".get_class($c)."\n";
    $a = $c->all();
    echo "support_all_type=".gettype($a)." json=".json_encode($a)."\n";
    echo "items_prop=";
    // invade
    $ref = new ReflectionClass($c);
    $p = $ref->getProperty('items');
    $p->setAccessible(true);
    $items = $p->getValue($c);
    echo gettype($items)." ".json_encode($items)."\n";
} catch (Throwable $e) {
    echo "support EX: ".$e->getMessage()."\n";
}

echo "schema_again\n";
try {
    $components = $page->getWidgetsSchemaComponents($page->getWidgets());
    echo "type=".gettype($components)." class=".(is_object($components)?get_class($components):'-')."\n";
    if (is_object($components) && method_exists($components, 'all')) {
        echo "calling all on result...\n";
        $all = $components->all();
        echo "inner_all_type=".gettype($all)."\n";
    }
} catch (Throwable $e) {
    echo "schema EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

echo "DONE\n";
