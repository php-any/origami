<?php

namespace App\Database;

use App\Services\DatabaseManager;
use Database\Seeders\DatabaseSeeder;

/**
 * 数据库初始化（迁移 + 种子），仅由 Artisan 命令调用
 */
class DatabaseBootstrap
{
    private const MODEL_DIR = __DIR__ . '/../Models';

    public static function migrate(?string $dbPath = null, bool $cli = false): void
    {
        if ($dbPath === null) {
            $dbPath = app_make(DatabaseManager::class)->path();
        }

        self::out($cli, '=== 数据库迁移 ===');
        self::out($cli, '数据库路径: ' . $dbPath);

        $db = app_make(DatabaseManager::class)->connect($dbPath);

        $result = \Database\migrate($db, self::MODEL_DIR);
        self::out($cli, 'Schema 同步: 新建表 ' . $result->createdCount . ' 个, 新增列 ' . $result->alteredCount . ' 个');
        foreach ($result->created as $item) {
            self::out($cli, '  创建表: ' . $item->table);
        }
    }

    public static function seed(bool $cli = false): void
    {
        (new DatabaseSeeder())->run($cli);
    }

    public static function migrateAndSeed(?string $dbPath = null, bool $cli = false): void
    {
        self::migrate($dbPath, $cli);
        self::seed($cli);
        self::out($cli, '数据库初始化完成');
    }

    private static function out(bool $cli, string $message): void
    {
        if ($cli) {
            echo $message . "\n";
            return;
        }
        \Log::info($message);
    }
}
