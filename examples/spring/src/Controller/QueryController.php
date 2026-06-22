<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\Middleware;
use Net\Http\Request;
use Net\Http\Response;
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
    public function singleTableUsers(Request $request, Response $response): void {
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

        $response->success([
            'rows' => $data,
            'total' => count($data),
            'params' => ['min_age' => $minAge, 'limit' => $limit],
        ], '单表查询：users WHERE age >= ? ORDER BY age DESC LIMIT ?');
    }

    /**
     * 单表模糊搜索
     * GET /api/queries/users/search?keyword=张
     */
    #[GetMapping(path: "/users/search")]
    public function searchUsers(Request $request, Response $response): void {
        $keyword = $request->input('keyword') ?? '';
        $rows = $this->queryService->searchUsers($keyword);
        $data = QueryDemoService::rowsToArray($rows);

        $response->success([
            'rows' => $data,
            'total' => count($data),
        ], '单表模糊查询：name/email LIKE');
    }

    /**
     * 聚合查询：按分类统计
     * GET /api/queries/products/stats
     */
    #[GetMapping(path: "/products/stats")]
    public function productStats(Request $request, Response $response): void {
        $rows = $this->queryService->aggregateByCategory();
        $data = QueryDemoService::rowsToArray($rows);

        $response->success([
            'rows' => $data,
            'total' => count($data),
        ], '聚合查询：GROUP BY category');
    }

    /**
     * 子查询：高于品类均价的商品
     * GET /api/queries/products/above-avg
     */
    #[GetMapping(path: "/products/above-avg")]
    public function productsAboveAvg(Request $request, Response $response): void {
        $rows = $this->queryService->productsAboveCategoryAvg();
        $data = QueryDemoService::rowsToArray($rows);

        $response->success([
            'rows' => $data,
            'total' => count($data),
        ], '子查询：price > 品类均价');
    }

    /**
     * INNER JOIN：订单 + 商品
     * GET /api/queries/orders/join-products
     */
    #[GetMapping(path: "/orders/join-products")]
    public function joinOrderProducts(Request $request, Response $response): void {
        $rows = $this->queryService->innerJoinOrderProducts();
        $data = QueryDemoService::rowsToArray($rows);

        $response->success([
            'rows' => $data,
            'total' => count($data),
        ], 'INNER JOIN：orders + products');
    }

    /**
     * LEFT JOIN：用户 + 订单
     * GET /api/queries/users/join-orders
     */
    #[GetMapping(path: "/users/join-orders")]
    public function joinUserOrders(Request $request, Response $response): void {
        $rows = $this->queryService->leftJoinUserOrders();
        $data = QueryDemoService::rowsToArray($rows);

        $response->success([
            'rows' => $data,
            'total' => count($data),
        ], 'LEFT JOIN：users + orders（含无订单用户）');
    }

    /**
     * 三表 JOIN：订单明细
     * GET /api/queries/orders/details
     */
    #[GetMapping(path: "/orders/details")]
    public function orderDetails(Request $request, Response $response): void {
        $rows = $this->queryService->orderDetails();
        $data = QueryDemoService::rowsToArray($rows);

        $response->success([
            'rows' => $data,
            'total' => count($data),
        ], '三表 JOIN：users + orders + products');
    }

    /**
     * 条件 JOIN + GROUP BY：已完成订单消费统计
     * GET /api/queries/orders/completed-stats
     */
    #[GetMapping(path: "/orders/completed-stats")]
    public function completedOrderStats(Request $request, Response $response): void {
        $rows = $this->queryService->completedOrderStats();
        $data = QueryDemoService::rowsToArray($rows);

        $response->success([
            'rows' => $data,
            'total' => count($data),
        ], 'JOIN + GROUP BY：已完成订单消费汇总');
    }
}
