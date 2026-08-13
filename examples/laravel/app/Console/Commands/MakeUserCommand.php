<?php

namespace App\Console\Commands;

use App\Models\User;
use App\Services\AuthService;
use App\Services\DatabaseManager;
use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
use Cli\Annotation\Command as CommandAttribute;

#[CommandAttribute(name: 'make:user', description: 'Create a new user (usage: make:user name email [password])')]
class MakeUserCommand extends Command
{
    protected function defineInput(): InputDefinition
    {
        return (new InputDefinition())
            ->addArgument('name', false, 'User display name')
            ->addArgument('email', false, 'User email address')
            ->addArgument('password', false, 'Password (default: secret123)')
            ->addOption('interactive', 'i', false, null, 'Prompt for missing fields');
    }

    public function handle(): void
    {
        $name = $this->argument('name');
        $email = $this->argument('email');
        $password = $this->argument('password') ?? 'secret123';

        $interactive = $this->optionEnabled('interactive') || $name === null || $email === null;

        if ($interactive) {
            $this->output->section('Create User (interactive)');
            $name = $this->ask('Name', $name);
            $email = $this->ask('Email', $email);

            if (!$this->confirm("Create user {$name} <{$email}>?", true)) {
                $this->warn('Cancelled.');
                return;
            }

            $password = $this->secret('Password (leave empty for default)');
            if ($password === '') {
                $password = 'secret123';
            }
        }

        if ($name === null || $email === null || $name === '' || $email === '') {
            $this->error('Usage: laravel make:user <name> <email> [password]');
            $this->comment('Or use: laravel make:user --interactive');
            return;
        }

        app_make(DatabaseManager::class)->connect();

        if (User::query()->where('email', $email)->exists()) {
            $this->error("User already exists: {$email}");
            return;
        }

        $user = User::query()->create([
            'name' => $name,
            'email' => $email,
            'password' => AuthService::hashPassword($password),
            'created_at' => date('Y-m-d H:i:s'),
        ]);

        $this->info(sprintf('Created user #%s %s <%s>', $user->id, $user->name, $user->email));
    }
}
