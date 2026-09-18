<?php

namespace tests\php;

/**
 * 嵌套 ob_start / ob_get_clean，对齐 Filament Grid::toEmbeddedHtml 模式。
 */
function NestedOb_inner(): string
{
    ob_start();
    echo 'INNER';
    return ob_get_clean();
}

function NestedOb_outer(): string
{
    ob_start();
    echo 'A';
    echo NestedOb_inner();
    echo 'B';
    return ob_get_clean();
}

$out = NestedOb_outer();
if ($out !== 'AINNERB') {
    \Log::fatal('nested ob 失败: '.var_export($out, true));
}

// Filament 风格：方法内 ?> HTML <?= ?> 混写
class NestedOb_GridLike
{
    public function toEmbeddedHtml(): string
    {
        $child = 'CHILD_HTML';

        ob_start(); ?>

        <div><?= $child ?></div>

        <?php return ob_get_clean();
    }

    public function toEmbeddedHtmlViaCall(): string
    {
        ob_start(); ?>

        <div><?= $this->childHtml() ?></div>

        <?php return ob_get_clean();
    }

    public function childHtml(): string
    {
        ob_start();
        echo 'NESTED_CHILD';
        return ob_get_clean();
    }
}

$g = new NestedOb_GridLike();
$h1 = $g->toEmbeddedHtml();
if (!str_contains($h1, 'CHILD_HTML')) {
    \Log::fatal('GridLike HTML 混写失败: '.var_export($h1, true));
}

$h2 = $g->toEmbeddedHtmlViaCall();
if (!str_contains($h2, 'NESTED_CHILD')) {
    \Log::fatal('GridLike 嵌套 ob+混写失败: '.var_export($h2, true));
}

\Log::info('nested_ob_html_embed 测试通过');
