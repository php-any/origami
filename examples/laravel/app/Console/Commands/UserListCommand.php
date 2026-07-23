<?php

namespace App\Console\Commands;

use App\Services\DatabaseManager;
use App\Services\UserService;
use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
use Cli\Annotation\Command as CommandAttribute;

#[CommandAttribute(name: 'user:list', description: 'List all users')]
class UserListCommand extends Command
{
    protected function defineInput(): InputDefinition
    {
        return (new InputDefinition())
            ->addOption('limit', 'l', true, null, 'Maximum number of rows');
    }

    public function handle(): void
    {
        app_make(DatabaseManager::class);

        $users = app_make(UserService::class)->all();
        $limit = $this->option('limit');
        if ($limit !== null && $limit !== '') {
            $users = array_slice($users, 0, (int) $limit);
        }

        $this->output->title('Users');

        if (count($users) === 0) {
            $this->warn('No users found — try: laravel migrate 或 laravel make:user');
            return;
        }

        $rows = [];
        foreach ($users as $user) {
            $rows[] = [
                (string) $user->id,
                $user->name,
                $user->email,
                $user->created_at ?? '-',
            ];
        }
        $this->table(['ID', 'Name', 'Email', 'Created'], $rows);

        $this->comment('Total: ' . count($users));
        $this->newLine();
    }
}
