<?php

namespace App\Console\Commands;

use App\Models\Post;
use App\Services\DatabaseManager;
use Bootstrap\Console\Command;
use Cli\Annotation\Command as CommandAttribute;

#[CommandAttribute(name: 'post:list', description: 'List all posts')]
class PostListCommand extends Command
{
    public function handle(): void
    {
        app_make(DatabaseManager::class)->connect();

        $posts = Post::query()->orderBy('id')->get();

        $this->output->title('Posts');

        if (count($posts) === 0) {
            $this->warn('No posts found.');
            return;
        }

        $rows = [];
        foreach ($posts as $post) {
            $rows[] = [(string) $post->id, (string) $post->user_id, $post->title];
        }
        $this->table(['ID', 'User', 'Title'], $rows);

        $this->comment('Total: ' . count($posts));
        $this->newLine();
    }
}
