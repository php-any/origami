<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;

class ServeCommand extends Command
{
    protected $signature = 'serve'
        . ' {--host=127.0.0.1 : 服务器监听地址}'
        . ' {--port=8000 : 服务器监听端口}';

    protected $description = '使用 Origami 内置 HTTP 服务器启动 Laravel 应用';

    public function handle()
    {
        $host = $this->option('host');
        $port = $this->option('port');

        putenv("SERVER_HOST={$host}");
        putenv("SERVER_PORT={$port}");

        $this->info("Laravel 开发服务器启动在: http://{$host}:{$port}");
        $this->comment('  按 Ctrl+C 停止服务器');

        // 加载 server.php 启动 HTTP 服务器
        require base_path('server.php');
    }
}
