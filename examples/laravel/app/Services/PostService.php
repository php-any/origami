<?php

namespace App\Services;

use App\Models\Post;
use App\Models\User;
use Container\Annotation\Singleton;

#[Singleton]
class PostService
{
    public function __construct(
        private UserService $userService,
    ) {}

    public function all(): array
    {
        return Post::query()->orderByDesc('id')->get()->all();
    }

    public function find(int $id): ?Post
    {
        return Post::query()->find($id);
    }

    public function create(int $userId, string $title, string $body): array
    {
        $post = Post::query()->create([
            'user_id' => $userId,
            'title' => $title,
            'body' => $body,
            'created_at' => date('Y-m-d H:i:s'),
        ]);

        $author = $this->userService->find((int) $post->user_id);

        return $this->toArray($post, $author);
    }

    public function allWithAuthors(): array
    {
        $posts = $this->all();
        $userIds = array_map(static fn (Post $post) => (int) $post->user_id, $posts);
        $authors = $this->userService->findMany($userIds);

        $list = [];
        foreach ($posts as $post) {
            $author = $authors[(int) $post->user_id] ?? null;
            $list[] = $this->toArray($post, $author);
        }

        return $list;
    }

    public function findWithAuthor(int $id): ?array
    {
        $post = $this->find($id);
        if ($post === null) {
            return null;
        }

        $author = $this->userService->find((int) $post->user_id);

        return $this->toArray($post, $author);
    }

    public function toArray(Post $post, ?User $author = null): array
    {
        $data = [
            'id' => $post->id,
            'user_id' => $post->user_id,
            'title' => $post->title,
            'body' => $post->body,
            'created_at' => $post->created_at,
        ];
        if ($author !== null) {
            $data['author'] = $author->name;
        }

        return $data;
    }
}
