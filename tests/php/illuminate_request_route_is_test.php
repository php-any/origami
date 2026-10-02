<?php

class RequestRouteIsFixture {
    public function fluent(): self { return $this; }
    public function named(...$patterns): bool {
        return array_any($patterns, fn ($pattern) => $pattern === 'admin.home');
    }
}
$route = new RequestRouteIsFixture();
$request = Illuminate\Http\Request::create('/admin');
if ($request->routeIs('admin.home') !== false) {
    throw new Exception('Request without a route must return false');
}
$request->setRouteResolver(fn () => $route);
if ($request->routeIs('missing', 'admin.home') !== true
    || $request->routeIs('missing') !== false || $request->routeIs() !== false) {
    throw new Exception('Route named patterns must be bound as variadic arguments');
}
$request->setRouteResolver(fn () => $route->fluent());
$active = fn (): bool => $request->routeIs('admin.home');
if ($active() !== true) {
    throw new Exception('Fluent route receiver lost object identity');
}
echo "Illuminate request routeIs OK\n";
