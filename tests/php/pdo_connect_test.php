<?php

$pdo = PDO::connect('sqlite::memory:');

if (!($pdo instanceof PDO)) {
    echo "FAIL return type\n";
    exit(1);
}

if ($pdo->getAttribute(PDO::ATTR_SERVER_VERSION) === '') {
    echo "FAIL server version\n";
    exit(1);
}

$pdo->exec('CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT)');
$pdo->exec("INSERT INTO items (name) VALUES ('first')");

$name = $pdo->query('SELECT name FROM items')->fetchColumn();
if ($name !== 'first') {
    echo "FAIL query\n";
    exit(1);
}

echo "PASS\n";
