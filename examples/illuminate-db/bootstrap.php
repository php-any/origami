<?php
/**
 * Illuminate Database（Eloquent）独立引导。
 *
 * 本文件将 Illuminate\Database\Capsule\Manager 配置为全局可用，
 * 独立于 Laravel 框架使用。
 */

use Illuminate\Database\Capsule\Manager as Capsule;

/**
 * 初始化数据库连接。
 *
 * @param string|null $database 数据库路径（SQLite 使用；:memory: 为内存库）
 * @return Capsule
 */
function bootstrap_capsule(?string $database = ':memory:'): Capsule
{
    static $capsule = null;

    if ($capsule !== null) {
        return $capsule;
    }

    $capsule = new Capsule();

    $config = [
        'driver'   => 'sqlite',
        'database' => $database,
        'prefix'   => '',
    ];

    // 支持通过环境变量覆盖数据库路径
    $envDb = getenv('DB_DATABASE');
    if ($envDb !== false && $envDb !== '') {
        $config['database'] = $envDb;
    }

    $capsule->addConnection($config);

    // 设为全局可访问
    $capsule->setAsGlobal();

    // 启动 Eloquent ORM
    $capsule->bootEloquent();

    return $capsule;
}
