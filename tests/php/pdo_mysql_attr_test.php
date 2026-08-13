<?php

namespace tests\php;

/**
 * PDO::MYSQL_ATTR_LOCAL_INFILE_DIRECTORY 与 PHP 8.4 Pdo\Mysql::ATTR_*。
 */

$dir = \PDO::MYSQL_ATTR_LOCAL_INFILE_DIRECTORY;
if (!is_int($dir)) {
    Log::fatal('PDO::MYSQL_ATTR_LOCAL_INFILE_DIRECTORY 应为 int');
}

// 类常量默认值在解析期求值；缺失常量会让条件 class 也解析失败
if (false) {
    class PdoMysqlAttr_ConditionalStub
    {
        public const ATTR_LOCAL_INFILE_DIRECTORY = \PHP_VERSION_ID >= 80100
            ? \PDO::MYSQL_ATTR_LOCAL_INFILE_DIRECTORY
            : 1015;
    }
}

if (!class_exists('Pdo\\Mysql', false)) {
    Log::fatal('Pdo\\Mysql 应已由运行时注册');
}

$sslCa = \Pdo\Mysql::ATTR_SSL_CA;
if ($sslCa !== \PDO::MYSQL_ATTR_SSL_CA) {
    Log::fatal('Pdo\\Mysql::ATTR_SSL_CA 应等于 PDO::MYSQL_ATTR_SSL_CA');
}

$localDir = \Pdo\Mysql::ATTR_LOCAL_INFILE_DIRECTORY;
if ($localDir !== \PDO::MYSQL_ATTR_LOCAL_INFILE_DIRECTORY) {
    Log::fatal('Pdo\\Mysql::ATTR_LOCAL_INFILE_DIRECTORY 应对齐 PDO 常量');
}

Log::info('pdo_mysql_attr 测试通过');
