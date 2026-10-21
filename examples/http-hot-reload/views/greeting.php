<?php

$controller = htmlspecialchars($data['controller'], ENT_QUOTES, 'UTF-8');
$message = htmlspecialchars($data['message'], ENT_QUOTES, 'UTF-8');

return '<!doctype html><html lang="zh-CN"><meta charset="utf-8">'
    .'<title>Origami HTTP 热重载</title><body>'
    .'<h1>view-v1</h1><p>'.$controller.'</p><p>'.$message.'</p>'
    .'<p>修改控制器、服务或视图文件，刷新即可验证热重载。</p>'
    .'</body></html>';
