<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;
use Net\Annotation\PutMapping;
use Net\Annotation\DeleteMapping;
use Net\Annotation\Middleware;
use Net\Http\Response;
use Spring\DTO\Request\CreateProductRequest;
use Spring\DTO\Request\SearchProductsQuery;
use Spring\DTO\Request\UpdateProductRequest;
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
    public function listProducts(Response $response): void {
        $products = $this->productService->findAll();

        $response->success([
            'list' => $products,
            'total' => count($products),
        ]);
    }

    #[GetMapping(path: "/product/{id}")]
    public function getProduct(int $id, Response $response): void {
        $product = $this->productService->findById($id);

        if (!$product) {
            $response->error('商品不存在', 404);
            return;
        }

        $response->success($product);
    }

    #[PostMapping(path: "/products")]
    public function createProduct(CreateProductRequest $request, Response $response): void {
        $product = $this->productService->create([
            'name' => $request->name,
            'price' => $request->price,
            'category' => $request->category,
            'description' => $request->description,
        ]);

        $response->success($product, 'created', 201);
    }

    #[PutMapping(path: "/product/{id}")]
    public function updateProduct(int $id, UpdateProductRequest $request, Response $response): void {
        $product = $this->productService->update($id, $request->toUpdateArray());

        if (!$product) {
            $response->error('商品不存在', 404);
            return;
        }

        $response->success($product, 'updated');
    }

    #[DeleteMapping(path: "/product/{id}")]
    public function deleteProduct(int $id, Response $response): void {
        $result = $this->productService->delete($id);

        if (!$result) {
            $response->error('商品不存在', 404);
            return;
        }

        $response->success(null, 'deleted');
    }

    #[GetMapping(path: "/products/search")]
    public function searchProducts(SearchProductsQuery $query, Response $response): void {
        $products = $this->productService->search($query->keyword, $query->category);

        $response->success([
            'list' => $products,
            'total' => count($products),
        ]);
    }
}
