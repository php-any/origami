<?php

class EvRouteStub
{
    public function getActionName(): string { return 'Home@index'; }
    public function gatherMiddleware(): array { return ['web']; }
}

class EvRequestStub
{
    public function route() { return new EvRouteStub(); }
    public function ip() { return '127.0.0.1'; }
    public function root() { return 'http://localhost'; }
    public function fullUrl() { return 'http://localhost/foo'; }
    public function method() { return 'GET'; }
    public function headers() { return ['host' => ['localhost']]; }
    public function input() { return []; }
    public function files() { return new class { public function all() { return []; } }; }
    public function hasSession() { return false; }
    public function server($k) { return microtime(true) - 1; }
}

class EvResponseStub
{
    public function getStatusCode() { return 200; }
    public function getContent() { return ''; }
    public function headers() { return new class { public function get($k) { return 'text/html'; } }; }
}

function ev_optional($value)
{
    return new class ($value) {
        private $value;
        public function __construct($v) { $this->value = $v; }
        public function __call($m, $a) { return $this->value === null ? null : $this->value->$m(...$a); }
    };
}

class EvRW
{
    public function headers($headers)
    {
        $headers = ev_collect($headers)
            ->map(fn ($header) => implode(', ', $header))
            ->all();

        return $headers;
    }

    public function payload($payload) { return $payload; }

    public function recordRequest($event)
    {
        $startTime = $event->request->server('REQUEST_TIME_FLOAT');

        $entry = EvIncomingEntry::make([
            'ip_address' => $event->request->ip(),
            'uri' => str_replace($event->request->root(), '', $event->request->fullUrl()) ?: '/',
            'method' => $event->request->method(),
            'controller_action' => ev_optional($event->request->route())->getActionName(),
            'middleware' => array_values(ev_optional($event->request->route())->gatherMiddleware() ?? []),
            'headers' => $this->headers($event->request->headers()),
            'payload' => $this->payload($event->request->input()),
            'session' => $this->payload([]),
            'response_status' => $event->response->getStatusCode(),
            'response' => 'Empty Response',
            'duration' => $startTime ? floor((microtime(true) - $startTime) * 1000) : null,
            'memory' => function_exists('memory_get_peak_usage')
                ? round(memory_get_peak_usage(true) / 1024 / 1024, 1)
                : 0.0,
        ]);

        return $entry;
    }
}

class EvIncomingEntry
{
    public static function make($x) { return $x; }
}

class EvCollection
{
    private $x;
    public function __construct($x) { $this->x = $x; }
    public function map($cb) {
        $out = [];
        foreach ($this->x as $k => $v) { $out[$k] = $cb($v); }
        return new EvCollection($out);
    }
    public function all() { return $this->x; }
}

function ev_collect($x) { return new EvCollection($x); }

$event = (object) ['request' => new EvRequestStub(), 'response' => new EvResponseStub()];
$a = (new EvRW())->recordRequest($event);
if (($a['controller_action'] ?? '') !== 'Home@index') {
    Log::fatal('controller_action');
}
Log::info('exact vendor recordRequest ok');
