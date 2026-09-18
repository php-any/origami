<?php
/**
 * 对齐 Filament：collect()->map 延迟闭包 + getComponentProperties 式 ...getData()。
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

echo "step1 invoke deferred props closures\n";
try {
    $components = $page->getWidgetsSchemaComponents($widgets);
    echo "components=".count($components)."\n";
    foreach ($components as $i => $comp) {
        echo "c$i class=".get_class($comp)."\n";
        try {
            $props = $comp->getComponentProperties();
            echo "  props_type=".gettype($props)." is_array=".(is_array($props)?'y':'n')." count=".(is_array($props)?count($props):-1)."\n";
            echo "  props_keys=".json_encode(is_array($props)?array_keys($props):null)."\n";
        } catch (Throwable $e) {
            echo "  getComponentProperties EX: ".$e->getMessage()."\n";
            echo "  @ ".$e->getFile().":".$e->getLine()."\n";
        }
        try {
            $data = $comp->getData();
            echo "  getData_type=".gettype($data)."\n";
            if ($data === null) echo "  getData IS NULL\n";
            if (is_array($data)) echo "  getData_keys=".json_encode(array_keys($data))."\n";
        } catch (Throwable $e) {
            echo "  getData EX: ".$e->getMessage()."\n";
        }
    }
} catch (Throwable $e) {
    echo "schema EX: ".$e->getMessage()."\n";
}

echo "step2 mimic getComponentProperties spread\n";
try {
    $components = $page->getWidgetsSchemaComponents($widgets);
    foreach ($components as $i => $comp) {
        $properties = ['record' => null];
        $merged = [...$properties, ...$comp->getData()];
        echo "mimic$i count=".count($merged)."\n";
    }
    echo "mimic_ok\n";
} catch (Throwable $e) {
    echo "mimic EX: ".$e->getMessage()."\n";
}

echo "DONE\n";
