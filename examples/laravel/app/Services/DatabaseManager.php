<?php

namespace App\Services;

use function Database\Sql\open;

/**
 * 数据库连接管理：Origami Sql 连接 + Eloquent Capsule 双轨。
 */
class DatabaseManager
{
    public function path(): string
    {
        $database = config('database');
        if (is_array($database)) {
            $connections = $database['connections'] ?? null;
            if (is_array($connections)) {
                $sqlite = $connections['sqlite'] ?? null;
                if (is_array($sqlite) && isset($sqlite['database']) && is_string($sqlite['database']) && $sqlite['database'] !== '') {
                    return $sqlite['database'];
                }
            }
        }

        return dirname(dirname(__DIR__)) . '/storage/laravel.db';
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
        bootstrap_eloquent($dbPath);

        return $db;
    }
}
