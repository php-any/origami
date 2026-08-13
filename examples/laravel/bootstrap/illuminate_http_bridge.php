<?php

/**
 * Illuminate HTTP 桥接：将 Origami Net\Http\Request 转为 Illuminate\Http\Request，
 * 并从 Foundation 容器解析控制器依赖。
 */

use Illuminate\Http\Request as IlluminateRequest;
use Symfony\Component\HttpFoundation\Request as SymfonyRequest;

function illuminate_resolve_controller_param(string $abstract, $netRequest = null)
{
    $abstract = ltrim($abstract, '\\');

    if ($abstract === IlluminateRequest::class) {
        return illuminate_http_request_from_net($netRequest);
    }

    return illuminate_container()->make($abstract);
}

function illuminate_http_request_from_net($netRequest = null): IlluminateRequest
{
    if ($netRequest === null) {
        return IlluminateRequest::create('/', 'GET');
    }

    $method = strtoupper((string) $netRequest->method());
    $fullUrl = (string) $netRequest->fullUrl();
    $path = parse_url($fullUrl, PHP_URL_PATH) ?: '/';
    $queryString = parse_url($fullUrl, PHP_URL_QUERY);
    $query = [];
    if (is_string($queryString) && $queryString !== '') {
        parse_str($queryString, $query);
    }

    $contentType = (string) $netRequest->header('Content-Type');
    $body = $netRequest->body();
    $parameters = $query;
    $content = '';

    if (is_array($body)) {
        $parameters = array_merge($parameters, $body);
        $content = json_encode($body);
    } elseif (is_string($body) && $body !== '') {
        $content = $body;
        if (str_contains($contentType, 'application/json')) {
            $decoded = json_decode($body, true);
            if (is_array($decoded)) {
                $parameters = array_merge($parameters, $decoded);
            }
        } elseif (str_contains($contentType, 'application/x-www-form-urlencoded')) {
            $parsed = [];
            parse_str($body, $parsed);
            if (is_array($parsed)) {
                $parameters = array_merge($parameters, $parsed);
            }
        }
    }

    $all = $netRequest->all();
    if ($all !== null) {
        if ($all instanceof \Traversable) {
            foreach ($all as $key => $value) {
                if (!array_key_exists($key, $parameters)) {
                    $parameters[$key] = $value;
                }
            }
        } elseif (is_array($all)) {
            $parameters = array_merge($parameters, $all);
        }
    }

    $uri = $path;
    if (is_string($queryString) && $queryString !== '') {
        $uri .= '?' . $queryString;
    }

    $server = [
        'REQUEST_METHOD' => $method,
        'REQUEST_URI' => $uri,
        'HTTP_HOST' => parse_url($fullUrl, PHP_URL_HOST) ?: 'localhost',
        'CONTENT_TYPE' => $contentType,
        'HTTP_ACCEPT' => (string) $netRequest->header('Accept'),
        'HTTP_X_REQUESTED_WITH' => (string) $netRequest->header('X-Requested-With'),
        'HTTP_X_CSRF_TOKEN' => (string) $netRequest->header('X-CSRF-Token'),
    ];

    $base = SymfonyRequest::create($uri, $method, $parameters, [], [], $server, $content);

    return IlluminateRequest::createFromBase($base);
}
