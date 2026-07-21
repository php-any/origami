<?php

namespace App\Console\Commands;

use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
use Cli\Annotation\Command as CommandAttribute;

#[CommandAttribute(name: 'inspire', description: 'Display an inspiring quote')]
class InspireCommand extends Command
{
    private array $quotes = [
        'Simplicity is the ultimate sophistication.',
        'Make it work, make it right, make it fast.',
        '代码如诗，架构如律。',
        'The best way to predict the future is to invent it.',
    ];

    protected function defineInput(): InputDefinition
    {
        return (new InputDefinition())
            ->addOption('quiet', 'q', false, null, 'Minimal output');
    }

    public function handle(): void
    {
        $quote = $this->quotes[array_rand($this->quotes)];

        if ($this->input->hasOption('quiet')) {
            $this->output->writeln($quote);
            return;
        }

        $this->output->newLine();
        $this->output->writeln('  "' . $quote . '"');
        $this->output->newLine();
    }
}
