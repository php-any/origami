<?php

namespace tests\php;

/**
 * dirname() 第二参数 $levels 测试
 */

$path = '/var/www/html/app/Console/Commands/ServeCommand.php';

if (dirname($path) !== '/var/www/html/app/Console/Commands') {
    Log::fatal('dirname 单层测试失败: ' . dirname($path));
}

if (dirname($path, 3) !== '/var/www/html/app') {
    Log::fatal('dirname levels=3 测试失败: ' . dirname($path, 3));
}

if (dirname($path, 4) !== '/var/www/html') {
    Log::fatal('dirname levels=4 测试失败: ' . dirname($path, 4));
}

Log::info('dirname 测试通过');
