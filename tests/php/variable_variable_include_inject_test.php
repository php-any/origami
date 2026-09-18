<?php

namespace tests\php;

/**
 * include 注入的变量必须能被 $$name 读到，且 $$name = $$name ?? null 不得擦除。
 */
$dir = sys_get_temp_dir().DIRECTORY_SEPARATOR.'origami_vv_inject';
@mkdir($dir, 0777, true);
$view = $dir.DIRECTORY_SEPARATOR.'v.php';
file_put_contents($view, <<<'PHP'
<?php
$__key = 'trigger';
$via = $$__key;
$$__key = $$__key ?? null;
return [isset($via), is_object($via), isset($trigger), is_object($trigger)];
PHP);

class VVInject_Slot
{
    public $html = 'X';
}

$r = (static function () use ($view) {
    extract(['trigger' => new \tests\php\VVInject_Slot()], EXTR_SKIP);
    return require $view;
})();

if ($r !== [true, true, true, true]) {
    \Log::fatal('$$ 读/写注入变量失败: '.var_export($r, true));
}

\Log::info('variable_variable_include_inject 测试通过');
