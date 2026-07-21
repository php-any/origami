<?php

namespace App\Console\Commands;

use Bootstrap\Console\Command;
use App\Models\Post;
use App\Services\DatabaseManager;
use Cli\Annotation\Command as CommandAttribute;
use Database\DB;

#[CommandAttribute(name: 'post:list', description: 'List all posts')]
class PostListCommand extends Command
{
    public function handle(): void
    {
        app_make(DatabaseManager::class);

        $posts = DB::model(Post::class)->orderBy('id ASC')->get();

        $this->output->title('Posts');

        if (count($posts) === 0) {
            $this->output->warning('No posts found.');
            return;
        }

        $table = $this->table();
        $table->setHeaders(['ID', 'User', 'Title']);
        foreach ($posts as $post) {
            $table->addRow([(string) $post->id, (string) $post->user_id, $post->title]);
        }
        $table->render();

        $this->output->comment('Total: ' . count($posts));
        $this->output->newLine();
    }
}
