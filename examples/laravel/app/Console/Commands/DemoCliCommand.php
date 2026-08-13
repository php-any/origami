<?php

namespace App\Console\Commands;

use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
use Cli\Annotation\Command as CommandAttribute;

/**
 * CLI 能力演示命令：参数/选项、格式化输出、表格、进度条、交互输入
 */
#[CommandAttribute(name: 'demo:cli', description: 'Demonstrate CLI helpers (args, table, progress, prompt)')]
class DemoCliCommand extends Command
{
    protected function defineInput(): InputDefinition
    {
        return (new InputDefinition())
            ->addArgument('name', false, 'Greeting target', 'World')
            ->addOption('verbose', 'v', false, null, 'Show parsed input')
            ->addOption('limit', 'l', true, '3', 'Table row count')
            ->addOption('interactive', 'i', false, null, 'Run prompt examples')
            ->addOption('help', 'h', false, null, 'Show usage');
    }

    public function handle(): void
    {
        if ($this->optionEnabled('help')) {
            $this->showUsage();
            return;
        }

        $name = (string) ($this->argument('name') ?? 'World');

        $this->output->title('CLI Demo');
        $this->info("Hello, {$name}!");
        $this->output->success('Formatted success message');
        $this->warn('Formatted warning message');
        $this->error('Formatted error message (stderr)');

        if ($this->optionEnabled('verbose')) {
            $this->output->section('Parsed Input');
            $this->comment('Arguments: ' . json_encode($this->arguments(), JSON_UNESCAPED_UNICODE));
            $this->comment('Options: ' . json_encode($this->options(), JSON_UNESCAPED_UNICODE));
        }

        $this->output->section('Table Output');
        $limit = max(1, (int) ($this->option('limit') ?? 3));
        $features = [
            ['Arguments / Options', 'Symfony ArgvInput via illuminate/console'],
            ['Styled Output', 'Illuminate\\Console\\OutputStyle'],
            ['Table', 'Symfony\\Component\\Console\\Helper\\Table'],
            ['Progress Bar', 'SymfonyStyle::createProgressBar'],
            ['User Input', 'Command::ask / confirm / choice'],
        ];
        $rows = [];
        foreach (array_slice($features, 0, $limit) as $index => $row) {
            $rows[] = [(string) ($index + 1), $row[0], $row[1]];
        }
        $this->table(['#', 'Feature', 'Class'], $rows);

        $this->output->section('Progress Bar');
        $bar = $this->output->createProgressBar(20);
        $bar->start();
        for ($i = 0; $i < 20; $i++) {
            $bar->advance();
        }
        $bar->finish();
        $this->newLine();

        if ($this->optionEnabled('interactive')) {
            $this->output->section('Interactive Prompts');
            $answer = $this->ask('Your name', $name);
            $confirmed = $this->confirm('Continue demo?', true);
            $choice = $this->choice('Pick a color', ['red', 'green', 'blue'], 'green');

            $this->info("ask => {$answer}");
            $this->info('confirm => ' . ($confirmed ? 'yes' : 'no'));
            $this->info("choice => {$choice}");
        } else {
            $this->comment('Add --interactive to try ask / confirm / choice prompts.');
        }

        $this->newLine();
        $this->output->success('CLI demo finished.');
    }

    private function showUsage(): void
    {
        $this->output->title('demo:cli Usage');
        $this->output->listing([
            'laravel demo:cli',
            'laravel demo:cli Alice',
            'laravel demo:cli Alice --verbose',
            'laravel demo:cli --limit=5',
            'laravel demo:cli --interactive',
        ]);
    }
}
