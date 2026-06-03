<?php

use App\Controller\HomeController;
use Symfony\Component\EventDispatcher\EventDispatcher;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\RequestStack;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Controller\ControllerResolver;
use Symfony\Component\HttpKernel\EventListener\RouterListener;
use Symfony\Component\HttpKernel\HttpKernel;
use Symfony\Component\Routing\Matcher\UrlMatcher;
use Symfony\Component\Routing\RequestContext;
use Symfony\Component\Routing\Route;
use Symfony\Component\Routing\RouteCollection;

require __DIR__ . '/vendor/autoload.php';

/**
 * 最小 Symfony HttpKernel 引导（无 DI 容器，便于 Origami 逐步验证）
 */
function symfony_create_kernel(): HttpKernel
{
    $routes = new RouteCollection();
    $routes->add('home', new Route('/', ['_controller' => [HomeController::class, 'index']], [], [], '', [], ['GET']));
    $routes->add('health', new Route('/health', ['_controller' => [HomeController::class, 'health']], [], [], '', [], ['GET']));

    $context = new RequestContext('', 'GET', 'localhost', 'http', 80);
    $matcher = new UrlMatcher($routes, $context);

    $dispatcher = new EventDispatcher();
    $dispatcher->addSubscriber(new RouterListener($matcher, new RequestStack(), $context));

    $resolver = new ControllerResolver();
    return new HttpKernel($dispatcher, $resolver);
}

function symfony_handle(HttpKernel $kernel, string $method, string $uri, string $host = 'localhost'): Response
{
    $_SERVER['REQUEST_METHOD'] = $method;
    $_SERVER['REQUEST_URI'] = $uri;
    $_SERVER['HTTP_HOST'] = $host;
    $_SERVER['SERVER_NAME'] = $host;
    $_SERVER['SERVER_PORT'] = '80';
    $_SERVER['SCRIPT_FILENAME'] = __DIR__ . '/public/index.php';
    $_SERVER['SCRIPT_NAME'] = '/index.php';
    $_SERVER['PHP_SELF'] = '/index.php';

    $request = Request::createFromGlobals();
    return $kernel->handle($request);
}
