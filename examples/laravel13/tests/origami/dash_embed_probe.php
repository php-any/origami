<?php
/**
 * Diagnose Dashboard embed + HTTP /admin remaining Fatal.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

$html = (string) app('livewire')->mount(\App\Filament\Pages\Dashboard::class);
echo "dash_len=".strlen($html)."\n";
echo "has_fi_wi=".(str_contains($html,'fi-wi-')?'yes':'no')."\n";
echo "has_lazy=".(str_contains($html,'wire:snapshot') && str_contains($html,'lazy')?'maybe':'?')."\n";
echo "has_stats_text=".(str_contains($html,'订单')||str_contains($html,'Stats')||str_contains($html,'fi-wi-stats')?'yes':'no')."\n";
echo "has_livewire_child=".(str_contains($html,'wire:id')?('yes count='.substr_count($html,'wire:id')):'no')."\n";
if (preg_match('/"children":(\[[^\]]*\]|\{[^}]*\})/', $html, $m)) {
    echo "snapshot_children=".$m[1]."\n";
} else {
    echo "snapshot_children=NOT_FOUND\n";
}
echo "snippet=".substr(preg_replace('/\s+/',' ',strip_tags($html)),0,300)."\n";

// Schema components embed path
$page = app(\App\Filament\Pages\Dashboard::class);
$comps = $page->getWidgetsSchemaComponents($page->getWidgets());
echo "schema_count=".count($comps)."\n";
foreach ($comps as $i => $comp) {
    try {
        echo "comp$i class=".(is_object($comp)?get_class($comp):gettype($comp))."\n";
        if (is_object($comp) && method_exists($comp, 'toEmbeddedHtml')) {
            $eh = $comp->toEmbeddedHtml();
            echo "comp$i embed_len=".strlen((string)$eh)."\n";
        } elseif (is_object($comp) && method_exists($comp, 'toHtml')) {
            $eh = $comp->toHtml();
            echo "comp$i html_len=".strlen((string)$eh)."\n";
        }
    } catch (Throwable $e) {
        echo "comp$i EX: ".$e->getMessage()."\n";
    }
}
echo "DONE\n";
