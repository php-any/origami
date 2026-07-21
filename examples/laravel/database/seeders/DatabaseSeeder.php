<?php

namespace Database\Seeders;

use App\Models\Post;
use App\Models\User;
use App\Services\AuthService;
use Database\DB;

/**
 * 数据库种子（对应 Laravel database/seeders/DatabaseSeeder.php）
 */
class DatabaseSeeder
{
    public function run(bool $cli = true): void
    {
        $count = DB::query('SELECT COUNT(*) AS cnt FROM users');
        if ($count[0]->cnt > 0) {
            $this->out($cli, '已有数据，跳过种子');
            return;
        }

        $this->out($cli, '插入种子数据...');

        $users = [
            ['name' => 'Alice', 'email' => 'alice@example.com', 'password' => 'secret123'],
            ['name' => 'Bob', 'email' => 'bob@example.com', 'password' => 'secret123'],
        ];
        foreach ($users as $row) {
            $user = new User();
            $user->name = $row['name'];
            $user->email = $row['email'];
            $user->password = AuthService::hashPassword($row['password']);
            DB::insert($user);
        }

        $posts = [
            ['user_id' => 1, 'title' => '欢迎使用 Laravel 风格示例', 'body' => '这是一个基于 Origami 的 Laravel 风格 Web 框架演示，包含路由、控制器、中间件、Eloquent 风格 Model 与 Artisan CLI。'],
            ['user_id' => 1, 'title' => '声明式路由', 'body' => 'Web 与 API 路由统一在 routes/*.php 中声明，目录结构遵循 Laravel 约定。'],
            ['user_id' => 2, 'title' => 'IoC 容器与依赖注入', 'body' => '服务类通过 #[Singleton] 注册到容器，控制器构造函数自动注入依赖，体验接近 Laravel 的 Service Container。'],
        ];
        foreach ($posts as $row) {
            $post = new Post();
            $post->user_id = $row['user_id'];
            $post->title = $row['title'];
            $post->body = $row['body'];
            DB::insert($post);
        }

        $this->out($cli, '种子数据插入完成');
    }

    private function out(bool $cli, string $message): void
    {
        if ($cli) {
            echo $message . "\n";
            return;
        }
        \Log::info($message);
    }
}
