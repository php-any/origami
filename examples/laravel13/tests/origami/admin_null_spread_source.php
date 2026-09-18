<?php
/**
 * 定位 /admin 中 ...null 的来源。
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
echo "page=".get_class($page)."\n";

try {
    if (method_exists($page, 'mount')) {
        $page->mount();
    }
} catch (Throwable $e) {
    echo "mount EX: ".$e->getMessage()."\n";
}

$widgets = $page->getWidgets();
echo "widgets=".json_encode($widgets)."\n";

try {
    $components = $page->getWidgetsSchemaComponents($widgets);
    echo "schema_count=".(is_countable($components)?count($components):gettype($components))."\n";
} catch (Throwable $e) {
    echo "schema EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

// 直接渲染废弃的 widgets 视图
try {
    $html = view('filament-widgets::components.widgets', [
        'widgets' => $widgets,
        'data' => [],
        'columns' => 2,
    ])->render();
    echo "view_len=".strlen($html)."\n";
} catch (Throwable $e) {
    echo "view EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

// data 故意为 null
try {
    $html = view('filament-widgets::components.widgets', [
        'widgets' => [\App\Filament\Widgets\StatsOverviewWidget::class],
        'data' => null,
        'columns' => 2,
    ])->render();
    echo "null_data_len=".strlen($html)."\n";
} catch (Throwable $e) {
    echo "null_data EX: ".$e->getMessage()."\n";
}

echo "DONE\n";
