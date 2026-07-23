<?php

namespace App\Console\Commands;

use Bootstrap\Console\Command;
use Cli\Annotation\Command as CommandAttribute;

#[CommandAttribute(name: 'inspire', description: 'Display an inspiring quote')]
class InspireCommand extends Command
{
    /**
     * Laravel 风格 signature（与 #[Command] 双轨；Origami 仍靠注解发现）。
     *
     * @var string
     */
    protected $signature = 'inspire {--quiet : Minimal output}';

    /**
     * @var string
     */
    protected $description = 'Display an inspiring quote';

    private array $quotes = [
        'Simplicity is the ultimate sophistication.',
        'Make it work, make it right, make it fast.',
        '代码如诗，架构如律。',
        'The best way to predict the future is to invent it.',
    ];

    public function handle(): void
    {
        $quote = $this->quotes[array_rand($this->quotes)];

        if ($this->optionEnabled('quiet')) {
            $this->line($quote);
            return;
        }

        $this->newLine();
        $this->line('  "' . $quote . '"');
        $this->newLine();
    }
}
