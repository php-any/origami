<?php

require __DIR__.'/../../vendor/autoload.php';

$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

echo "has method macro: ".(method_exists(Illuminate\Database\Eloquent\Builder::class, 'macro') ? 'yes' : 'no')."\n";
echo "class_uses: ";
print_r(class_uses(Illuminate\Database\Eloquent\Builder::class));

try {
    Illuminate\Database\Eloquent\Builder::macro('fooTest', fn () => 'bar');
    echo "macro call ok\n";
} catch (Throwable $e) {
    echo "macro call error: ".$e->getMessage()."\n";
}
