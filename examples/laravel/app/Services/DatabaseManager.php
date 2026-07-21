<?php

namespace App\Services;

use function Database\Sql\open;

/**
 * 数据库连接管理（由 AppServiceProvider 注册为 singleton）
 */
class DatabaseManager
{
    public function path(): string
    {
        return config('database.connections.sqlite.database')
            ?? dirname(dirname(__DIR__)) . '/storage/laravel.db';
    }

    public function connect(?string $dbPath = null)
    {
        if ($dbPath === null) {
            $dbPath = $this->path();
        }

        $dir = dirname($dbPath);
        if (!is_dir($dir)) {
            mkdir($dir, 0755, true);
        }

        $db = open('sqlite', $dbPath);
        $db->ping();
        \Database\registerDefaultConnection($db);

        return $db;
    }
}
