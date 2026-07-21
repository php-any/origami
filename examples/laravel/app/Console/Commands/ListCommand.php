<?php

namespace App\Console\Commands;

use Bootstrap\Console\Command;
use Cli\Annotation\Command as CommandAttribute;
use Cli\Annotation\CommandRegistry;

/**
 * 默认 Artisan 命令：无参数 / help 时由 main 分发到此（Symfony Console list 风格）
 */
#[CommandAttribute(name: 'list', description: 'List commands')]
class ListCommand extends Command
{

    public function handle(): void
    {
        $bin = 'laravel';

        $this->output->writeln(CommandRegistry::getLongVersion());
        $this->output->newLine();

        $this->output->writeln('Usage:');
        $this->output->writeln("  command [options] [arguments]");
        $this->output->newLine();

        $this->output->writeln('Options:');
        $this->renderOption('-h, --help', 'Display help for the given command. When no command is given display help for the list command');
        $this->renderOption('-q, --quiet', 'Do not output any message');
        $this->renderOption('-V, --version', 'Display this application version');
        $this->renderOption('    --ansi|--no-ansi', 'Force (or disable --no-ansi) ANSI output');
        $this->renderOption('-n, --no-interaction', 'Do not ask any interactive question');
        $this->renderOption('-v|vv|vvv, --verbose', 'Increase the verbosity of messages: 1 for normal output, 2 for more verbose output and 3 for debug');
        $this->output->newLine();

        $commands = CommandRegistry::getCommands();
        $global = [];
        $namespaces = [];

        foreach ($commands as $command) {
            $name = (string) ($command['name'] ?? '');
            $description = (string) ($command['description'] ?? '');
            if ($name === '') {
                continue;
            }
            if (!str_contains($name, ':')) {
                $global[$name] = $description;
                continue;
            }
            [$namespace] = explode(':', $name, 2);
            $namespaces[$namespace][$name] = $description;
        }

        ksort($global);
        ksort($namespaces);

        $maxNameLen = 0;
        foreach ($global as $name => $_) {
            $maxNameLen = max($maxNameLen, strlen($name));
        }
        foreach ($namespaces as $items) {
            foreach ($items as $name => $_) {
                $maxNameLen = max($maxNameLen, strlen($name));
            }
        }

        $this->output->writeln('Available commands:');
        foreach ($global as $name => $description) {
            $this->output->describeCommand($name, $description, $maxNameLen);
        }

        foreach ($namespaces as $namespace => $items) {
            ksort($items);
            $this->output->describeNamespace($namespace);
            foreach ($items as $name => $description) {
                $this->output->describeCommand($name, $description, $maxNameLen, 2);
            }
        }

        $this->output->newLine();
        $this->output->comment("Use \"{$bin} list\" to see all available commands.");
        $this->output->comment("Use \"{$bin} <command> --help\" for more information about a command.");
    }

    private function renderOption(string $option, string $description): void
    {
        $width = 24;
        $padding = max(1, $width - strlen($option) + 2);
        $this->output->writeln('  ' . $option . str_repeat(' ', $padding) . $description);
    }
}
