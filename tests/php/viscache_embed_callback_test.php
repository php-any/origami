<?php

namespace tests\php;

/**
 * 复现：声明 string 的方法经 try/finally + 箭头回调，在「有可见组件」路径返回闭包。
 * 对齐 Filament Schema::toEmbeddedHtml / withVisibilityCache。
 */

class VisCache_Component
{
    public static bool $enabled = false;

    public static function withVisibilityCache(callable $callback): mixed
    {
        $wasEnabled = self::$enabled;
        if (!$wasEnabled) {
            self::$enabled = true;
        }
        try {
            return $callback();
        } finally {
            if (!$wasEnabled) {
                self::$enabled = false;
            }
        }
    }
}

class VisCache_Schema
{
    /** @var array<int, object> */
    public array $components = [];

    public function toEmbeddedHtml(): string
    {
        return VisCache_Component::withVisibilityCache(fn (): string => $this->renderEmbeddedHtml());
    }

    protected function renderEmbeddedHtml(): string
    {
        $hasVisible = false;
        $componentsWithVisibility = array_map(
            function ($component) use (&$hasVisible): array {
                $visible = true;
                if ($visible) {
                    $hasVisible = true;
                }
                return [$component, $visible];
            },
            $this->components,
        );

        if (!$hasVisible) {
            return '';
        }

        ob_start(); ?>
        <div>
            <?php foreach ($componentsWithVisibility as [$comp, $vis]) { ?>
                <?php if ($vis) { ?>
                    <?= $comp->toHtml() ?>
                <?php } ?>
            <?php } ?>
        </div>
        <?php return ob_get_clean();
    }
}

class VisCache_Comp
{
    public function toHtml(): string
    {
        return VisCache_Component::withVisibilityCache(fn (): string => $this->render());
    }

    protected function render(): string
    {
        return '<span>x</span>';
    }
}

$s = new VisCache_Schema();
$s->components = [new VisCache_Comp()];

try {
    $html = $s->toEmbeddedHtml();
    if (!is_string($html)) {
        \Log::fatal('toEmbeddedHtml 应返回 string, 实际 '.gettype($html));
    }
    if (!str_contains($html, '<span>x</span>')) {
        \Log::fatal('HTML 内容不对: '.$html);
    }
    \Log::info('viscache_embed 测试通过');
} catch (\Throwable $e) {
    \Log::fatal('异常: '.$e->getMessage());
}
