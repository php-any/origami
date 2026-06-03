<?php

use Symfony\Component\Routing\Route;
use Symfony\Component\Routing\RouteCollection;
use Symfony\Component\Routing\RequestContext;
use Symfony\Component\Routing\Matcher\UrlMatcher;
use Symfony\Component\Routing\Generator\UrlGenerator;

// step02: symfony/routing

$route = new Route('/users/{id}', ['_controller' => 'App\\Controller\\UserController::show'], ['id' => '\d+']);
$compiled = $route->compile();
step_check('Route::compile', $compiled->getRegex() !== '');
step_check('Route regex has id', str_contains($compiled->getRegex(), 'id'));

$routes = new RouteCollection();
$routes->add('home', new Route('/', ['_controller' => 'App\\Controller\\HomeController::index']));
$routes->add('user', new Route('/users/{id}', ['_controller' => 'App\\Controller\\UserController::show'], ['id' => '\d+']));

$context = new RequestContext('', 'GET', 'localhost', 'http', 80);
$matcher = new UrlMatcher($routes, $context);

$match = $matcher->match('/users/42');
step_check('UrlMatcher /users/42', isset($match['_route']) && $match['_route'] === 'user');
step_check('UrlMatcher param id', ($match['id'] ?? null) === '42');

$home = $matcher->match('/');
step_check('UrlMatcher /', ($home['_route'] ?? '') === 'home');

$generator = new UrlGenerator($routes, $context);
$url = $generator->generate('user', ['id' => 99]);
step_check('UrlGenerator', str_contains($url, '99'), 'url=' . $url);

step_check('step02_routing', true, 'symfony/routing');
