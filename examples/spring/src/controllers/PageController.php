<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\GetMapping;
use Net\Annotation\Route;

/**
 * 静态 HTML 页面（原 html 示例，使用 view 渲染模板）
 */
#[Controller]
#[Route(prefix: "/")]
class PageController {

    private function render($response, string $page): void {
        $pagesDir = dirname(dirname(__DIR__)) . '/pages/';
        $response->view($pagesDir . $page . '.html', []);
    }

    #[GetMapping(path: "/")]
    public function home($request, $response) {
        $this->render($response, "index");
    }

    #[GetMapping(path: "/index")]
    public function index($request, $response) {
        $this->render($response, "index");
    }

    #[GetMapping(path: "/about")]
    public function about($request, $response) {
        $this->render($response, "about");
    }

    #[GetMapping(path: "/products")]
    public function products($request, $response) {
        $this->render($response, "products");
    }

    #[GetMapping(path: "/contact")]
    public function contact($request, $response) {
        $this->render($response, "contact");
    }
}
