<?php

namespace tests\php;

/**
 * htmlspecialchars 对齐 PHP 内置：带 __toString 的对象可转成字符串再转义。
 * Laravel 异常页会对 Request（实现 __toString）走 e() / htmlspecialchars。
 */

class Htmlspecialchars_ToStringObj
{
    public function __toString()
    {
        return '<b>x</b>';
    }
}

$s = htmlspecialchars(new Htmlspecialchars_ToStringObj());
if ($s !== '&lt;b&gt;x&lt;/b&gt;') {
    \Log::fatal('对象 __toString 转义错误: ' . var_export($s, true));
}

if (htmlspecialchars('<a>') !== '&lt;a&gt;') {
    \Log::fatal('普通字符串转义错误');
}

\Log::info('htmlspecialchars_tostring 测试通过');
