<?php

namespace Spring\Service;

use Container\Singleton;
use Database\DB;
use Spring\Model\Entity\UserEntity;

/**
 * 数据库查询示例服务
 *
 * 演示：先 query 查数据，再 toEntity 转实体；JOIN/聚合直接返回行对象。
 */
#[Singleton]
class QueryDemoService {

    /**
     * 单表查询：查出行数据 → 映射为 UserEntity
     */
    public function singleTableQuery(int $minAge = 0, int $limit = 10): array {
        if ($limit <= 0) {
            $limit = 10;
        }
        $rows = DB::query(
            "SELECT * FROM users WHERE age >= ? ORDER BY age DESC LIMIT ?",
            $minAge,
            $limit
        );
        return DB::toEntity(UserEntity::class, $rows);
    }

    /**
     * 单表模糊搜索
     */
    public function searchUsers(string $keyword): array {
        $pattern = "%" . $keyword . "%";
        $rows = DB::query(
            "SELECT * FROM users WHERE name LIKE ? OR email LIKE ? ORDER BY name ASC",
            $pattern,
            $pattern
        );
        return DB::toEntity(UserEntity::class, $rows);
    }

    public function aggregateByCategory(): array {
        return DB::query(
            "SELECT category, COUNT(*) AS product_count, AVG(price) AS avg_price, MIN(price) AS min_price, MAX(price) AS max_price FROM products GROUP BY category ORDER BY product_count DESC"
        );
    }

    public function productsAboveCategoryAvg(): array {
        return DB::query("
            SELECT p.id, p.name, p.price, p.category
            FROM products p
            WHERE p.price > (
                SELECT AVG(price) FROM products WHERE category = p.category
            )
            ORDER BY p.category, p.price DESC
        ");
    }

    public function innerJoinOrderProducts(): array {
        return DB::query("
            SELECT orders.id, orders.quantity, orders.total_price, orders.status,
                   p.name AS product_name, p.category AS product_category
            FROM orders
            INNER JOIN products p ON orders.product_id = p.id
            ORDER BY orders.id DESC
        ");
    }

    public function leftJoinUserOrders(): array {
        return DB::query("
            SELECT u.id AS user_id, u.name AS user_name, u.email,
                   o.id AS order_id, o.total_price, o.status
            FROM users u
            LEFT JOIN orders o ON u.id = o.user_id
            ORDER BY u.id, o.id
        ");
    }

    public function orderDetails(): array {
        return DB::query("
            SELECT
                o.id AS order_id,
                u.name AS user_name,
                u.email AS user_email,
                p.name AS product_name,
                p.category AS product_category,
                o.quantity,
                o.total_price,
                o.status,
                o.created_at
            FROM orders o
            INNER JOIN users u ON o.user_id = u.id
            INNER JOIN products p ON o.product_id = p.id
            ORDER BY o.id DESC
        ");
    }

    public function completedOrderStats(): array {
        return DB::query("
            SELECT
                u.id AS user_id,
                u.name AS user_name,
                COUNT(o.id) AS order_count,
                SUM(o.total_price) AS total_spent
            FROM users u
            INNER JOIN orders o ON u.id = o.user_id
            WHERE o.status = 'completed'
            GROUP BY u.id, u.name
            ORDER BY total_spent DESC
        ");
    }

    public static function rowsToArray(array $rows): array {
        $result = [];
        foreach ($rows as $row) {
            if (is_array($row)) {
                $result[] = $row;
                continue;
            }
            if (method_exists($row, 'toArray')) {
                $result[] = $row->toArray();
                continue;
            }
            $item = [];
            foreach ($row as $key => $value) {
                $item[$key] = $value;
            }
            $result[] = $item;
        }
        return $result;
    }

    public static function entitiesToArray(array $entities): array {
        $result = [];
        foreach ($entities as $entity) {
            if (method_exists($entity, 'toArray')) {
                $result[] = $entity->toArray();
            } else {
                $result[] = $entity;
            }
        }
        return $result;
    }
}
