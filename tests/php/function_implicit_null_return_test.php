<?php

namespace tests\php;

/**
 * PHP：函数/闭包/方法没有 return 时应返回 null，而不是最后一条语句的值。
 * Laravel View::render($callback) 与 Livewire 根标签注入依赖此语义。
 */

function implicitNullReturn_fn()
{
    $x = 42;
}

if (implicitNullReturn_fn() !== null) {
    Log::fatal('函数无 return 应返回 null, got=' . var_export(implicitNullReturn_fn(), true));
}

$fn = function () {
    $y = 'sections';
};
if ($fn() !== null) {
    Log::fatal('闭包无 return 应返回 null, got=' . var_export($fn(), true));
}

class ImplicitNullReturn_Host
{
    public $sections = [];

    public function extract()
    {
        $this->sections = ['ok'];
    }
}

$host = new ImplicitNullReturn_Host();
if ($host->extract() !== null) {
    Log::fatal('方法无 return 应返回 null, got=' . var_export($host->extract(), true));
}

$contents = '<div class="root">ok</div>';
$callback = function ($view) {
    $view->sections = ['ok'];
};
$response = $callback($host);
$html = is_null($response) ? $contents : $response;
if ($html !== $contents) {
    Log::fatal('View::render 回调无 return 时应采用 HTML, got=' . var_export($html, true));
}

$arrow = fn ($n) => $n * 2;
if ($arrow(21) !== 42) {
    Log::fatal('箭头函数仍应返回表达式值, got=' . var_export($arrow(21), true));
}

Log::info('函数无 return 返回 null 测试通过');
