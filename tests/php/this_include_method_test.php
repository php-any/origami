<?php

namespace tests\php;

/**
 * $this->include($path, $context) 必须是方法调用，不能当成 include 语句。
 * Symfony 异常页：<?= $this->include('views/exception.html.php', $context); ?>
 */
class ThisIncludeMethod_Host
{
    public function include(string $name, $context = null)
    {
        return 'INC:'.$name.':'.(is_array($context) ? 'arr' : 'null');
    }
}

$host = new ThisIncludeMethod_Host();
$one = $host->include('favicon.png.base64');
$two = $host->include('views/exception.html.php', ['e' => 1]);
if ($one !== 'INC:favicon.png.base64:null' || $two !== 'INC:views/exception.html.php:arr') {
    Log::fatal('this->include 方法调用失败: '.var_export([$one, $two], true));
}

Log::info('$this->include 方法调用测试通过');
