<?php
/**
 * Same request: Grid toEmbeddedHtml vs child toHtml.
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
$viaChild = (string) $child->toHtml();
echo "via_child_len=".strlen($viaChild)."\n";

$viaGrid = (string) $grid->toEmbeddedHtml();
echo "via_grid_len=".strlen($viaGrid)."\n";
echo "via_grid_has_wire=".(str_contains($viaGrid,'wire:')?'yes':'no')."\n";
echo "via_grid_snip=".substr(preg_replace('/\s+/',' ',$viaGrid),0,180)."\n";

// Again after viaGrid — child still ok?
$viaChild2 = (string) $grid->getChildSchema()->toHtml();
echo "via_child2_len=".strlen($viaChild2)."\n";

echo "DONE\n";
