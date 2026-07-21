<?php

namespace App\Services;

use App\Models\Post;
use App\Models\User;
use Container\Annotation\Singleton;
use Database\DB;

#[Singleton]
class PostService
{
    public function __construct(
        private UserService $userService,
    ) {}

    private function query(): DB
    {
        return DB::model(Post::class);
    }

    public function all(): array
    {
        return $this->query()->orderBy('id DESC')->get();
    }

    public function find(int $id): ?Post
    {
        return $this->query()->where('id = ?', $id)->first();
    }

    public function create(int $userId, string $title, string $body): array
    {
        $post = new Post();
        $post->user_id = $userId;
        $post->title = $title;
        $post->body = $body;

        $result = DB::insert($post);
        $post->id = $result->insertId;

        $author = $this->userService->find($post->user_id);

        return $this->toArray($post, $author);
    }

    public function allWithAuthors(): array
    {
        $posts = $this->all();
        $userIds = array_map(fn (Post $post) => $post->user_id, $posts);
        $authors = $this->userService->findMany($userIds);

        $list = [];
        foreach ($posts as $post) {
            $author = $authors[$post->user_id] ?? null;
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

        $author = $this->userService->find($post->user_id);

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
