<?php

namespace App\Console\Commands;

use App\Database\DatabaseBootstrap;
use Bootstrap\Console\Command;
use Cli\Annotation\Command as CommandAttribute;

#[CommandAttribute(name: 'migrate', description: 'Run database migrations and seeders')]
class MigrateCommand extends Command
{
    public function handle(): void
    {
        $this->info('Running migrations...');
        DatabaseBootstrap::migrateAndSeed(null, true);
        $this->info('Migrations completed.');
    }
}
