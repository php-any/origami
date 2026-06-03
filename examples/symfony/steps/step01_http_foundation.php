<?php

use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\ParameterBag;

// step01: symfony/http-foundation

$_SERVER['REQUEST_METHOD'] = 'GET';
$_SERVER['REQUEST_URI'] = '/hello?foo=bar';
$_SERVER['QUERY_STRING'] = 'foo=bar';
$_GET['foo'] = 'bar';
$_SERVER['HTTP_HOST'] = 'localhost';
$_SERVER['SERVER_NAME'] = 'localhost';
$_SERVER['SERVER_PORT'] = '80';
$_SERVER['SCRIPT_FILENAME'] = dirname(__DIR__) . '/public/index.php';
$_SERVER['SCRIPT_NAME'] = '/index.php';

$req = Request::createFromGlobals();
step_check('Request::createFromGlobals', $req instanceof Request);
step_check('Request pathInfo', $req->getPathInfo() === '/hello', 'path=' . $req->getPathInfo());
step_check('Request query', $req->query->get('foo') === 'bar');

$req2 = Request::create('/api/test', 'POST', ['name' => 'origami'], [], [], [
    'CONTENT_TYPE' => 'application/json',
], '{"id":1}');
step_check('Request::create POST', $req2->getMethod() === 'POST');
step_check('Request content', $req2->getContent() === '{"id":1}');

$resp = new Response('hello', 200, ['X-Test' => '1']);
step_check('Response status', $resp->getStatusCode() === 200);
step_check('Response content', $resp->getContent() === 'hello');

try {
    $json = new JsonResponse(['ok' => true]);
    step_check('JsonResponse', str_contains($json->getContent(), '"ok"'));
} catch (\Throwable $e) {
    step_check('JsonResponse', false, $e->getMessage());
}

$bag = new ParameterBag(['a' => 1, 'b' => 2]);
step_check('ParameterBag all', count($bag->all()) === 2);

step_check('step01_http_foundation', true, 'symfony/http-foundation');
