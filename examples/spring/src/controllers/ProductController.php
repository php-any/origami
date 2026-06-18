<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;
use Net\Annotation\PutMapping;
use Net\Annotation\DeleteMapping;
use Net\Annotation\Middleware;
use Net\Http\Request;
use Net\Http\Response;
use Spring\Service\ProductService;
use Spring\Middleware\AuthInterceptor;
use Spring\Middleware\LogInterceptor;

#[Middleware(AuthInterceptor::class)]
#[Middleware(LogInterceptor::class)]
#[Controller]
#[Route(prefix: "/api")]
class ProductController {

    public function __construct(
        private ProductService $productService,
    ) {}

    #[GetMapping(path: "/products")]
    public function listProducts(Request $request, Response $response): void {
        $products = $this->productService->findAll();

        $response->header("Content-Type", "application/json; charset=utf-8");
        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => $products,
            "total" => count($products)
        ]);
    }

    #[GetMapping(path: "/product/{id}")]
    public function getProduct(Request $request, Response $response): void {
        $id = (int)$request->pathValue('id');
        $product = $this->productService->findById($id);

        if (!$product) {
            $response->status(404)->json([
                "code" => 404,
                "message" => "商品不存在",
                "data" => null
            ]);
            return;
        }

        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => $product
        ]);
    }

    #[PostMapping(path: "/products")]
    public function createProduct(Request $request, Response $response): void {
        $body = $request->body();

        if (!isset($body['name']) || !isset($body['price'])) {
            $response->status(400)->json([
                "code" => 400,
                "message" => "缺少必要参数：name 和 price",
                "data" => null
            ]);
            return;
        }

        $product = $this->productService->create($body);

        $response->status(201)->json([
            "code" => 201,
            "message" => "created",
            "data" => $product
        ]);
    }

    #[PutMapping(path: "/product/{id}")]
    public function updateProduct(Request $request, Response $response): void {
        $id = (int)$request->pathValue('id');
        $body = $request->body();

        $product = $this->productService->update($id, $body);

        if (!$product) {
            $response->status(404)->json([
                "code" => 404,
                "message" => "商品不存在",
                "data" => null
            ]);
            return;
        }

        $response->json([
            "code" => 200,
            "message" => "updated",
            "data" => $product
        ]);
    }

    #[DeleteMapping(path: "/product/{id}")]
    public function deleteProduct(Request $request, Response $response): void {
        $id = (int)$request->pathValue('id');
        $result = $this->productService->delete($id);

        if (!$result) {
            $response->status(404)->json([
                "code" => 404,
                "message" => "商品不存在",
                "data" => null
            ]);
            return;
        }

        $response->json([
            "code" => 200,
            "message" => "deleted",
            "data" => null
        ]);
    }

    #[GetMapping(path: "/products/search")]
    public function searchProducts(Request $request, Response $response): void {
        $keyword = $request->input('keyword') ?? '';
        $category = $request->input('category') ?? '';

        $products = $this->productService->search($keyword, $category);

        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => $products,
            "total" => count($products)
        ]);
    }
}
