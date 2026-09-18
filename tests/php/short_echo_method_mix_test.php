<?php

namespace tests\php;

/**
 * 缩小 <?= 在方法内混写时变量/echo 丢失问题。
 */
class ShortEcho_Mix
{
    public function plainEcho(): string
    {
        $child = 'CHILD_HTML';
        ob_start();
        echo $child;
        return ob_get_clean();
    }

    public function shortEchoOnly(): string
    {
        $child = 'CHILD_HTML';
        ob_start();
        echo 'A';
        ?><?= $child ?><?php
        echo 'B';
        return ob_get_clean();
    }

    public function shortEchoInDiv(): string
    {
        $child = 'CHILD_HTML';
        ob_start(); ?>
        <div><?= $child ?></div>
        <?php return ob_get_clean();
    }

    public function shortEchoCall(): string
    {
        ob_start(); ?>
        <div><?= $this->val() ?></div>
        <?php return ob_get_clean();
    }

    public function val(): string
    {
        return 'VAL_HTML';
    }
}

$o = new ShortEcho_Mix();
$a = $o->plainEcho();
if ($a !== 'CHILD_HTML') {
    \Log::fatal('plainEcho: '.var_export($a, true));
}
$b = $o->shortEchoOnly();
if ($b !== 'ACHILD_HTMLB') {
    \Log::fatal('shortEchoOnly: '.var_export($b, true));
}
$c = $o->shortEchoInDiv();
if (!str_contains($c, 'CHILD_HTML')) {
    \Log::fatal('shortEchoInDiv: '.var_export($c, true));
}
$d = $o->shortEchoCall();
if (!str_contains($d, 'VAL_HTML')) {
    \Log::fatal('shortEchoCall: '.var_export($d, true));
}

\Log::info('short_echo_method_mix 测试通过');
