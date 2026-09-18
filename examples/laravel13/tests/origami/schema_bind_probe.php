<?php
/**
 * Reproduce schema container binding for Dashboard widgets.
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;
use Filament\Schemas\Schema;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

$page = app(\App\Filament\Pages\Dashboard::class);

// Full content schema path (what Dashboard actually renders)
try {
    $schema = $page->getContentSchema();
    echo "content_schema=".get_class($schema)."\n";
    if (method_exists($schema, 'getComponents')) {
        $comps = $schema->getComponents();
        echo "content_comps=".count($comps)."\n";
    }
    if (method_exists($schema, 'toHtml')) {
        $h = $schema->toHtml();
        echo "content_html_len=".strlen((string)$h)."\n";
    } elseif (method_exists($schema, 'toEmbeddedHtml')) {
        $h = $schema->toEmbeddedHtml();
        echo "content_embed_len=".strlen((string)$h)."\n";
    }
} catch (Throwable $e) {
    echo "content EX: ".$e->getMessage()."\n".$e->getFile().":".$e->getLine()."\n";
}

// Try binding via schema() builder
try {
    $widgets = $page->getWidgets();
    $schemaComps = $page->getWidgetsSchemaComponents($widgets);
    echo "raw_schema_count=".count($schemaComps)."\n";
    $s = Schema::make($page)
        ->components($schemaComps);
    echo "schema_made\n";
    $bound = $s->getComponents();
    echo "bound_count=".count($bound)."\n";
    foreach ($bound as $i => $c) {
        try {
            $cont = method_exists($c,'getContainer') ? $c->getContainer() : null;
            echo "bound$i container=".(is_object($cont)?get_class($cont):var_export($cont,true))."\n";
            if (method_exists($c,'toEmbeddedHtml')) {
                $eh = $c->toEmbeddedHtml();
                echo "bound$i embed_len=".strlen((string)$eh)."\n";
            }
        } catch (Throwable $e) {
            echo "bound$i EX: ".$e->getMessage()."\n";
        }
    }
} catch (Throwable $e) {
    echo "bind EX: ".$e->getMessage()."\n".$e->getFile().":".$e->getLine()."\n";
}

echo "DONE\n";
