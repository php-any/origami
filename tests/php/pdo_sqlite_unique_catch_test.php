<?php

namespace tests\php;

use Exception;

/**
 * SQLite UNIQUE 应抛出可 catch 的 PDOException（Laravel Connection::runQueryCallback 依赖）。
 */

$path = sys_get_temp_dir().DIRECTORY_SEPARATOR.'origami_pdo_unique_catch.sqlite';
@unlink($path);
$pdo = new \PDO('sqlite:'.$path);
$pdo->setAttribute(\PDO::ATTR_ERRMODE, \PDO::ERRMODE_EXCEPTION);
$pdo->exec('CREATE TABLE t(uuid TEXT PRIMARY KEY)');
$pdo->exec("INSERT INTO t(uuid) VALUES('a')");

$caught = null;
try {
    $pdo->exec("INSERT INTO t(uuid) VALUES('a')");
} catch (Exception $e) {
    $caught = 'ex:'.get_class($e);
} catch (\PDOException $e) {
    $caught = 'pdo';
} catch (\Throwable $e) {
    $caught = 'th:'.get_class($e);
}
@unlink($path);

if ($caught === null || strpos((string) $caught, 'ex:') !== 0) {
    Log::fatal('UNIQUE 应被 catch (Exception) 接到 PDOException，实际: '.var_export($caught, true));
}

Log::info('pdo_sqlite_unique_catch 测试通过');
