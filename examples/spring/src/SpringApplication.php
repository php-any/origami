<?php

namespace Spring;

use Net\Annotation\Application;
use Spring\Config\DatabaseBootstrap;

/**
 * Spring 应用引导类
 *
 * index.php 通过 $server->boot(SpringApplication::class) 加载本类后，#[Application] 会：
 * 1. 按 scan 参数扫描控制器/服务并注册路由与 IoC 绑定
 * 2. 扫描完成后调用 boot() 一次
 * 3. 将 exit() 注册为 shutdown 回调
 */
#[Application(name: 'spring', scan: __DIR__)]
class SpringApplication {

    public static function boot(): void {
        $dbPath = dirname(__DIR__) . '/spring.db';
        DatabaseBootstrap::init($dbPath);

        \Log::info("========================================");
        \Log::info("Spring Demo v1.0.0 引导完成");
        \Log::info("IoC 容器已就绪（服务通过 #[Singleton] 等注解自动注册）");
        \Log::info("WebSocket 聊天室: ws://0.0.0.0:8080/ws/chat");
        \Log::info("========================================");
    }

    public static function exit(): void {
        \Log::info("Spring Demo 正在关闭...");
    }
}
