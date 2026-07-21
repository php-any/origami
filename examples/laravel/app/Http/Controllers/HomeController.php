<?php

namespace App\Http\Controllers;

use App\Services\PostService;
use Bootstrap\View\View;
use Net\Http\Response;

class HomeController
{
    public function __construct(
        private PostService $postService,
    ) {}

    public function index(Response $response): void
    {
        $posts = $this->postService->all();
        $recent = array_slice($posts, 0, 3);

        View::render($response, 'home', [
            'title' => '首页',
            'posts' => $recent,
        ]);
    }
}
