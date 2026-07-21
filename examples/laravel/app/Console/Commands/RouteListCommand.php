<?php

namespace App\Console\Commands;

use App\Application;
use Bootstrap\Routing\Route;
use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
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
        $filter = strtoupper((string) ($this->input->getOption('method') ?? ''));

        $this->output->title('Registered Routes');

        $table = $this->table();
        $table->setHeaders(['Method', 'URI', 'Action']);

        $count = 0;
        foreach ($routes as $route) {
            if ($filter !== '' && strtoupper($route['method']) !== $filter) {
                continue;
            }

            $action = $route['controller'] . '@' . $route['action'];
            $table->addRow([$route['method'], $route['path'], $action]);
            $count++;
        }

        $table->render();
        $this->output->comment("Total: {$count} routes");
        $this->output->newLine();
    }
}
