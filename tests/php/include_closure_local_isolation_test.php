<?php

namespace tests\php;

/**
 * Laravel Filesystem::getRequire 每次都是新闭包：子视图的 $column 不得写穿父视图。
 * Origami 若把 include 局部挂到进程级 $GLOBALS，Filament tables index 里
 * 子视图 foreach ($arr as $column => $v) 会把父级 $column 变成 3，随后 getName() 失败。
 */

class InclIsolate_Col
{
    public function __construct(public string $name)
    {
    }

    public function getName(): string
    {
        return $this->name;
    }
}

$dir = sys_get_temp_dir().DIRECTORY_SEPARATOR.'origami_incl_col_iso';
@mkdir($dir, 0777, true);
$child = $dir.DIRECTORY_SEPARATOR.'child.php';
$parent = $dir.DIRECTORY_SEPARATOR.'parent.php';

file_put_contents($child, <<<'PHP'
<?php
foreach ([10, 20, 30] as $column => $v) {
}
return $column;
PHP);

file_put_contents($parent, '<?php
$columns = [
    \'a\' => new \\tests\\php\\InclIsolate_Col(\'a\'),
    \'b\' => new \\tests\\php\\InclIsolate_Col(\'b\'),
];
$names = [];
foreach ($columns as $column) {
    (static function () {
        extract([], EXTR_SKIP);
        return require '.var_export($child, true).';
    })();
    if (!is_object($column)) {
        return [\'fail\', $column];
    }
    $names[] = $column->getName();
}
return [\'ok\', $names];
');

$got = (static function () use ($parent) {
    extract([], EXTR_SKIP);
    return require $parent;
})();

@unlink($child);
@unlink($parent);
@rmdir($dir);

if (!is_array($got) || ($got[0] ?? '') !== 'ok' || ($got[1] ?? null) !== ['a', 'b']) {
    \Log::fatal('getRequire 闭包间 $column 被写穿: '.var_export($got, true));
}

\Log::info('include 闭包局部隔离测试通过');
