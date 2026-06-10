<?php

namespace App;

use Net\Annotation\Application;
use Spring\Config\AppConfig;
use Spring\Config\DatabaseBootstrap;

/**
 * Spring 应用引导类
 *
 * #[Application] 在 flash 加载本文件时会：
 * 1. 扫描 scan 目录下的控制器并注册路由
 * 2. 扫描完成后由框架调用 boot() 一次
 * 3. 将 exit() 注册为 shutdown 回调，脚本结束时自动执行
 */
#[Application(name: 'spring', scan: __DIR__)]
class SpringApplication {

    public static function boot(): void {
        $dbPath = dirname(__DIR__) . '/' . AppConfig::DB_PATH;
        DatabaseBootstrap::init($dbPath);

        Log::info("========================================");
        Log::info(AppConfig::APP_NAME . " v" . AppConfig::APP_VERSION . " 引导完成");
        Log::info("========================================");
    }

    public static function exit(): void {
        Log::info(AppConfig::APP_NAME . " 正在关闭...");
    }
}
