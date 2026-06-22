<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\GetMapping;
use Net\Annotation\Route;
use Net\Http\Request;
use Net\Http\Response;

/**
 * 静态 HTML 页面（原 html 示例，使用 view 渲染模板）
 */
#[Controller]
#[Route(prefix: "/")]
class PageController {

    private function render(Response $response, string $page): void {
        $pagesDir = dirname(dirname(__DIR__)) . '/pages/';
        $response->view($pagesDir . $page . '.html', []);
    }

    #[GetMapping(path: "/")]
    public function home(Request $request, Response $response): void {
        $this->render($response, "index");
    }

    #[GetMapping(path: "/index")]
    public function index(Request $request, Response $response): void {
        $this->render($response, "index");
    }

    #[GetMapping(path: "/about")]
    public function about(Request $request, Response $response): void {
        $this->render($response, "about");
    }

    #[GetMapping(path: "/products")]
    public function products(Request $request, Response $response): void {
        $this->render($response, "products");
    }

    #[GetMapping(path: "/contact")]
    public function contact(Request $request, Response $response): void {
        $this->render($response, "contact");
    }

    #[GetMapping(path: "/chat")]
    public function chat(Request $request, Response $response): void {
        $this->render($response, "chat");
    }
}
