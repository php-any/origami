<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Filament\Schemas\Schema;
use Filament\Schemas\Components\Component;

$schema = Schema::make();
echo "schema_class=".get_class($schema)."\n";

// 反射调用 renderEmbeddedHtml / toEmbeddedHtml
try {
    $ref = new ReflectionClass($schema);
    $m = $ref->getMethod('renderEmbeddedHtml');
    $m->setAccessible(true);
    $html = $m->invoke($schema);
    echo "render_type=".gettype($html)." len=".(is_string($html)?strlen($html):-1)."\n";
} catch (Throwable $e) {
    echo "render EX: ".$e->getMessage()."\n";
}

try {
    $html = $schema->toEmbeddedHtml();
    echo "embed_type=".gettype($html)." is_string=".(is_string($html)?'y':'n')."\n";
    if (is_object($html)) echo "embed_class=".get_class($html)."\n";
} catch (Throwable $e) {
    echo "embed EX: ".$e->getMessage()."\n";
    echo $e->getFile().":".$e->getLine()."\n";
}

// 逐步：只 withVisibilityCache + schema render
try {
    $html = Component::withVisibilityCache(fn (): string => $schema->renderEmbeddedHtml());
    echo "via_cache_type=".gettype($html)."\n";
} catch (Throwable $e) {
    echo "via_cache EX: ".$e->getMessage()."\n";
}

echo "DONE\n";
