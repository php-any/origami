<?php

namespace App\Console\Commands;

use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
use App\Models\User;
use App\Services\DatabaseManager;
use Cli\Annotation\Command as CommandAttribute;
use Database\DB;

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
        $name = $this->input->getArgument('name');
        $email = $this->input->getArgument('email');
        $password = $this->input->getArgument('password') ?? 'secret123';

        $interactive = $this->input->hasOption('interactive') || $name === null || $email === null;

        if ($interactive) {
            $this->output->section('Create User (interactive)');
            $name = $this->prompt->ask('Name', $name);
            $email = $this->prompt->ask('Email', $email);

            if (!$this->prompt->confirm("Create user {$name} <{$email}>?", true)) {
                $this->output->warning('Cancelled.');
                return;
            }

            $password = $this->prompt->secret('Password (leave empty for default)');
            if ($password === '') {
                $password = 'secret123';
            }
        }

        if ($name === null || $email === null || $name === '' || $email === '') {
            $this->output->error('Usage: laravel make:user <name> <email> [password]');
            $this->output->comment('Or use: laravel make:user --interactive');
            return;
        }

        app_make(DatabaseManager::class);

        $exists = DB::model(User::class)->where('email = ?', $email)->first();
        if ($exists !== null) {
            $this->output->error("User already exists: {$email}");
            return;
        }

        $user = new User();
        $user->name = $name;
        $user->email = $email;
        $user->password = \App\Services\AuthService::hashPassword($password);
        DB::insert($user);

        $this->output->success(sprintf('Created user #%s %s <%s>', $user->id, $user->name, $user->email));
    }
}
