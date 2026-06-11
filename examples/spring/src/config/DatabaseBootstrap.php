<?php

namespace Spring\Config;

use Database\DB;
use Database\Sql\open;
use Spring\Model\Entity\OrderEntity;
use Spring\Model\Entity\ProductEntity;
use Spring\Model\Entity\UserEntity;

/**
 * SQLite 数据库初始化：建表、种子数据
 */
class DatabaseBootstrap {

    public static function init($dbPath = null) {
        if ($dbPath === null) {
            $dbPath = __DIR__ . '/../../spring.db';
        }

        Log::info("=== 初始化 SQLite 数据库 ===");
        Log::info("数据库路径: " . $dbPath);

        $db = open("sqlite", $dbPath);
        $db->ping();
        Database\registerDefaultConnection($db);

        self::createTables($db);
        self::seedData();

        Log::info("数据库初始化完成");
    }

    private static function createTables($db) {
        $db->exec("
            CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name VARCHAR(100) NOT NULL,
                email VARCHAR(100) UNIQUE NOT NULL,
                age INTEGER DEFAULT 0,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP
            )
        ");

        $db->exec("
            CREATE TABLE IF NOT EXISTS products (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name VARCHAR(200) NOT NULL,
                price REAL NOT NULL DEFAULT 0,
                category VARCHAR(100) DEFAULT '未分类',
                description TEXT,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP
            )
        ");

        $db->exec("
            CREATE TABLE IF NOT EXISTS orders (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                user_id INTEGER NOT NULL,
                product_id INTEGER NOT NULL,
                quantity INTEGER NOT NULL DEFAULT 1,
                total_price REAL NOT NULL DEFAULT 0,
                status VARCHAR(50) DEFAULT 'pending',
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY (user_id) REFERENCES users (id),
                FOREIGN KEY (product_id) REFERENCES products (id)
            )
        ");
    }

    private static function seedData() {
        $userCount = DB::sql("SELECT COUNT(*) AS cnt FROM users");
        if ($userCount[0]->cnt > 0) {
            Log::info("数据库已有数据，跳过种子数据");
            return;
        }

        Log::info("插入种子数据...");

        $users = [
            ["name" => "张三", "email" => "zhangsan@example.com", "age" => 25],
            ["name" => "李四", "email" => "lisi@example.com", "age" => 30],
            ["name" => "王五", "email" => "wangwu@example.com", "age" => 28],
            ["name" => "赵六", "email" => "zhaoliu@example.com", "age" => 35],
        ];
        foreach ($users as $row) {
            $user = new UserEntity();
            $user->name = $row['name'];
            $user->email = $row['email'];
            $user->age = $row['age'];
            DB::insert($user);
        }

        $products = [
            ["name" => "iPhone 15 Pro", "price" => 7999.00, "category" => "电子产品", "description" => "Apple 最新旗舰手机"],
            ["name" => "MacBook Pro 14", "price" => 14999.00, "category" => "电子产品", "description" => "专业级笔记本电脑"],
            ["name" => "AirPods Pro", "price" => 1899.00, "category" => "电子产品", "description" => "无线降噪耳机"],
            ["name" => "iPad Air", "price" => 4799.00, "category" => "电子产品", "description" => "轻薄平板电脑"],
            ["name" => "机械键盘", "price" => 599.00, "category" => "配件", "description" => "Cherry 轴机械键盘"],
            ["name" => "显示器支架", "price" => 299.00, "category" => "配件", "description" => "铝合金显示器支架"],
        ];
        foreach ($products as $row) {
            $product = new ProductEntity();
            $product->name = $row['name'];
            $product->price = $row['price'];
            $product->category = $row['category'];
            $product->description = $row['description'];
            DB::insert($product);
        }

        $orders = [
            ["user_id" => 1, "product_id" => 1, "quantity" => 1, "total_price" => 7999.00, "status" => "completed"],
            ["user_id" => 1, "product_id" => 3, "quantity" => 2, "total_price" => 3798.00, "status" => "completed"],
            ["user_id" => 2, "product_id" => 2, "quantity" => 1, "total_price" => 14999.00, "status" => "pending"],
            ["user_id" => 3, "product_id" => 5, "quantity" => 1, "total_price" => 599.00, "status" => "completed"],
            ["user_id" => 2, "product_id" => 4, "quantity" => 1, "total_price" => 4799.00, "status" => "cancelled"],
        ];
        foreach ($orders as $row) {
            $order = new OrderEntity();
            $order->user_id = $row['user_id'];
            $order->product_id = $row['product_id'];
            $order->quantity = $row['quantity'];
            $order->total_price = $row['total_price'];
            $order->status = $row['status'];
            DB::insert($order);
        }

        Log::info("种子数据插入完成");
    }
}
