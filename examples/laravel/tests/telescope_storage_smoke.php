<?php

/**
 * Telescope 真实存储路径：IncomingEntry → Repository::store → find。
 */

require dirname(__DIR__) . '/bootstrap/app.php';
bootstrap_app();
bootstrap_telescope();

use Illuminate\Support\Collection;
use Illuminate\Support\Str;
use Laravel\Telescope\EntryType;
use Laravel\Telescope\IncomingEntry;

$repo = telescope_entries_repository();

$uuid = (string) Str::uuid();
$batchId = (string) Str::uuid();

$entry = new IncomingEntry([
    'level' => 'info',
    'message' => 'telescope storage smoke',
    'context' => [],
], $uuid);
$entry->type(EntryType::LOG);
$entry->batchId($batchId);

try {
    $repo->store(new Collection([$entry]));
} catch (\Throwable $e) {
    echo "FAIL: store: " . $e->getMessage() . "\n";
    exit(1);
}

try {
    $found = $repo->find($uuid);
} catch (\Throwable $e) {
    echo "FAIL: find: " . $e->getMessage() . "\n";
    exit(1);
}

if ((string) $found->id !== $uuid && (string) ($found->uuid ?? '') !== $uuid) {
    // EntryResult 用 id 存 uuid
    if ((string) $found->id !== $uuid) {
        echo "FAIL: uuid mismatch id=" . (string) $found->id . "\n";
        exit(1);
    }
}

if (($found->type ?? null) !== EntryType::LOG) {
    echo "FAIL: type\n";
    var_export($found->type ?? null);
    echo "\n";
    exit(1);
}

$content = $found->content ?? null;
if (!is_array($content) || ($content['message'] ?? null) !== 'telescope storage smoke') {
    echo "FAIL: content\n";
    var_export($content);
    echo "\n";
    exit(1);
}

echo "PASS\n";
