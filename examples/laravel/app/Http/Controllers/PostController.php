<?php

namespace App\Http\Controllers;

use App\Services\PostService;
use App\Services\UserService;
use Bootstrap\View\View;
use Net\Http\Response;

class PostController
{
    public function __construct(
        private PostService $postService,
        private UserService $userService,
    ) {}

    public function index(Response $response): void
    {
        $posts = $this->postService->all();

        View::render($response, 'posts.index', [
            'title' => '文章列表',
            'posts' => $posts,
        ]);
    }

    public function show(int $id, Response $response): void
    {
        $post = $this->postService->find($id);
        if ($post === null) {
            $response->error('文章不存在', 404);
            return;
        }

        $author = $this->userService->find($post->user_id);

        View::render($response, 'posts.show', [
            'title' => $post->title,
            'post' => $post,
            'author' => $author,
        ]);
    }
}
