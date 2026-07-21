<?php

namespace App\Controller;

use Net\Annotation\Controller;
use Net\Annotation\GetMapping;
use Net\Annotation\Route;
use Net\Http\Response;

#[Controller]
#[Route(prefix: "/")]
class PageController
{
    #[GetMapping(path: "/")]
    public function home(Response $response): void
    {
        $pagesDir = dirname(dirname(__DIR__)) . "/public/";
        $response->view($pagesDir . "index.html", []);
    }
}
