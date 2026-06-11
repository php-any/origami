<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\Middleware;
use Spring\Middleware\LogInterceptor;
use Spring\Service\QueryDemoService;

#[Middleware(LogInterceptor::class)]
#[Controller]
#[Route(prefix: "/api/queries")]
class QueryController {

    public function __construct(
        private QueryDemoService $queryService,
    ) {}

    /**
     * 单表查询：条件 + 排序 + 分页
     * GET /api/queries/users?min_age=25&limit=5
     */
    #[GetMapping(path: "/simple/users")]
    public function singleTableUsers($request, $response) {
        $minAge = 0;
        $limit = 10;
        $minAgeInput = $request->input('min_age');
        $limitInput = $request->input('limit');
        if ($minAgeInput !== null) {
            $minAge = (int)$minAgeInput;
        }
        if ($limitInput !== null) {
            $limit = (int)$limitInput;
        }
        if ($limit <= 0) {
            $limit = 10;
        }

        $rows = $this->queryService->singleTableQuery($minAge, $limit);
        $data = QueryDemoService::rowsToArray($rows);

        $response->header("Content-Type", "application/json; charset=utf-8");
        $response->json([
            "code" => 200,
            "message" => "单表查询：users WHERE age >= ? ORDER BY age DESC LIMIT ?",
            "data" => $data,
            "total" => count($data),
            "params" => ["min_age" => $minAge, "limit" => $limit],
        ]);
    }

    /**
     * 单表模糊搜索
     * GET /api/queries/users/search?keyword=张
     */
    #[GetMapping(path: "/users/search")]
    public function searchUsers($request, $response) {
        $keyword = $request->input('keyword') ?? '';
        $rows = $this->queryService->searchUsers($keyword);
        $data = QueryDemoService::rowsToArray($rows);

        $response->json([
            "code" => 200,
            "message" => "单表模糊查询：name/email LIKE",
            "data" => $data,
            "total" => count($data),
        ]);
    }

    /**
     * 聚合查询：按分类统计
     * GET /api/queries/products/stats
     */
    #[GetMapping(path: "/products/stats")]
    public function productStats($request, $response) {
        $rows = $this->queryService->aggregateByCategory();
        $data = QueryDemoService::rowsToArray($rows);

        $response->json([
            "code" => 200,
            "message" => "聚合查询：GROUP BY category",
            "data" => $data,
            "total" => count($data),
        ]);
    }

    /**
     * 子查询：高于品类均价的商品
     * GET /api/queries/products/above-avg
     */
    #[GetMapping(path: "/products/above-avg")]
    public function productsAboveAvg($request, $response) {
        $rows = $this->queryService->productsAboveCategoryAvg();
        $data = QueryDemoService::rowsToArray($rows);

        $response->json([
            "code" => 200,
            "message" => "子查询：price > 品类均价",
            "data" => $data,
            "total" => count($data),
        ]);
    }

    /**
     * INNER JOIN：订单 + 商品
     * GET /api/queries/orders/join-products
     */
    #[GetMapping(path: "/orders/join-products")]
    public function joinOrderProducts($request, $response) {
        $rows = $this->queryService->innerJoinOrderProducts();
        $data = QueryDemoService::rowsToArray($rows);

        $response->json([
            "code" => 200,
            "message" => "INNER JOIN：orders + products",
            "data" => $data,
            "total" => count($data),
        ]);
    }

    /**
     * LEFT JOIN：用户 + 订单
     * GET /api/queries/users/join-orders
     */
    #[GetMapping(path: "/users/join-orders")]
    public function joinUserOrders($request, $response) {
        $rows = $this->queryService->leftJoinUserOrders();
        $data = QueryDemoService::rowsToArray($rows);

        $response->json([
            "code" => 200,
            "message" => "LEFT JOIN：users + orders（含无订单用户）",
            "data" => $data,
            "total" => count($data),
        ]);
    }

    /**
     * 三表 JOIN：订单明细
     * GET /api/queries/orders/details
     */
    #[GetMapping(path: "/orders/details")]
    public function orderDetails($request, $response) {
        $rows = $this->queryService->orderDetails();
        $data = QueryDemoService::rowsToArray($rows);

        $response->json([
            "code" => 200,
            "message" => "三表 JOIN：users + orders + products",
            "data" => $data,
            "total" => count($data),
        ]);
    }

    /**
     * 条件 JOIN + GROUP BY：已完成订单消费统计
     * GET /api/queries/orders/completed-stats
     */
    #[GetMapping(path: "/orders/completed-stats")]
    public function completedOrderStats($request, $response) {
        $rows = $this->queryService->completedOrderStats();
        $data = QueryDemoService::rowsToArray($rows);

        $response->json([
            "code" => 200,
            "message" => "JOIN + GROUP BY：已完成订单消费汇总",
            "data" => $data,
            "total" => count($data),
        ]);
    }
}
