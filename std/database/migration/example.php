<?php

use Database\Sql\open;

$db = open("sqlite", ":memory:");
$db->ping();
\Database\registerDefaultConnection($db);

$entityDir = __DIR__ . "/fixtures";
$result = \Database\migrate($db, $entityDir);

echo "created: " . $result->createdCount . "\n";
echo "altered: " . $result->alteredCount . "\n";

$result2 = \Database\migrate($db, $entityDir);
echo "second run created: " . $result2->createdCount . "\n";
echo "second run altered: " . $result2->alteredCount . "\n";

$rows = \Database\DB::query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name");
foreach ($rows as $row) {
    echo "table: " . $row->name . "\n";
}
