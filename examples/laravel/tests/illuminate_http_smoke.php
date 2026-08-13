<?php

require dirname(__DIR__) . '/vendor/autoload.php';

use Illuminate\Http\Request;
use Illuminate\Http\Response;

// 避免依赖 URI query + array_replace 合并路径（NamedZVal 键在 array_replace 中会丢）
$base = \Symfony\Component\HttpFoundation\Request::create('/users', 'GET', [
    'id' => '42',
    'q' => 'origami',
]);
$request = Request::createFromBase($base);

$path = method_exists($request, 'path') ? $request->path() : ltrim($request->getPathInfo(), '/');
if ($path !== 'users') {
    echo "FAIL: path=";
    var_export($path);
    echo "\n";
    exit(1);
}

$q = $request->get('q');
if ($q !== 'origami') {
    echo "FAIL: q=";
    var_export([$q, $request->query->all()]);
    echo "\n";
    exit(1);
}

$id = $request->get('id');
if ($id != 42 && $id !== '42') {
    echo "FAIL: id=";
    var_export($id);
    echo "\n";
    exit(1);
}

$response = new Response('hello http', 201, ['X-Origami' => '1']);
if ((int) $response->getStatusCode() !== 201) {
    echo "FAIL: status=";
    var_export($response->getStatusCode());
    echo "\n";
    exit(1);
}
if ($response->getContent() !== 'hello http') {
    echo "FAIL: content=";
    var_export($response->getContent());
    echo "\n";
    exit(1);
}

echo "PASS\n";
