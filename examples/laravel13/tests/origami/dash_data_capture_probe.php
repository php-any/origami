<?php
/**
 * 复现 Page::getWidgetsSchemaComponents 内 ...$data 捕获。
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

// 模拟 Livewire::make 的第二参闭包被延迟调用
$widgets = $page->getWidgets();
$components = $page->getWidgetsSchemaComponents($widgets);
echo "components=".count($components)."\n";

foreach ($components as $i => $comp) {
    echo "c$i=".get_class($comp)."\n";
    // Livewire 组件上的 data / params
    try {
        if (method_exists($comp, 'getChildSchema')) {
            // Schema Livewire component
        }
        // Filament Schemas Components Livewire
        $ref = new ReflectionObject($comp);
        foreach (['data', 'properties', 'componentProperties', 'params'] as $prop) {
            if ($ref->hasProperty($prop)) {
                $p = $ref->getProperty($prop);
                $p->setAccessible(true);
                $v = $p->getValue($comp);
                echo "  prop $prop type=".gettype($v);
                if (is_object($v)) echo " class=".get_class($v);
                if (is_callable($v) && !is_array($v)) {
                    try {
                        $r = $v();
                        echo " invoke_type=".gettype($r);
                        if ($r === null) echo " NULL!";
                    } catch (Throwable $e) {
                        echo " invoke_EX=".$e->getMessage();
                    }
                }
                echo "\n";
            }
        }
    } catch (Throwable $e) {
        echo "  reflect EX: ".$e->getMessage()."\n";
    }
}

// 直接测试参数默认 + 箭头捕获
function capture_data(array $widgets, array $data = []): array {
    $fns = [];
    foreach ($widgets as $w) {
        $fns[] = fn (): array => [
            ...[],
            ...$data,
        ];
    }
    return $fns;
}
$fns = capture_data(['A']);
try {
    $r = $fns[0]();
    echo "capture_ok count=".count($r)." data_was_array=y\n";
} catch (Throwable $e) {
    echo "capture EX: ".$e->getMessage()."\n";
}

echo "DONE\n";
