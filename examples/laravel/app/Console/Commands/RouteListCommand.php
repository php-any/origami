<?php

namespace App\Console\Commands;

use App\Application;
use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
use Bootstrap\Routing\Route;
use Cli\Annotation\Command as CommandAttribute;
use Net\Http\Server;

#[CommandAttribute(name: 'route:list', description: 'List all registered routes')]
class RouteListCommand extends Command
{
    protected function defineInput(): InputDefinition
    {
        return (new InputDefinition())
            ->addOption('method', 'm', true, null, 'Filter by HTTP method');
    }

    public function handle(): void
    {
        $server = new Server('127.0.0.1', port: 8080);
        $server->boot(Application::class);

        $routes = Route::getRoutes();
        $filter = strtoupper((string) ($this->option('method') ?? ''));

        $this->output->title('Registered Routes');

        $rows = [];
        $count = 0;
        foreach ($routes as $route) {
            if ($filter !== '' && strtoupper($route['method']) !== $filter) {
                continue;
            }

            $action = $route['controller'] . '@' . $route['action'];
            $rows[] = [$route['method'], $route['path'], $action];
            $count++;
        }

        $this->table(['Method', 'URI', 'Action'], $rows);
        $this->comment("Total: {$count} routes");
        $this->newLine();
    }
}
