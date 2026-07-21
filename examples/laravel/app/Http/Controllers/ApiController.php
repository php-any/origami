<?php

namespace App\Http\Controllers;

use Net\Http\Response;

class ApiController
{
    public function health(Response $response): void
    {
        $response->success([
            'status' => 'ok',
            'framework' => 'Origami Laravel Demo',
        ]);
    }
}
