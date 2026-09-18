<?php

namespace tests\php;

/**
 * Symfony Routing：Route 编译、匹配、生成；laravel13 原生层与 vendor PHP 共用断言。
 */

$autoload = dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';
$native = class_exists(\Symfony\Component\Routing\Route::class, false);
if (!$native) {
    if (!is_file($autoload)) {
        Log::info('skip: 缺少 laravel13 vendor，跳过 Symfony Routing 测试');
        return;
    }
    require $autoload;
}

use Symfony\Component\Routing\Exception\MethodNotAllowedException;
use Symfony\Component\Routing\Exception\ResourceNotFoundException;
use Symfony\Component\Routing\Exception\RouteNotFoundException;
use Symfony\Component\Routing\Generator\UrlGenerator;
use Symfony\Component\Routing\Generator\UrlGeneratorInterface;
use Symfony\Component\Routing\Matcher\UrlMatcher;
use Symfony\Component\Routing\RequestContext;
use Symfony\Component\Routing\Route;
use Symfony\Component\Routing\RouteCollection;
use Symfony\Component\Routing\RouteCompiler;

$route = new Route('/users/{id}', [], [], ['utf8' => true], '', [], ['GET']);
if ($route->getPath() !== '/users/{id}') {
    Log::fatal('Route::getPath 失败: '.$route->getPath());
}
$compiled = $route->compile();
$prefix = $compiled->getStaticPrefix();
if ($prefix !== '/users') {
    Log::fatal('CompiledRoute::getStaticPrefix 失败: '.$prefix);
}
$regex = $compiled->getRegex();
if (!preg_match($regex, '/users/42')) {
    Log::fatal('编译正则未匹配 /users/42: '.$regex);
}
if (preg_match($regex, '/posts/42')) {
    Log::fatal('编译正则不应匹配 /posts/42: '.$regex);
}

$opt = new Route('/blog/{page}', ['page' => '1'], [], ['utf8' => true]);
$optRx = $opt->compile()->getRegex();
if (!preg_match($optRx, '/blog') || !preg_match($optRx, '/blog/2')) {
    Log::fatal('可选参数路由编译失败: '.$optRx);
}
if (($opt->getDefaults()['page'] ?? null) !== '1') {
    Log::fatal('Route defaults 关联数组丢失: '.json_encode($opt->getDefaults()));
}
if (substr($optRx, -1) !== 'u') {
    Log::fatal('utf8 选项未写入编译正则: '.$optRx);
}

$root = new Route('/', [], [], ['utf8' => true]);
if (!preg_match($root->compile()->getRegex(), '/')) {
    Log::fatal('根路径 / 编译失败');
}

$col = new RouteCollection();
$col->add('users.show', $route);
$col->add('home', $root);
if ($col->count() !== 2) {
    Log::fatal('RouteCollection::count 失败');
}
if ($col->get('users.show') === null) {
    Log::fatal('RouteCollection::get 失败');
}

$ctx = new RequestContext('', 'GET', 'localhost', 'http');
$matcher = new UrlMatcher($col, $ctx);
$matched = $matcher->match('/users/7');
if (($matched['_route'] ?? '') !== 'users.show' || ($matched['id'] ?? '') !== '7') {
    Log::fatal('UrlMatcher::match 失败: '.json_encode($matched));
}

try {
    $matcher->match('/missing');
    Log::fatal('缺失路由应抛 ResourceNotFoundException');
} catch (ResourceNotFoundException $e) {
    // ok
}

$postOnly = new Route('/submit', [], [], ['utf8' => true], '', [], ['POST']);
$col->add('submit', $postOnly);
try {
    $matcher->match('/submit');
    Log::fatal('方法不匹配应抛 MethodNotAllowedException');
} catch (MethodNotAllowedException $e) {
    // ok
}

$gen = new UrlGenerator($col, $ctx);
$url = $gen->generate('users.show', ['id' => 9]);
if ($url !== '/users/9') {
    Log::fatal('UrlGenerator::generate 失败: '.$url);
}
$abs = $gen->generate('users.show', ['id' => 9], UrlGeneratorInterface::ABSOLUTE_URL);
if (strpos($abs, 'http://localhost/users/9') === false) {
    Log::fatal('ABSOLUTE_URL 失败: '.$abs);
}
try {
    $gen->generate('nope');
    Log::fatal('未知路由应抛 RouteNotFoundException');
} catch (RouteNotFoundException $e) {
    // ok
}

$inline = new Route('/item/{id<\d+>}', [], [], ['utf8' => true]);
if ($inline->getRequirement('id') !== '\d+') {
    Log::fatal('内联 requirement 失败: '.$inline->getRequirement('id'));
}
if (!preg_match($inline->compile()->getRegex(), '/item/12')) {
    Log::fatal('内联 requirement 编译后应匹配 /item/12');
}

if (RouteCompiler::SEPARATORS === '') {
    Log::fatal('RouteCompiler::SEPARATORS 应非空');
}

Log::info('symfony_routing 测试通过');
