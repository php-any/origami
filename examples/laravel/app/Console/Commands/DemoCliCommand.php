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
        if ($this->input->hasOption('help')) {
            $this->showUsage();
            return;
        }

        $name = (string) $this->input->getArgument('name');

        $this->output->title('CLI Demo');
        $this->output->info("Hello, {$name}!");
        $this->output->success('Formatted success message');
        $this->output->warning('Formatted warning message');
        $this->output->error('Formatted error message (stderr)');

        if ($this->input->hasOption('verbose')) {
            $this->output->section('Parsed Input');
            $this->output->comment('Arguments: ' . json_encode($this->input->getArguments(), JSON_UNESCAPED_UNICODE));
            $this->output->comment('Options: ' . json_encode($this->input->getOptions(), JSON_UNESCAPED_UNICODE));
        }

        $this->output->section('Table Output');
        $limit = max(1, (int) $this->input->getOption('limit'));
        $table = $this->table();
        $table->setHeaders(['#', 'Feature', 'Class']);
        $features = [
            ['Arguments / Options', 'Bootstrap\\Console\\ArgvInput'],
            ['Styled Output', 'Bootstrap\\Console\\Output'],
            ['Table', 'Bootstrap\\Console\\Table'],
            ['Progress Bar', 'Bootstrap\\Console\\ProgressBar'],
            ['User Input', 'Bootstrap\\Console\\Prompt'],
        ];
        foreach (array_slice($features, 0, $limit) as $index => $row) {
            $table->addRow([(string) ($index + 1), $row[0], $row[1]]);
        }
        $table->render();

        $this->output->section('Progress Bar');
        $bar = $this->createProgressBar(20);
        $bar->start();
        for ($i = 0; $i < 20; $i++) {
            $bar->advance();
        }
        $bar->finish();

        if ($this->input->hasOption('interactive')) {
            $this->output->section('Interactive Prompts');
            $answer = $this->prompt->ask('Your name', $name);
            $confirmed = $this->prompt->confirm('Continue demo?', true);
            $choice = $this->prompt->choice('Pick a color', ['red', 'green', 'blue'], 'green');

            $this->output->info("ask => {$answer}");
            $this->output->info('confirm => ' . ($confirmed ? 'yes' : 'no'));
            $this->output->info("choice => {$choice}");
        } else {
            $this->output->comment('Add --interactive to try ask / confirm / choice prompts.');
        }

        $this->output->newLine();
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
