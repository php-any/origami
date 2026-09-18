<?php
/**
 * After Livewire mount Dashboard: inspect content schema render.
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

try {
    $component = Livewire::new(\App\Filament\Pages\Dashboard::class);
    echo "new_ok class=".get_class($component)."\n";
    // mimic mount lifecycle partially
    if (method_exists($component, 'mount')) {
        try { $component->mount(); echo "mount_method_ok\n"; } catch (Throwable $e) { echo "mount_method EX: ".$e->getMessage()."\n"; }
    }
    $schema = $component->getSchema('content');
    echo "schema=".(is_object($schema)?get_class($schema):var_export($schema,true))."\n";
    if ($schema) {
        $comps = $schema->getComponents();
        echo "top_comps=".count($comps)."\n";
        foreach ($comps as $i => $c) {
            echo "top$i=".get_class($c)."\n";
            try {
                if (method_exists($c, 'getChildSchema')) {
                    $child = $c->getChildSchema();
                    echo "  childSchema=".(is_object($child)?get_class($child):gettype($child))."\n";
                }
                if (method_exists($c, 'getChildSchemas')) {
                    $children = $c->getChildSchemas();
                    echo "  childSchemas=".count($children)."\n";
                }
                $eh = $c->toEmbeddedHtml();
                echo "  embed_len=".strlen((string)$eh)." has_wi=".(str_contains((string)$eh,'fi-wi-')||str_contains((string)$eh,'wire:id')?'yes':'no')."\n";
            } catch (Throwable $e) {
                echo "  EX: ".$e->getMessage()."\n";
            }
        }
        try {
            $html = $schema->toEmbeddedHtml();
            echo "schema_embed_len=".strlen((string)$html)."\n";
            echo "schema_has_wire=".(substr_count((string)$html,'wire:id'))."\n";
            echo "schema_has_lazy=".(str_contains((string)$html,'__lazyLoad')||str_contains((string)$html,'wire:init')?'yes':'no')."\n";
        } catch (Throwable $e) {
            echo "schema_embed EX: ".$e->getMessage()."\n".$e->getFile().":".$e->getLine()."\n";
        }
    }
} catch (Throwable $e) {
    echo "EX: ".$e->getMessage()."\n".$e->getFile().":".$e->getLine()."\n";
}
echo "DONE\n";
