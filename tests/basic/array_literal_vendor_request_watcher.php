<?php

class VrwRouteStub
{
    public function getActionName(): string { return 'Home@index'; }
    public function gatherMiddleware(): array { return ['web']; }
}

class VrwRequestStub
{
    public function route() { return new VrwRouteStub(); }
    public function ip() { return '127.0.0.1'; }
    public function root() { return 'http://localhost'; }
    public function fullUrl() { return 'http://localhost/foo'; }
    public function method() { return 'GET'; }
    public function headers() { return new class { public function all() { return ['host' => ['localhost']]; } }; }
    public function input() { return []; }
    public function files() { return new class { public function all() { return []; } }; }
    public function hasSession() { return false; }
    public function server($k) { return microtime(true) - 1; }
}

class VrwResponseStub
{
    public function getStatusCode() { return 200; }
    public function getContent() { return ''; }
    public function headers() { return new class { public function get($k) { return 'text/html'; } }; }
}

function vrw_optional($value)
{
    return new class ($value) {
        private $value;
        public function __construct($v) { $this->value = $v; }
        public function __call($m, $a) { return $this->value === null ? null : $this->value->$m(...$a); }
    };
}

class VrwRW
{
    public function recordRequest($event): void
    {
        $startTime = $event->request->server('REQUEST_TIME_FLOAT');

        $entry = [
            'ip_address' => $event->request->ip(),
            'uri' => str_replace($event->request->root(), '', $event->request->fullUrl()) ?: '/',
            'method' => $event->request->method(),
            'controller_action' => vrw_optional($event->request->route())->getActionName(),
            'middleware' => array_values(vrw_optional($event->request->route())->gatherMiddleware() ?? []),
            'response_status' => $event->response->getStatusCode(),
            'duration' => $startTime ? floor((microtime(true) - $startTime) * 1000) : null,
            'memory' => function_exists('memory_get_peak_usage')
                ? round(memory_get_peak_usage(true) / 1024 / 1024, 1)
                : 0.0,
        ];

        if (($entry['controller_action'] ?? '') !== 'Home@index') {
            Log::fatal('controller_action');
        }
        if (($entry['middleware'][0] ?? '') !== 'web') {
            Log::fatal('middleware');
        }
    }
}

$event = (object) ['request' => new VrwRequestStub(), 'response' => new VrwResponseStub()];
(new VrwRW())->recordRequest($event);
Log::info('vendor-like request watcher array ok');
