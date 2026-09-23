<?php

/**
 * UrlGenerator 原生实现冒烟。未 AddClass 时跳过（无 vendor 时无法加载 PHP 版）。
 */
if (!class_exists('Illuminate\\Routing\\UrlGenerator', false)) {
    echo "skip: UrlGenerator 原生类未注册\n";
    return;
}

$request = \Illuminate\Http\Request::create('http://example.com', 'GET', [], [], [], [
    'HTTP_REFERER' => 'http://example.com/previous-page',
]);

class FakeNamedRoute
{
    public function __construct(private string $uri) {}

    public function uri(): string
    {
        return $this->uri;
    }
}

class FakeRouteCollection
{
    private array $named = [];

    public function addNamed(string $name, string $uri): void
    {
        $this->named[$name] = new FakeNamedRoute($uri);
    }

    public function getByName(string $name): ?FakeNamedRoute
    {
        return $this->named[$name] ?? null;
    }

    public function getByAction(string $action): ?FakeNamedRoute
    {
        return $this->named[$action] ?? null;
    }
}

$routes = new FakeRouteCollection();
$routes->addNamed('panel.dashboard', '/admin/{panel}');

$url = new \Illuminate\Routing\UrlGenerator($routes, $request);

$generated = $url->route('panel.dashboard', ['panel' => 'main']);
$expected = 'http://example.com/admin/main';
if ($generated !== $expected) {
    echo "FAIL route() expected {$expected} got {$generated}\n";
    exit(1);
}

echo "OK illuminate_url_generator_test\n";
