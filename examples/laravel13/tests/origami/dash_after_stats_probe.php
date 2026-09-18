<?php
/**
 * 验证：先 mount widgets 再 mount Dashboard 是否污染 __mountParamsContainer。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

$check = function ($label) {
    try {
        $hook = new \Livewire\Features\SupportLazyLoading\SupportLazyLoading();
        $ref = new ReflectionClass($hook);
        $m = $ref->getMethod('registerContainerComponent');
        $m->setAccessible(true);
        $m->invoke($hook);
        $c = app('livewire')->new('__mountParamsContainer');
        echo "$label register+new_ok\n";
    } catch (Throwable $e) {
        echo "$label EX: ".$e->getMessage()."\n";
    }
};

$check('before');

echo "mount_stats\n";
try {
    $html = app('livewire')->mount(\App\Filament\Widgets\StatsOverviewWidget::class);
    echo "stats_len=".strlen((string)$html)."\n";
} catch (Throwable $e) {
    echo "stats EX: ".$e->getMessage()."\n";
}
$check('after_stats');

echo "mount_dash\n";
try {
    $html = app('livewire')->mount(\App\Filament\Pages\Dashboard::class);
    echo "dash_len=".strlen((string)$html)."\n";
    echo "dash_has_lazy=".(str_contains((string)$html,'__lazyLoad')?'y':'n')."\n";
    echo "dash_has_仪表盘=".(str_contains((string)$html,'仪表盘')||str_contains((string)$html,'Dashboard')?'y':'n')."\n";
    echo "head=".substr(preg_replace('/\s+/',' ',strip_tags((string)$html)),0,180)."\n";
} catch (Throwable $e) {
    echo "dash EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}
$check('after_dash');

echo "DONE\n";
