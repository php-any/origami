<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';

$boot = [
    Illuminate\Foundation\Bootstrap\LoadEnvironmentVariables::class,
    Illuminate\Foundation\Bootstrap\LoadConfiguration::class,
];
foreach ($boot as $b) {
    $app->make($b)->bootstrap($app);
}

$configProviders = $app->make('config')->get('app.providers') ?? [];
echo 'config_count='.count($configProviders)."\n";
echo 'config_has_view='.(in_array('Illuminate\\View\\ViewServiceProvider', $configProviders, true) ? 'y' : 'n')."\n";

$partitioned = (new Illuminate\Support\Collection($configProviders))
    ->partition(fn ($provider) => str_starts_with($provider, 'Illuminate\\'));

echo 'partition_count='.$partitioned->count()."\n";
$pkg = $app->make(Illuminate\Foundation\PackageManifest::class)->providers();
echo 'package_count='.count($pkg)."\n";

$removed = $partitioned->splice(1, 0, [$pkg]);
echo 'removed_count='.$removed->count()."\n";
echo 'partition_after_splice_count='.$partitioned->count()."\n";

$collapsed = $partitioned->collapse();
echo 'collapsed_count='.$collapsed->count()."\n";

$final = $collapsed->toArray();
echo 'final_count='.count($final)."\n";
echo 'final_has_view='.(in_array('Illuminate\\View\\ViewServiceProvider', $final, true) ? 'y' : 'n')."\n";
echo 'final_sample='.json_encode(array_slice($final, 0, 3))."\n";

$manifest = $app->getCachedServicesPath();
$cached = file_exists($manifest) ? include $manifest : null;
echo 'cached_providers_count='.(is_array($cached) ? count($cached['providers'] ?? []) : 0)."\n";
echo 'arrays_equal='.(is_array($cached) && ($cached['providers'] ?? null) === $final ? 'y' : 'n')."\n";
