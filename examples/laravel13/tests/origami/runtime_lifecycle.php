<?php

// Test-only routes exercise the official application without changing app/vendor.
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);
$kernel->bootstrap();

class LifecycleGlobal {
    public function handle($request, $next) {
        if ($request->path() === '__runtime/file') { return $next($request); }
        $id = $request->query('id', 'one');
        register_shutdown_function(function () use ($id) { echo '|shutdown:'.$id; });
        if ($request->path() === '__runtime/short') {
            return new Illuminate\Http\Response('short:'.$id);
        }
        $response = $next($request);
        if ($request->path() === '__runtime/order') {
            $response->setContent('global:'.$id.'>'.$response->getContent().'<global');
        }
        return $response;
    }
    public function terminate($request, $response) {
        if ($request->path() !== '__runtime/file') { echo '|global-term:'.$request->query('id', 'one'); }
    }
}
class LifecycleGroup {
    public function handle($request, $next) {
        $response = $next($request);
        $response->setContent('group>'.$response->getContent().'<group');
        return $response;
    }
    public function terminate($request, $response) { echo '|group-term'; }
}
class LifecycleRoute {
    public function handle($request, $next, $label) {
        $response = $next($request);
        $response->setContent($label.'>'.$response->getContent().'<'.$label);
        return $response;
    }
    public function terminate($request, $response) { echo '|route-term'; }
}
$kernel->setGlobalMiddleware([LifecycleGlobal::class]);
$kernel->setMiddlewareGroups(['lifecycle' => [LifecycleGroup::class]]);
$kernel->setMiddlewareAliases(['lifecycle.route' => LifecycleRoute::class]);
$router = $app['router'];
$app->scoped('runtime.scoped', function () { return new stdClass; });
// Resolve before the worker starts: the request must forget this cached object.
$app->make('runtime.scoped')->count = 99;
class LifecycleThirdPartyManager {
    public int $count = 0;
    public $owner;
    public $self;
    public array $links;
    public function touch(): int { return ++$this->count; }
}
class LifecycleThirdPartyFacade extends Illuminate\Support\Facades\Facade {
    protected static function getFacadeAccessor() { return 'runtime.manager'; }
}
$graphCaptured = [];
$manager = new LifecycleThirdPartyManager();
$manager->owner = $app;
$manager->self = $manager;
$manager->links = ['captured' => &$graphCaptured, 'alias' => $manager];
$app->instance('runtime.manager', $manager);
$app->instance('runtime.manager.alias', $manager);
$app['events']->listen('runtime.graph', function ($id) use (&$graphCaptured) { $graphCaptured[] = $id; });
// Prime the Facade cache before requests; its object must map to the same graph.
LifecycleThirdPartyFacade::getFacadeRoot();
$router->get('/__runtime/graph', function () {
    $manager = app('runtime.manager');
    $id = request()->query('id', 'one');
    app('events')->dispatch('runtime.graph', [$id]);
    if ($manager !== app('runtime.manager.alias') || $manager->self !== $manager || $manager->links['alias'] !== $manager || $manager->owner !== app()) {
        throw new RuntimeException('singleton aliases or cycle lost');
    }
    if ($manager->links['captured'] !== [$id]) { throw new RuntimeException('callback capture disconnected or leaked'); }
    return (string) LifecycleThirdPartyFacade::touch();
});
$router->get('/__runtime/order', function () { return 'body'; })
    ->middleware(['lifecycle', 'lifecycle.route:route']);
$router->get('/__runtime/short', function () { return 'unreachable'; });
class LifecycleCustomResponse extends Illuminate\Http\Response {
    public function sendContent(): static {
        echo 'custom-send';
        return $this;
    }
}
class LifecycleThrowingResponse extends Illuminate\Http\Response {
    public function sendContent(): static {
        echo 'partial-custom';
        throw new RuntimeException('custom response exception after headers');
    }
}
$router->get('/__runtime/custom-send', function () { return new LifecycleCustomResponse('stored-content'); });
$router->get('/__runtime/custom-throw', function () { return new LifecycleThrowingResponse('stored-content'); });

$router->get('/__runtime/stream', function () {
    return new Symfony\Component\HttpFoundation\StreamedResponse(function () {
        echo 'first|';
        flush();
        echo 'second';
    });
});
$router->get('/__runtime/chunks', function () {
    return new Symfony\Component\HttpFoundation\StreamedResponse(['a', 'b']);
});
$router->get('/__runtime/exit', function () { echo 'exit-body'; exit; });
$router->get('/__runtime/empty', function () { return new Illuminate\Http\Response('discard', 204); });
$router->get('/__runtime/not-modified', function () { return new Illuminate\Http\Response('discard', 304); });
$router->get('/__runtime/file', function () { return new Symfony\Component\HttpFoundation\BinaryFileResponse(__FILE__); });
$router->get('/__runtime/cancel', function () {
    register_shutdown_function(function () { runtimeShutdownSignal('shutdown'); });
    runtimeShutdownSignal('started');
    sleep(20);
    runtimeShutdownSignal('resumed');
    return 'unexpected';
});
$router->get('/__runtime/throw', function () { throw new RuntimeException('lifecycle exception'); });
$app->make(Illuminate\Contracts\Debug\ExceptionHandler::class)->renderable(function (RuntimeException $exception, $request) {
    if ($request->path() === '__runtime/throw') { return new Illuminate\Http\Response('handled-exception', 500); }
});
$router->get('/__runtime/stream-throw', function () {
    return new Symfony\Component\HttpFoundation\StreamedResponse(function () {
        echo 'partial-stream';
        throw new RuntimeException('stream exception after headers');
    });
});
$router->get('/__runtime/panic', function () { runtimePanicSignal(); });
$router->get('/__runtime/scoped', function () use ($app) {
    $scoped = app('runtime.scoped');
    $scoped->count = ($scoped->count ?? 0) + 1;
    return (string) $scoped->count;
});
$router->get('/__runtime/events', function () {
    $events = app('events');
    $events->listen('runtime.local', function () { return 'local'; });
    return (string) count($events->getRawListeners()['runtime.local']);
});
// Exercise vendor feature callbacks after the worker rebinding step. Static
// method closures have lexical class scope but no instance object identity.
class LifecycleLivewireComponent extends Livewire\Component {}
$router->get('/__runtime/livewire-call', function () {
    $component = new LifecycleLivewireComponent();
    $context = new Livewire\Mechanisms\HandleComponents\ComponentContext($component);
    $called = false;
    Livewire\trigger('call', $component, '$refresh', [], $context, function () use (&$called) { $called = true; }, [], 0);
    return $called ? 'magic-action' : 'missing-magic-action';
});
return $app;
