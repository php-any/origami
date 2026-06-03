<?php

namespace App\Controller;

use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;

class HomeController
{
    public function index(Request $request): Response
    {
        return new JsonResponse([
            'message' => 'Hello from Origami Symfony!',
            'path' => $request->getPathInfo(),
            'method' => $request->getMethod(),
        ]);
    }

    public function health(): Response
    {
        return new JsonResponse(['status' => 'ok']);
    }
}
