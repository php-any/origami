<?php

/**
 * Telescope 桥接：建表、Facade、EntriesRepository 绑定。
 */

use Bootstrap\Telescope\OrigamiEntriesRepository;
use Illuminate\Support\Facades\Facade;
use Laravel\Telescope\Contracts\ClearableRepository;
use Laravel\Telescope\Contracts\EntriesRepository;
use Laravel\Telescope\Contracts\PrunableRepository;
use Laravel\Telescope\Storage\DatabaseEntriesRepository;

function bootstrap_telescope(): void
{
    static $done = false;
    if ($done) {
        return;
    }

    $app = illuminate_container();
    bootstrap_eloquent();

    // Origami 上 CombGenerator/orderedUuid 依赖的大整数乘法未完备；统一用 uuid4
    \Illuminate\Support\Str::createUuidsUsing(static function () {
        return \Ramsey\Uuid\Uuid::uuid4();
    });

    Facade::setFacadeApplication($app);

    // Illuminate\Support\Carbon::format 在 Origami 上 parent:: 会死循环；Date/now 改走 nesbot Carbon
    \Illuminate\Support\DateFactory::useFactory(new \Bootstrap\Telescope\CarbonDateFactory());

    // 确保 now() / Date facade 可用（IncomingEntry::recordedAt）
    if (!$app->bound('date')) {
        $app->singleton('date', function () {
            return new \Illuminate\Support\DateFactory();
        });
    }

    // Telescope::store 捕获异常后走 ExceptionHandler
    if (!$app->bound(\Illuminate\Contracts\Debug\ExceptionHandler::class)) {
        $app->singleton(
            \Illuminate\Contracts\Debug\ExceptionHandler::class,
            \Illuminate\Foundation\Exceptions\Handler::class
        );
    }

    if (!$app->bound('db')) {
        $app->instance('db', eloquent_capsule()->getDatabaseManager());
    }

    $config = $app->make('config');
    if ($config->get('telescope') === null) {
        $path = config_path('telescope.php');
        if (is_file($path)) {
            $config->set('telescope', require $path);
        }
    }

    $connection = (string) $config->get('telescope.storage.database.connection', 'sqlite');
    $chunk = (int) $config->get('telescope.storage.database.chunk', 1000);

    telescope_migrate_tables($connection);

    $app->singleton(EntriesRepository::class, function () use ($connection, $chunk) {
        return new OrigamiEntriesRepository($connection, $chunk);
    });
    $app->singleton(ClearableRepository::class, function ($app) {
        return $app->make(EntriesRepository::class);
    });
    $app->singleton(PrunableRepository::class, function ($app) {
        return $app->make(EntriesRepository::class);
    });
    $app->alias(EntriesRepository::class, DatabaseEntriesRepository::class);
    $app->alias(EntriesRepository::class, OrigamiEntriesRepository::class);

    $done = true;
}

function telescope_migrate_tables(string $connection): void
{
    $db = eloquent_capsule()->getConnection($connection);

    // 用原生 SQL，避开 Blueprint 在 Origami 上的异常路径问题
    $db->statement('CREATE TABLE IF NOT EXISTS telescope_entries (
        sequence INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid VARCHAR(36) NOT NULL,
        batch_id VARCHAR(36) NOT NULL,
        family_hash VARCHAR(255) NULL,
        should_display_on_index INTEGER NOT NULL DEFAULT 1,
        type VARCHAR(20) NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME NULL
    )');

    $db->statement('CREATE UNIQUE INDEX IF NOT EXISTS telescope_entries_uuid_unique ON telescope_entries (uuid)');
    $db->statement('CREATE INDEX IF NOT EXISTS telescope_entries_batch_id_index ON telescope_entries (batch_id)');
    $db->statement('CREATE INDEX IF NOT EXISTS telescope_entries_type_index ON telescope_entries (type)');

    $db->statement('CREATE TABLE IF NOT EXISTS telescope_entries_tags (
        entry_uuid VARCHAR(36) NOT NULL,
        tag VARCHAR(255) NOT NULL,
        PRIMARY KEY (entry_uuid, tag)
    )');

    $db->statement('CREATE TABLE IF NOT EXISTS telescope_monitoring (
        tag VARCHAR(255) NOT NULL PRIMARY KEY
    )');
}

function telescope_entries_repository(): EntriesRepository
{
    bootstrap_telescope();

    return illuminate_container()->make(EntriesRepository::class);
}

/**
 * HTTP 侧 Telescope：鉴权放行、filter、内存 cache（pause recording）。
 */
function bootstrap_telescope_http(): void
{
    static $done = false;
    if ($done) {
        return;
    }

    bootstrap_telescope();

    $app = illuminate_container();

    if (!$app->bound('cache')) {
        $app->singleton('cache', function () {
            return new \Illuminate\Cache\Repository(new \Illuminate\Cache\ArrayStore());
        });
        $app->alias('cache', \Illuminate\Contracts\Cache\Repository::class);
    }

    \Laravel\Telescope\Telescope::auth(static function () {
        return true;
    });

    \Laravel\Telescope\Telescope::filter(static function () {
        return true;
    });

    $done = true;
}

/**
 * 将 EntryResult 转为 JSON 数组（避免 createdAt 非 Carbon 时崩溃）。
 */
function telescope_entry_to_array($entry): array
{
    if ($entry instanceof \Laravel\Telescope\EntryResult) {
        try {
            $entry->generateAvatar();
        } catch (\Throwable $e) {
            // avatar 可选
        }

        $created = $entry->createdAt;
        if (is_string($created)) {
            $createdAt = $created;
        } elseif (is_object($created) && method_exists($created, 'toDateTimeString')) {
            $createdAt = $created->toDateTimeString();
        } else {
            $createdAt = (string) $created;
        }

        $content = $entry->content;
        if (!is_array($content)) {
            $content = [];
        }

        return [
            'id' => (string) $entry->id,
            'sequence' => is_numeric($entry->sequence) ? (int) $entry->sequence : $entry->sequence,
            'batch_id' => (string) $entry->batchId,
            'type' => (string) $entry->type,
            'content' => $content,
            'tags' => [],
            'family_hash' => $entry->familyHash,
            'created_at' => $createdAt,
        ];
    }

    return (array) $entry;
}

function telescope_watcher_status(string $watcherClass): string
{
    if (!config('telescope.enabled', false)) {
        return 'disabled';
    }

    try {
        if (cache('telescope:pause-recording')) {
            return 'paused';
        }
    } catch (\Throwable $e) {
        // cache 未绑定时视为 recording
    }

    $watcher = config('telescope.watchers.' . $watcherClass);
    if (!$watcher || (isset($watcher['enabled']) && !$watcher['enabled'])) {
        return 'off';
    }

    return 'enabled';
}
