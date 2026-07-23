<?php

/**
 * Eloquent Capsule 引导。
 * Capsule 构造会改写 config；构造后重新 instance 一份完整 Fluent 配置。
 */

use Illuminate\Container\Container as IlluminateContainer;
use Illuminate\Database\Capsule\Manager as Capsule;
use Illuminate\Support\Fluent;

function bootstrap_eloquent(?string $databasePath = null): Capsule
{
    static $capsule = null;

    if ($capsule !== null) {
        return $capsule;
    }

    $dbConfig = config('database');
    if (!is_array($dbConfig)) {
        $dbConfig = [];
    }

    $default = (string) ($dbConfig['default'] ?? 'sqlite');
    $connections = $dbConfig['connections'] ?? [];
    if (!is_array($connections)) {
        $connections = [];
    }

    if ($databasePath !== null && isset($connections['sqlite']) && is_array($connections['sqlite'])) {
        $connections['sqlite']['database'] = $databasePath;
    }

    $makeConfig = static function () use ($default, $connections): Fluent {
        return new Fluent([
            'database.fetch' => PDO::FETCH_OBJ,
            'database.default' => $default,
            'database.connections' => $connections,
            'database.dbal.types' => [],
        ]);
    };

    $container = new IlluminateContainer();
    $container->instance('config', $makeConfig());

    $capsule = new Capsule($container);
    // Capsule::setupDefaultConfiguration 会改写 default；换回完整配置
    $container->instance('config', $makeConfig());

    $capsule->setAsGlobal();
    $capsule->bootEloquent();

    return $capsule;
}

function eloquent_capsule(): Capsule
{
    return bootstrap_eloquent();
}

/**
 * #[Table] Entity 扫描目录（Database\migrate）。
 */
function entity_model_dir(): string
{
    return base_path('app/Models/Entity');
}
