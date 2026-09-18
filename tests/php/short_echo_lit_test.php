<?php

namespace tests\php;

/**
 * <?= 字面量 vs 变量，确认是解析丢失还是变量槽问题。
 */
class ShortEcho_Lit
{
    public function lit(): string
    {
        ob_start();
        echo 'A';
        ?><?= 'CHILD_HTML' ?><?php
        echo 'B';
        return ob_get_clean();
    }

    public function varEcho(): string
    {
        $child = 'CHILD_HTML';
        ob_start();
        echo 'A';
        echo $child;
        echo 'B';
        return ob_get_clean();
    }
}

$o = new ShortEcho_Lit();
$a = $o->lit();
if ($a !== 'ACHILD_HTMLB') {
    \Log::fatal('lit: '.var_export($a, true));
}
$b = $o->varEcho();
if ($b !== 'ACHILD_HTMLB') {
    \Log::fatal('varEcho: '.var_export($b, true));
}
\Log::info('short_echo_lit 测试通过');
