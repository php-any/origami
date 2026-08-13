<?php

/**
 * Telescope 类可加载冒烟。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();

use Laravel\Telescope\Contracts\EntriesRepository;
use Laravel\Telescope\IncomingEntry;
use Laravel\Telescope\Storage\DatabaseEntriesRepository;
use Laravel\Telescope\Telescope;
use Laravel\Telescope\TelescopeServiceProvider;

if (!class_exists(Telescope::class)) {
    echo "FAIL: Telescope class missing\n";
    exit(1);
}

if (!class_exists(IncomingEntry::class)) {
    echo "FAIL: IncomingEntry missing\n";
    exit(1);
}

if (!interface_exists(EntriesRepository::class) && !class_exists(EntriesRepository::class)) {
    // Contracts 是 interface
    if (!interface_exists(EntriesRepository::class)) {
        echo "FAIL: EntriesRepository missing\n";
        exit(1);
    }
}

if (!class_exists(DatabaseEntriesRepository::class)) {
    echo "FAIL: DatabaseEntriesRepository missing\n";
    exit(1);
}

if (!class_exists(TelescopeServiceProvider::class)) {
    echo "FAIL: TelescopeServiceProvider missing\n";
    exit(1);
}

echo "PASS\n";
