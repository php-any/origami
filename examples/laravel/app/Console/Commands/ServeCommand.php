<?php

namespace App\Console\Commands;

use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
use Bootstrap\Routing\Route;
use Cli\Annotation\Command as CommandAttribute;

#[CommandAttribute(name: 'serve', description: 'Serve the application on the PHP development server')]
class ServeCommand extends Command
{
    protected function defineInput(): InputDefinition
    {
        return (new InputDefinition())
            ->addArgument('port', false, 'The port to serve the application on', '8080')
            ->addOption('host', null, true, '0.0.0.0', 'The host address to serve the application on');
    }

    public function handle(): void
    {
        $port = (int) strval($this->argument('port') ?? '8080');
        $host = (string) ($this->option('host') ?? '0.0.0.0');

        if ($port < 1 || $port > 65535) {
            $this->error('Invalid port. Must be between 1 and 65535.');
            return;
        }

        require dirname(__DIR__, 3) . '/bootstrap/http.php';

        $server = bootstrap_http_server($port, $host);
        $routes = Route::getRoutes();

        $this->newLine();
        $this->info(sprintf('Server running on [http://%s:%d].', $host, $port));
        $this->comment('Press Ctrl+C to stop the server');
        $this->newLine();

        if (count($routes) > 0) {
            $this->comment('Registered routes (' . count($routes) . '):');
            foreach ($routes as $route) {
                $method = str_pad($route['method'], 7);
                $this->comment('  ' . $method . ' ' . $route['path']);
            }
            $this->newLine();
        }

        $server->run();
    }
}
