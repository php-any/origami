<?php

require __DIR__.'/../../vendor/autoload.php';

class OrigamiContextualRepo
{
    public $connection;

    public function __construct(string $connection, ?int $chunkSize = null)
    {
        $this->connection = $connection;
    }
}

interface OrigamiEntriesContract
{
}

class OrigamiEntriesRepo implements OrigamiEntriesContract
{
    public $connection;

    public function __construct(string $connection, ?int $chunkSize = null)
    {
        $this->connection = $connection;
    }
}

$ref = new ReflectionClass(OrigamiContextualRepo::class);
$ctor = $ref->getConstructor();
$params = $ctor->getParameters();
$n0 = $params[0]->getName();
if ($n0 !== 'connection') {
    fwrite(STDERR, "ctor getName=".var_export($n0, true)." nameProp=".var_export($params[0]->name ?? null, true)."\n");
    exit(1);
}

$c = new Illuminate\Container\Container();
$c->when(OrigamiContextualRepo::class)
    ->needs('$connection')
    ->give(fn () => 'sqlite');

try {
    $repo = $c->make(OrigamiContextualRepo::class);
} catch (Throwable $e) {
    fwrite(STDERR, "make failed: ".$e->getMessage()."\n");
    exit(1);
}

if (!is_object($repo) || $repo->connection !== 'sqlite') {
    fwrite(STDERR, "connection=".var_export(is_object($repo) ? ($repo->connection ?? null) : $repo, true)."\n");
    exit(1);
}

$c2 = new Illuminate\Container\Container();
$c2->singleton(OrigamiEntriesContract::class, OrigamiEntriesRepo::class);
$c2->when(OrigamiEntriesRepo::class)
    ->needs('$connection')
    ->give(fn () => 'telescope');

try {
    $repo2 = $c2->make(OrigamiEntriesContract::class);
} catch (Throwable $e) {
    fwrite(STDERR, "interface make failed: ".$e->getMessage()."\n");
    exit(1);
}

if (!is_object($repo2) || $repo2->connection !== 'telescope') {
    fwrite(STDERR, "interface connection=".var_export(is_object($repo2) ? ($repo2->connection ?? null) : $repo2, true)."\n");
    exit(1);
}

echo "container_contextual_ok\n";
