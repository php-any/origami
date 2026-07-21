<?php

namespace App\Http\Controllers;

use App\Http\Requests\CreatePostRequest;
use App\Services\AuthService;
use App\Services\PostService;
use Net\Http\Request;
use Net\Http\Response;

class PostApiController
{
    public function __construct(
        private PostService $postService,
    ) {}

    public function index(Response $response): void
    {
        $list = $this->postService->allWithAuthors();
        $response->success(['list' => $list, 'total' => count($list)]);
    }

    public function show(int $id, Response $response): void
    {
        $item = $this->postService->findWithAuthor($id);
        if ($item === null) {
            $response->error('文章不存在', 404);
            return;
        }
        $response->success($item);
    }

    public function store(CreatePostRequest $request, Request $httpRequest, Response $response): void
    {
        $user = AuthService::userFromRequest($httpRequest);
        if ($user === null) {
            $response->error('Unauthenticated.', 401);
            return;
        }

        $item = $this->postService->create((int) $user['id'], $request->title, $request->body);
        $response->success($item, 'created', 201);
    }
}
