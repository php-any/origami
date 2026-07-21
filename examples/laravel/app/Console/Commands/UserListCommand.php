<?php

namespace App\Console\Commands;

use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
use App\Services\DatabaseManager;
use App\Services\UserService;
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
        $limit = $this->input->getOption('limit');
        if ($limit !== null && $limit !== '') {
            $users = array_slice($users, 0, (int) $limit);
        }

        $this->output->title('Users');

        if (count($users) === 0) {
            $this->output->warning('No users found — try: laravel migrate 或 laravel make:user');
            return;
        }

        $table = $this->table();
        $table->setHeaders(['ID', 'Name', 'Email', 'Created']);
        foreach ($users as $user) {
            $table->addRow([
                (string) $user->id,
                $user->name,
                $user->email,
                $user->created_at ?? '-',
            ]);
        }
        $table->render();

        $this->output->comment('Total: ' . count($users));
        $this->output->newLine();
    }
}
