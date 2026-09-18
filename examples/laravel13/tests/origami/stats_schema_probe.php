<?php
/**
 * 拆开 Stats Schema 的 renderEmbeddedHtml / toEmbeddedHtml。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use Illuminate\Support\Facades\Auth;
use Filament\Facades\Filament;
use Filament\Schemas\Components\Component;

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

$widget = app(\App\Filament\Widgets\StatsOverviewWidget::class);
echo "widget=".get_class($widget)."\n";

try {
    $stats = $widget->getStats();
    echo "stats_count=".count($stats)." first=".get_class($stats[0])."\n";
} catch (Throwable $e) {
    echo "getStats EX: ".$e->getMessage()."\n";
}

try {
    // Livewire mount 前先拿 schema
    $schema = $widget->getSchema('content');
    echo "schema=".get_class($schema)."\n";

    $ref = new ReflectionClass($schema);
    $m = $ref->getMethod('renderEmbeddedHtml');
    $m->setAccessible(true);

    try {
        $html = $m->invoke($schema);
        echo "render_type=".gettype($html)." len=".(is_string($html)?strlen($html):-1)."\n";
        if (is_object($html)) echo "render_obj=".get_class($html)."\n";
    } catch (Throwable $e) {
        echo "render EX: ".$e->getMessage()."\n";
        echo $e->getFile().":".$e->getLine()."\n";
    }

    try {
        $html = $schema->toEmbeddedHtml();
        echo "embed_type=".gettype($html)."\n";
        if (is_object($html)) echo "embed_obj=".get_class($html)."\n";
        if (is_string($html)) echo "embed_len=".strlen($html)."\n";
    } catch (Throwable $e) {
        echo "embed EX: ".$e->getMessage()."\n";
        echo $e->getFile().":".$e->getLine()."\n";
    }

    try {
        $html = Component::withVisibilityCache(function () use ($schema, $m) {
            $r = $m->invoke($schema);
            echo "inner_render_type=".gettype($r)."\n";
            return $r;
        });
        echo "cache_type=".gettype($html)."\n";
    } catch (Throwable $e) {
        echo "cache EX: ".$e->getMessage()."\n";
    }
} catch (Throwable $e) {
    echo "schema EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

echo "DONE\n";
