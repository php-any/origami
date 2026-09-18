<?php
/**
 * 逐步渲染 Dashboard / Widgets，定位展开与 Collection::all 递归。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;

$log = function ($msg) {
    echo $msg."\n";
    file_put_contents(storage_path('framework/admin-step.txt'), $msg."\n", FILE_APPEND);
};

@unlink(storage_path('framework/admin-step.txt'));
$admin = Admin::where('email', 'admin@example.com')->first();
Auth::guard('admin')->login($admin);
Filament::setCurrentPanel(Filament::getPanel('admin'));
$log('auth_ok');

try {
    $widgets = \App\Filament\Pages\Dashboard::getWidgets();
    $log('widgets_count='.(is_countable($widgets)?count($widgets):gettype($widgets)));
    if (is_array($widgets) || $widgets instanceof Traversable) {
        foreach ($widgets as $i => $w) {
            $cls = is_object($w) ? get_class($w) : (is_string($w) ? $w : gettype($w));
            $log("widget[$i]=$cls");
            try {
                if (is_string($w) && class_exists($w)) {
                    $props = $w::getDefaultProperties();
                    $log("  defaultProps=".json_encode($props));
                    $merged = [...$props];
                    $log("  spread_defaults_ok");
                } elseif ($w instanceof \Filament\Widgets\WidgetConfiguration) {
                    $merged = [...$w->widget::getDefaultProperties(), ...$w->getProperties()];
                    $log("  config_spread_ok count=".count($merged));
                }
            } catch (Throwable $e) {
                $log("  widget EX: ".$e->getMessage());
            }
        }
    }
} catch (Throwable $e) {
    $log('getWidgets EX: '.$e->getMessage());
}

$log('mount_stats');
try {
    $html = app('livewire')->mount(\App\Filament\Widgets\StatsOverviewWidget::class);
    $log('stats_len='.strlen((string)$html));
} catch (Throwable $e) {
    $log('stats EX: '.$e->getMessage());
    $log($e->getFile().':'.$e->getLine());
}

$log('mount_dashboard');
try {
    $html = app('livewire')->mount(\App\Filament\Pages\Dashboard::class);
    $log('dash_len='.strlen((string)$html));
} catch (Throwable $e) {
    $log('dash EX: '.$e->getMessage());
    $log($e->getFile().':'.$e->getLine());
    foreach (array_slice($e->getTrace(), 0, 10) as $i => $t) {
        $log("#$i ".($t['class']??'').($t['type']??'').($t['function']??'').' @ '.($t['file']??'').':'.($t['line']??''));
    }
}

$log('DONE');
