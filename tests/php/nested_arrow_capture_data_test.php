<?php
namespace tests\php;

/**
 * 嵌套箭头应捕获外层方法参数（含默认 []），对齐 Filament Page::getWidgetsSchemaComponents。
 */
class NestedArrow_CaptureHost
{
    public function build(array $widgets, array $data = []): array
    {
        $out = [];
        foreach ($widgets as $w) {
            $out[] = (fn (): array => [
                ...['w' => $w],
                ...$data,
            ]);
        }
        return $out;
    }

    public function buildViaMap(array $widgets, array $data = []): array
    {
        return array_values(array_map(
            fn (string $w, int $k): callable => fn (): array => [
                ...['w' => $w, 'k' => $k],
                ...$data,
            ],
            $widgets,
            array_keys($widgets),
        ));
    }
}

$h = new NestedArrow_CaptureHost();
$fns = $h->build(['A']);
$r = $fns[0]();
if (!is_array($r) || ($r['w'] ?? null) !== 'A') {
    \Log::fatal('单层箭头捕获失败: '.var_export($r, true));
}

$fns2 = $h->buildViaMap(['A', 'B']);
try {
    $r2 = $fns2[0]();
    if (!is_array($r2)) {
        \Log::fatal('嵌套箭头返回非数组: '.gettype($r2));
    }
    if (($r2['w'] ?? null) !== 'A') {
        \Log::fatal('嵌套箭头捕获 w 失败: '.var_export($r2, true));
    }
    // $data 默认 []，展开后不应报错；键 k 存在
    if (!array_key_exists('k', $r2)) {
        \Log::fatal('嵌套箭头缺少 k');
    }
    \Log::info('nested_arrow_capture_data 测试通过');
} catch (\Throwable $e) {
    \Log::fatal('嵌套箭头调用异常: '.$e->getMessage());
}
