<?php
/**
 * Bootstrap probe with Go-native Collection enabled.
 */
require __DIR__.'/../../vendor/autoload.php';

function step(string $label, callable $fn): void {
    echo "STEP {$label}... ";
    try {
        $fn();
        echo "OK\n";
    } catch (Throwable $e) {
        echo "FAIL: ".$e->getMessage()."\n";
        echo "  at ".$e->getFile().":".$e->getLine()."\n";
        exit(1);
    }
}

step('collect_class', function () {
    $c = collect([1]);
    echo '('.get_class($c).') ';
    if (!($c instanceof Illuminate\Support\Collection)) {
        throw new RuntimeException('not Collection instance');
    }
});

step('package_manifest_chain', function () {
    $packages = [
        ['name' => 'vendor/foo/pkg', 'extra' => ['laravel' => ['providers' => ['Foo\\Bar']]]],
        ['name' => 'vendor/baz/pkg', 'extra' => []],
    ];
    $ignore = [];
    $ignoreAll = false;
    $manifest = (new Illuminate\Support\Collection($packages))
        ->mapWithKeys(function ($package) {
            return [str_replace('vendor/', '', $package['name']) => $package['extra']['laravel'] ?? []];
        })
        ->each(function ($configuration) use (&$ignore) {
            $ignore = array_merge($ignore, $configuration['dont-discover'] ?? []);
        })
        ->reject(function ($configuration, $package) use ($ignore, $ignoreAll) {
            return $ignoreAll || in_array($package, $ignore);
        })
        ->filter()
        ->all();
    if (!is_array($manifest)) {
        throw new RuntimeException('manifest not array: '.gettype($manifest));
    }
    echo '(keys='.count($manifest).') ';
});

step('bootstrap_app', function () {
    global $app;
    $app = require __DIR__.'/../../bootstrap/app.php';
});

step('kernel_bootstrappers', function () {
    global $app;
    $kernel = $app->make(Illuminate\Contracts\Console\Kernel::class);
    $ref = new ReflectionClass($kernel);
    $prop = $ref->getProperty('bootstrappers');
    $prop->setAccessible(true);
    $bootstrappers = $prop->getValue($kernel);
    foreach ($bootstrappers as $bootstrapper) {
        echo "\n  bootstrapper {$bootstrapper}... ";
        try {
            $app->make($bootstrapper)->bootstrap($app);
            echo 'OK';
            if ($bootstrapper === Illuminate\Foundation\Bootstrap\RegisterProviders::class) {
                $loaded = $app->getLoadedProviders();
                $view = $loaded['Illuminate\\View\\ViewServiceProvider'] ?? false;
                $bound = $app->bound('blade.compiler');
                echo " view_sp=".($view?'y':'n')." blade_bound=".($bound?'y':'n');
                $providers = $app->make('config')->get('app.providers') ?? [];
                echo " provider_count=".count($providers);
                echo " has_view=".((in_array('Illuminate\\View\\ViewServiceProvider', $providers, true))?'y':'n');
            }
        } catch (Throwable $e) {
            echo 'FAIL: '.$e->getMessage();
            echo "\n  at ".$e->getFile().":".$e->getLine();
            throw $e;
        }
    }
    echo "\n";
});

step('blade_compiler', function () {
    global $app;
    $b = $app->make('blade.compiler');
    echo '('.get_class($b).') ';
});

step('view_provider_loaded', function () {
    global $app;
    $loaded = $app->getLoadedProviders();
    $ok = $loaded['Illuminate\\View\\ViewServiceProvider'] ?? false;
    if (!$ok) {
        throw new RuntimeException('ViewServiceProvider not loaded');
    }
});

echo "ALL_OK\n";
