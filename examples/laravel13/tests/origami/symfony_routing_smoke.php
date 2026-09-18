<?php

namespace tests\origami;

/**
 * Symfony Routing 原生加速层：Route 编译 / UrlMatcher / UrlGenerator。
 */

if (!class_exists(\Symfony\Component\Routing\Route::class, false)) {
    \Log::fatal('Symfony\\Component\\Routing\\Route 应已由 std/symfony/routing 预注册');
}

$route = new \Symfony\Component\Routing\Route('/admin/{id}', ['id' => null], [], ['utf8' => true], '', [], ['GET', 'HEAD']);
$compiled = $route->compile();
if ($compiled->getStaticPrefix() !== '/admin') {
    \Log::fatal('staticPrefix 失败: '.$compiled->getStaticPrefix());
}
if (!array_key_exists('id', $route->getDefaults())) {
    \Log::fatal('defaults 应保留 id 键: '.json_encode($route->getDefaults()));
}
if (substr($compiled->getRegex(), -1) !== 'u') {
    \Log::fatal('utf8 选项未写入编译正则: '.$compiled->getRegex());
}
if (!preg_match($compiled->getRegex(), '/admin/1')) {
    \Log::fatal('regex 未匹配 /admin/1: '.$compiled->getRegex());
}

$col = new \Symfony\Component\Routing\RouteCollection();
$col->add('admin.show', $route);
$ctx = new \Symfony\Component\Routing\RequestContext();
$matcher = new \Symfony\Component\Routing\Matcher\UrlMatcher($col, $ctx);
$hit = $matcher->match('/admin/5');
if (($hit['_route'] ?? null) !== 'admin.show' || ($hit['id'] ?? null) !== '5') {
    \Log::fatal('UrlMatcher 失败');
}

$gen = new \Symfony\Component\Routing\Generator\UrlGenerator($col, $ctx);
if ($gen->generate('admin.show', ['id' => 5]) !== '/admin/5') {
    \Log::fatal('UrlGenerator 失败: '.$gen->generate('admin.show', ['id' => 5]));
}

$dumper = new \Symfony\Component\Routing\Matcher\Dumper\CompiledUrlMatcherDumper($col);
$compiledRoutes = $dumper->getCompiledRoutes();
if (!is_array($compiledRoutes) || count($compiledRoutes) < 4) {
    \Log::fatal('CompiledUrlMatcherDumper::getCompiledRoutes 失败');
}
$cm = new \Symfony\Component\Routing\Matcher\CompiledUrlMatcher($compiledRoutes, $ctx);
$hit2 = $cm->match('/admin/5');
if (($hit2['_route'] ?? null) !== 'admin.show') {
    \Log::fatal('CompiledUrlMatcher 失败');
}

\Log::info('symfony_routing smoke 通过');
