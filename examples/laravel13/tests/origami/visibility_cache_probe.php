<?php
/**
 * withVisibilityCache 应调用 callback 并返回其结果，不是返回闭包本身。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Filament\Schemas\Components\Component;

$result = Component::withVisibilityCache(fn (): string => 'hello-html');
echo "type=".gettype($result)."\n";
echo "is_string=".(is_string($result)?'y':'n')."\n";
echo "is_callable=".(is_callable($result)?'y':'n')."\n";
if (is_object($result)) echo "class=".get_class($result)."\n";
echo "value=";
var_export($result);
echo "\n";

// Schema::toEmbeddedHtml 简化
use Filament\Schemas\Schema;
try {
    $schema = Schema::make();
    // might need more setup
    echo "schema_ok\n";
} catch (Throwable $e) {
    echo "schema EX: ".$e->getMessage()."\n";
}

echo "DONE\n";
