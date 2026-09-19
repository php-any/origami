<?php

namespace tests\php;

/**
 * 单行注释遇 ?> 即结束 PHP（与换行同等），其后 HTML 作为输出。
 * 对齐 PHP：// 只注释到当前 PHP 块结束。
 */

class LineCommentPhpClose_HtmlEmbed
{
    public function toHtml(): string
    {
        ob_start(); ?>
<div>before</div>
<?php // tabindex stays in comment until close.?>
<span>after</span>
<?php
        return ob_get_clean();
    }
}

$html = (new LineCommentPhpClose_HtmlEmbed())->toHtml();
if (strpos($html, '<div>before</div>') === false) {
    \Log::fatal('缺少 before HTML: ' . var_export($html, true));
}
if (strpos($html, '<span>after</span>') === false) {
    \Log::fatal('?> 未结束注释后的 HTML: ' . var_export($html, true));
}

$q = 1; // question mark ? is not a closer
if ($q !== 1) {
    \Log::fatal('普通 ? 不应结束注释');
}

\Log::info('line_comment_php_close 测试通过');
