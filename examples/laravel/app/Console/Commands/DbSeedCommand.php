<?php

namespace App\Console\Commands;

use App\Database\DatabaseBootstrap;
use App\Services\DatabaseManager;
use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
use Cli\Annotation\Command as CommandAttribute;

#[CommandAttribute(name: 'db:seed', description: 'Seed the database with records')]
class DbSeedCommand extends Command
{
    protected function defineInput(): InputDefinition
    {
        return (new InputDefinition())
            ->addOption('force', 'f', false, null, 'Run without confirmation');
    }

    public function handle(): void
    {
        if (!$this->optionEnabled('force')) {
            if (!$this->confirm('This will seed the database. Continue?', true)) {
                $this->warn('Seeding cancelled.');
                return;
            }
        }

        app_make(DatabaseManager::class);

        $this->info('Seeding database...');
        DatabaseBootstrap::seed(true);
        $this->info('Database seeded.');
    }
}
