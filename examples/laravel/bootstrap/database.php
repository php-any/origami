<?php

/**
 * Eloquent Capsule 引导。
 * 挂到 Foundation Application，避免双容器。
 */

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

    $container = illuminate_container();

    // Capsule 会改写 config；用 Fluent 临时覆盖 database.*，构造后再恢复 Repository
    $originalConfig = $container->make('config');
    $container->instance('config', $makeConfig());

    $capsule = new Capsule($container);
    $container->instance('config', $makeConfig());

    $capsule->setAsGlobal();
    $capsule->bootEloquent();

    // 恢复完整 Config Repository，并写回 database 连接
    if ($originalConfig instanceof \Illuminate\Config\Repository) {
        $originalConfig->set('database.default', $default);
        $originalConfig->set('database.connections', $connections);
        $container->instance('config', $originalConfig);
    }

    // 绑定 db 以便 Telescope / 其它组件解析
    if (!$container->bound('db')) {
        $container->instance('db', $capsule->getDatabaseManager());
    }

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
