<?php

namespace App\Console\Commands;

use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
use Bootstrap\Console\Style;
use Cli\Annotation\Command as CommandAttribute;
use Cli\Annotation\CommandRegistry;

#[CommandAttribute(name: 'about', description: 'Display basic information about your application')]
class AboutCommand extends Command
{
    protected function defineInput(): InputDefinition
    {
        return (new InputDefinition())
            ->addOption('json', 'j', false, null, 'Output as JSON');
    }

    public function handle(): void
    {
        $dbDefault = config('database.default', 'sqlite');
        $dbPath = config('database.connections.' . $dbDefault . '.database', '(unset)');

        $url = (string) config('app.url', '');
        $url = preg_replace('#^https?://#', '', $url) ?? $url;

        $debug = !empty(config('app.debug'));
        $sections = [
            'Environment' => [
                ['Application Name', (string) config('app.name', 'LaravelDemo')],
                ['Laravel Version', CommandRegistry::getAppVersion()],
                ['PHP Version', PHP_VERSION],
                ['Environment', (string) config('app.env', 'local')],
                ['Debug Mode', $debug ? Style::yellow(Style::bold('ENABLED')) : 'OFF'],
                ['URL', $url],
                ['Timezone', (string) config('app.timezone', '')],
                ['Locale', (string) config('app.locale', '')],
            ],
            'Cache' => [
                ['Config', Style::yellow(Style::bold('NOT CACHED'))],
                ['Events', Style::yellow(Style::bold('NOT CACHED'))],
                ['Routes', Style::yellow(Style::bold('NOT CACHED'))],
                ['Views', Style::yellow(Style::bold('NOT CACHED'))],
            ],
            'Drivers' => [
                ['Database', $dbDefault],
                ['Database Path', $dbPath],
            ],
        ];

        if ($this->input->hasOption('json')) {
            $data = [];
            foreach ($sections as $section => $rows) {
                $key = strtolower(preg_replace('/\s+/', '_', $section) ?? $section);
                $data[$key] = [];
                foreach ($rows as [$label, $value]) {
                    $data[$key][strtolower(str_replace(' ', '_', $label))] = $value;
                }
            }
            $this->output->writeln(json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
            return;
        }

        $first = true;
        foreach ($sections as $section => $rows) {
            if (!$first) {
                $this->output->newLine();
            }
            $first = false;

            $this->output->twoColumnDetail(Style::green(Style::bold($section)));
            foreach ($rows as [$label, $value]) {
                $this->output->twoColumnDetail($label, $value);
            }
        }

        $this->output->newLine();
    }
}
