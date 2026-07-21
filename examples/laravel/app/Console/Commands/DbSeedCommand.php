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
        if (!$this->input->hasOption('force')) {
            if (!$this->prompt->confirm('This will seed the database. Continue?', true)) {
                $this->output->warning('Seeding cancelled.');
                return;
            }
        }

        app_make(DatabaseManager::class);

        $this->output->info('Seeding database...');
        DatabaseBootstrap::seed(true);
        $this->output->success('Database seeded.');
    }
}
