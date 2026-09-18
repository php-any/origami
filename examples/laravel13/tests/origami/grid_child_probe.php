<?php
/**
 * Inspect Grid child schema components after Dashboard getSchema('content').
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;
use Livewire\Livewire;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

$component = Livewire::new(\App\Filament\Pages\Dashboard::class);
$schema = $component->getSchema('content');
$grid = $schema->getComponents()[0];
$child = $grid->getChildSchema();
echo "child_class=".get_class($child)."\n";

$comps = $child->getComponents(withHidden: true);
echo "child_comps_withHidden=".count($comps)."\n";
$comps2 = $child->getComponents();
echo "child_comps=".count($comps2)."\n";

foreach ($comps as $i => $c) {
    echo "c$i=".get_class($c);
    try {
        $hidden = method_exists($c,'isHidden') ? ($c->isHidden()?'yes':'no') : '?';
        echo " hidden=$hidden";
        $cont = method_exists($c,'getContainer') ? $c->getContainer() : null;
        echo " cont=".(is_object($cont)?'yes':'null');
    } catch (Throwable $e) {
        echo " EX=".$e->getMessage();
    }
    echo "\n";
    try {
        $h = $c->toEmbeddedHtml();
        echo "  embed=".strlen((string)$h)."\n";
    } catch (Throwable $e) {
        echo "  embedEX: ".$e->getMessage()."\n";
    }
}

try {
    $html = $child->toHtml();
    echo "child_toHtml_len=".strlen((string)$html)."\n";
    echo "child_toHtml=".substr(preg_replace('/\s+/',' ',(string)$html),0,200)."\n";
} catch (Throwable $e) {
    echo "toHtml EX: ".$e->getMessage()."\n";
}

try {
    $html = $child->toEmbeddedHtml();
    echo "child_toEmbedded_len=".strlen((string)$html)."\n";
} catch (Throwable $e) {
    echo "toEmbedded EX: ".$e->getMessage()."\n";
}

// Directly evaluate widget schema closure
try {
    $direct = $component->getWidgetsSchemaComponents($component->getWidgets());
    echo "direct_count=".count($direct)."\n";
} catch (Throwable $e) {
    echo "direct EX: ".$e->getMessage()."\n";
}

echo "DONE\n";
