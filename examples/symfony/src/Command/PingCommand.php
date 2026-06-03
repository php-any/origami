<?php

namespace App\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class PingCommand extends Command
{
    protected function configure(): void
    {
        $this->setName('app:ping');
        $this->setDescription('Ping command');
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $output->writeln('pong');
        return Command::SUCCESS;
    }
}
