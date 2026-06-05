<?php

namespace App;

use Net\Annotation\Application;

// 扫描整个 src 目录，包括 controllers 和 middleware
#[Application(name: 'spring', scan: __DIR__)]
function main($request, $response): void
{

}
