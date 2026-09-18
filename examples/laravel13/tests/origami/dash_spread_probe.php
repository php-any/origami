<?php
/**
 * 逐步检查 Dashboard widget 数据展开。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;
use Filament\Widgets\WidgetConfiguration;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

$page = app(\App\Filament\Pages\Dashboard::class);
$widgets = $page->getWidgets();
echo "widgets_count=".count($widgets)."\n";

$wd = $page->getWidgetData();
echo "widgetData_type=".gettype($wd)." is_array=".(is_array($wd)?'y':'n')."\n";

foreach ($widgets as $i => $widget) {
    echo "w$i class=".(is_string($widget)?$widget:get_class($widget))."\n";
    if ($widget instanceof WidgetConfiguration) {
        $dp = $widget->widget::getDefaultProperties();
        $props = $widget->getProperties();
        echo "  cfg default=".gettype($dp)." props=".gettype($props)."\n";
    } else {
        $dp = $widget::getDefaultProperties();
        echo "  default=".gettype($dp)." count=".(is_array($dp)?count($dp):-1)."\n";
        if ($dp === null) echo "  DEFAULT IS NULL\n";
    }
}

// 手动拼装与 Page::getWidgetsSchemaComponents 相同的数组展开
try {
    foreach (array_values($widgets) as $widgetKey => $widget) {
        $widgetClass = is_string($widget) ? $widget : $widget->widget;
        $arr = [
            ...$page->getWidgetData(),
            ...[],
            ...(($widget instanceof WidgetConfiguration) ? [
                ...$widget->widget::getDefaultProperties(),
                ...$widget->getProperties(),
            ] : $widget::getDefaultProperties()),
        ];
        echo "built$widgetKey type=".gettype($arr)." count=".count($arr)."\n";
    }
    echo "manual_ok\n";
} catch (Throwable $e) {
    echo "manual EX: ".$e->getMessage()."\n";
}

try {
    $c = $page->getWidgetsSchemaComponents($widgets);
    echo "schema_ok count=".count($c)."\n";
} catch (Throwable $e) {
    echo "schema EX: ".$e->getMessage()."\n";
}

echo "DONE\n";
