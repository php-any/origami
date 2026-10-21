<?php

namespace HotReload\Controller;

use HotReload\Service\GreetingService;
use HotReload\View\GreetingView;

class GreetingController {
    private GreetingService $service;

    public function __construct(GreetingService $service) {
        $this->service = $service;
    }

    private function data($request): array {
        $query = $request->query();
        $name = $query['name'] ?? 'Origami';
        return [
            'controller' => 'controller-v1',
            'message' => $this->service->greet($name),
        ];
    }

    public function show($request, $response): void {
        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write(json_encode($this->data($request)));
    }

    public function home($request, $response): void {
        $view = new GreetingView();
        $response->header('Content-Type', 'text/html; charset=utf-8');
        $response->write($view->render($this->data($request)));
    }
}
