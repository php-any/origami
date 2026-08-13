<?php

namespace App;

use Net\Annotation\Application;

/**
 * 应用引导类（业务入口）
 *
 * #[Application] 扫描 src/ 目录下的控制器与服务并注册路由。
 */
#[Application(name: "agents", scan: __DIR__)]
class AgentsApplication
{
    public static function boot(): void
    {
        \Log::info("Agents Web 应用已就绪");
    }

    public static function exit(): void
    {
        \Log::info("Agents Web 服务已关闭");
    }
}
