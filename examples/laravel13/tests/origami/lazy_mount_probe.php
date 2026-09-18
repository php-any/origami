<?php
/**
 * 在 Dashboard 嵌入路径上复现 lazy placeholder / mountParamsContainer。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

echo "mount_stats_with_lazy\n";
try {
    $html = app('livewire')->mount(\App\Filament\Widgets\StatsOverviewWidget::class, ['lazy' => true]);
    echo "lazy_stats_len=".strlen((string)$html)."\n";
    echo "has_intersect=".(str_contains((string)$html,'__lazyLoad')||str_contains((string)$html,'x-intersect')?'y':'n')."\n";
} catch (Throwable $e) {
    echo "lazy_stats EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

echo "mount_stats_no_lazy\n";
try {
    $html = app('livewire')->mount(\App\Filament\Widgets\StatsOverviewWidget::class);
    echo "stats_len=".strlen((string)$html)."\n";
} catch (Throwable $e) {
    echo "stats EX: ".$e->getMessage()."\n";
}

echo "DONE\n";
