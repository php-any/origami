<?php

/**
 * 编译 Blade 常见形态：第一个 <?php 块里 use enum，?> 后接 HTML，再开 <?php 用短名。
 * use 是文件作用域，后面的 PHP 块必须仍能解析 Alignment。
 */

enum UseAliasAfterHtml_Alignment: string
{
    case Center = 'center';
    case Start = 'start';
}

use UseAliasAfterHtml_Alignment as Alignment;

?>
<div class="root">
<?php

$v = Alignment::Center;
if ($v->value !== 'center') {
    \Log::fatal('HTML 后短名 Alignment::Center 失败: '.var_export($v, true));
}

\Log::info('use_alias_after_html 测试通过');
