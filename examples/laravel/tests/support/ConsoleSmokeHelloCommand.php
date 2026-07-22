<?php

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;

class ConsoleSmokeHelloCommand extends Command
{
    protected function configure()
    {
        $this->setName('smoke:hello');
        $this->addArgument('name', InputArgument::OPTIONAL, '', 'origami');
    }

    protected function execute($input, $output)
    {
        $output->writeln('hello:' . $input->getArgument('name'));
        return Command::SUCCESS;
    }
}
