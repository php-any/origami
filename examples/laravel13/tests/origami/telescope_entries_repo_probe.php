<?php

require __DIR__.'/../../vendor/autoload.php';

$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
$kernel->bootstrap();

$abstract = Laravel\Telescope\Contracts\EntriesRepository::class;
$concrete = Laravel\Telescope\Storage\DatabaseEntriesRepository::class;

$ref = new ReflectionClass($concrete);
$params = $ref->getConstructor()->getParameters();
$names = [];
foreach ($params as $p) {
    $names[] = $p->getName();
}
echo 'ctor_names='.json_encode($names)."\n";

if (($names[0] ?? null) !== 'connection') {
    fwrite(STDERR, "DatabaseEntriesRepository ctor getName expected connection\n");
    exit(1);
}

try {
    $repo = $app->make($abstract);
} catch (Throwable $e) {
    fwrite(STDERR, "app make EntriesRepository failed: ".$e->getMessage()."\n");
    exit(1);
}

if (!is_object($repo) || !($repo instanceof $concrete)) {
    fwrite(STDERR, "repo type=".get_debug_type($repo)."\n");
    exit(1);
}

echo "telescope_entries_repo_ok\n";
