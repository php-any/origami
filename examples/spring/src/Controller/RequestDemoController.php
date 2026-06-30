<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;
use Net\Annotation\RequestBody;
use Net\Annotation\PathVariable;
use Net\Annotation\RequestParam;
use Net\Annotation\Route;
use Net\Http\Request;
use Net\Http\Response;

/**
 * 请求参数绑定演示
 *
 * 两种使用方式：
 *
 * 1. 方法级注解 — 显式声明每个参数的来源
 *    #[PathVariable('userId', 'postId')]  // 这些参数来自路由路径
 *    #[RequestParam('sort', 'page')]      // 这些参数来自 URL 查询
 *    #[RequestBody('data', 'user')]       // 这些参数来自请求体
 *
 * 2. 自动推断 — 无注解时按类型 + 参数名自动查找
 *    string/int/float/bool → POST > Query > Route 按名称
 *    array                 → 全部数据合并
 *    自定义类               → JSON body 绑定
 *
 * 优先级：POST 表单 > URL 查询参数 > 路由参数
 */
#[Controller]
#[Route(prefix: "/api/request-demo")]
class RequestDemoController {

    // ================================================================
    // 方法级注解 — 显式声明绑定来源
    // ================================================================

    /**
     * #[PathVariable] + #[RequestParam] 区分路由参数和查询参数
     *
     * GET /api/request-demo/users/42/posts/99?sort=date&page=1
     * → {"userId": "42", "postId": "99", "sort": "date", "page": 1}
     */
    #[GetMapping(path: "/users/{userId}/posts/{postId}")]
    #[PathVariable('userId', 'postId')]
    #[RequestParam('sort', 'page')]
    public function getPost(
        string $userId,
        string $postId,
        string $sort,
        int $page
    ): array {
        return [
            'userId' => $userId,   // 路由 {userId}
            'postId' => $postId,   // 路由 {postId}
            'sort'   => $sort,     // ?sort=
            'page'   => $page,     // ?page=
        ];
    }

    /**
     * #[RequestBody] array — 接收整个 JSON body 为数组
     *
     * POST /api/request-demo/users
     * Content-Type: application/json
     * Body: {"name": "Alice", "email": "alice@example.com", "age": 25}
     */
    #[PostMapping(path: "/users")]
    #[RequestBody('data')]
    public function createUser(array $data): array {
        return $data;
    }

    /**
     * #[RequestBody] 自定义类 — JSON body → 对象
     *
     * POST /api/request-demo/users/dto
     * Content-Type: application/json
     * Body: {"name": "Bob", "email": "bob@example.com", "age": 30}
     */
    #[PostMapping(path: "/users/dto")]
    #[RequestBody('user')]
    public function createUserDto(UserInfo $user): array {
        return [
            'name'  => $user->name,
            'email' => $user->email,
            'age'   => $user->age,
        ];
    }

    /**
     * #[RequestBody] 单个表单字段
     *
     * POST /api/request-demo/login
     * Body: username=admin&password=secret
     */
    #[PostMapping(path: "/login")]
    #[RequestBody('username', 'password')]
    public function login(string $username, string $password): array {
        return ['username' => $username, 'status' => 'ok'];
    }

    /**
     * 混合注解：路由 + 查询 + Body
     *
     * POST /api/request-demo/order/10?discount=0.8
     * Body: item=laptop&qty=2
     */
    #[PostMapping(path: "/order/{orderId}")]
    #[PathVariable('orderId')]
    #[RequestParam('discount')]
    #[RequestBody('item', 'qty')]
    public function createOrder(
        int $orderId,
        float $discount,
        string $item,
        int $qty
    ): array {
        return [
            'orderId'  => $orderId,   // 路由 {orderId}
            'discount' => $discount,  // ?discount=
            'item'     => $item,      // POST body
            'qty'      => $qty,       // POST body
        ];
    }

    // ================================================================
    // 无注解 — 自动推断
    // ================================================================

    /**
     * 无注解：按参数名从 POST > Query > Route 自动查找
     *
     * GET /api/request-demo/users/42?name=Alice
     */
    #[GetMapping(path: "/users/{id}")]
    public function autoBind(string $id, string $name): array {
        return ['id' => $id, 'name' => $name];
    }

    /**
     * 选择性使用 $request 对象
     */
    #[GetMapping(path: "/debug")]
    #[RequestParam('key')]
    public function debug(Request $request, string $key): array {
        return [
            'key'     => $key,              // 注解声明来自查询参数
            'allData' => $request->all(),   // $request 对象方法
        ];
    }
}

class UserInfo {
    public string $name;
    public string $email;
    public int $age;
}
