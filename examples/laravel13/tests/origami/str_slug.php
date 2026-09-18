<?php

/**
 * Illuminate\Support\Str::slug 热路径（config/cache.php prefix）。
 */
$got = Illuminate\Support\Str::slug((string) 'Origami Admin');
if ($got !== 'origami-admin') {
    fwrite(STDERR, "Str::slug 失败: {$got}\n");
    exit(1);
}
echo "str_slug_ok\n";
