<?php

namespace App;

use Net\Annotation\Application;

// 引导：配置、模型与服务（控制器由 #[Application] 扫描 src/controllers 加载）
$baseDir = __DIR__;
include($baseDir . '/config/AppConfig.php');
include($baseDir . '/dto/ResponseDTO.php');
include($baseDir . '/model/User.php');
include($baseDir . '/model/Product.php');
include($baseDir . '/services/UserService.php');
include($baseDir . '/services/ProductService.php');
include($baseDir . '/services/AuthService.php');

#[Application(name: 'spring')]
function main($request, $response): void
{
    // #[Application] 会在本函数体首部注入 RegisterRoute，将请求分发到 @*Mapping 控制器方法
}
