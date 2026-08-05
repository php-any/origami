<?php

/**
 * 官方 Storage 组件可加载并可实例化。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope_http();

use Laravel\Telescope\Contracts\EntriesRepository;
use Laravel\Telescope\Storage\DatabaseEntriesRepository;

$config = illuminate_container()->make('config');
$connection = (string) $config->get('telescope.storage.database.connection', 'sqlite');
$chunk = (int) $config->get('telescope.storage.database.chunk', 1000);

$repo = new DatabaseEntriesRepository($connection, $chunk);
if (!$repo instanceof EntriesRepository) {
    echo "FAIL: DatabaseEntriesRepository is not EntriesRepository\n";
    exit(1);
}

echo "PASS\n";
