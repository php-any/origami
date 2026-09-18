<?php

namespace tests\php;

/**
 * sscanf 两参数返回数组；覆盖 Filament Color 常用格式。
 */

$rgb = sscanf('#ff00aa', '#%02x%02x%02x');
if (!is_array($rgb) || count($rgb) !== 3 || $rgb[0] !== 255 || $rgb[1] !== 0 || $rgb[2] !== 170) {
    Log::fatal('sscanf hex 失败: ' . var_export($rgb, true));
}

$oklch = sscanf('oklch(0.628 0.258 29.234)', 'oklch(%f %f %f)');
if (!is_array($oklch) || count($oklch) !== 3) {
    Log::fatal('sscanf oklch 失败: ' . var_export($oklch, true));
}
if (abs($oklch[0] - 0.628) > 0.0001 || abs($oklch[1] - 0.258) > 0.0001 || abs($oklch[2] - 29.234) > 0.0001) {
    Log::fatal('sscanf oklch 数值错误: ' . var_export($oklch, true));
}

$rgb2 = sscanf('rgb(1,2,3)', 'rgb(%d,%d,%d)');
if ($rgb2[0] !== 1 || $rgb2[1] !== 2 || $rgb2[2] !== 3) {
    Log::fatal('sscanf rgb 失败: ' . var_export($rgb2, true));
}

Log::info('sscanf 测试通过');
