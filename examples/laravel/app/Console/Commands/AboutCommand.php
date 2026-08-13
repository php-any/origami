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
                ['Debug Mode', $debug ? 'ENABLED' : 'OFF'],
                ['URL', $url],
                ['Timezone', (string) config('app.timezone', '')],
                ['Locale', (string) config('app.locale', '')],
            ],
            'Cache' => [
                ['Config', 'NOT CACHED'],
                ['Events', 'NOT CACHED'],
                ['Routes', 'NOT CACHED'],
                ['Views', 'NOT CACHED'],
            ],
            'Drivers' => [
                ['Database', $dbDefault],
                ['Database Path', $dbPath],
            ],
        ];

        if ($this->optionEnabled('json')) {
            // 用 list-of-pairs，避免 Origami 嵌套关联数组丢键
            $payload = [];
            foreach ($sections as $section => $rows) {
                $items = [];
                foreach ($rows as $pair) {
                    $items[] = [(string) $pair[0], (string) $pair[1]];
                }
                $payload[] = [(string) $section, $items];
            }
            $this->line(json_encode($payload, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
            return;
        }

        $first = true;
        foreach ($sections as $section => $rows) {
            if (!$first) {
                $this->newLine();
            }
            $first = false;

            $this->twoColumnDetail(Style::green(Style::bold($section)));
            foreach ($rows as $pair) {
                $label = (string) $pair[0];
                $value = (string) $pair[1];
                if ($section === 'Environment' && $label === 'Debug Mode' && $value === 'ENABLED') {
                    $value = Style::yellow(Style::bold('ENABLED'));
                }
                if ($section === 'Cache') {
                    $value = Style::yellow(Style::bold($value));
                }
                $this->twoColumnDetail($label, $value);
            }
        }

        $this->newLine();
    }
}
