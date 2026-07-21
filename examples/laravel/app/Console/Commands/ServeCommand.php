<?php

namespace App\Console\Commands;

use Bootstrap\Routing\Route;
use Bootstrap\Console\Command;
use Bootstrap\Console\InputDefinition;
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
        $port = (int) strval($this->input->getArgument('port') ?? '8080');
        $host = (string) ($this->input->getOption('host') ?? '0.0.0.0');

        if ($port < 1 || $port > 65535) {
            $this->output->error('Invalid port. Must be between 1 and 65535.');
            return;
        }

        require dirname(__DIR__, 3) . '/bootstrap/http.php';

        $server = bootstrap_http_server($port, $host);
        $routes = Route::getRoutes();

        $this->output->newLine();
        $this->output->info(sprintf('Server running on [http://%s:%d].', $host, $port));
        $this->output->comment('Press Ctrl+C to stop the server');
        $this->output->newLine();

        if (count($routes) > 0) {
            $this->output->comment('Registered routes (' . count($routes) . '):');
            foreach ($routes as $route) {
                $method = str_pad($route['method'], 7);
                $this->output->comment('  ' . $method . ' ' . $route['path']);
            }
            $this->output->newLine();
        }

        $server->run();
    }
}
