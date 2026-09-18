<?php

namespace tests\php;

/**
 * 拆开：include 作用域里 extract 注入 vs $$key ?? 覆盖。
 */
$dir = sys_get_temp_dir().DIRECTORY_SEPARATOR.'origami_view_require2';
@mkdir($dir, 0777, true);

$viewBefore = $dir.DIRECTORY_SEPARATOR.'before.php';
file_put_contents($viewBefore, '<?php return [isset($trigger), is_object($trigger), is_object($trigger)?$trigger->html:null];');

$viewAfter = $dir.DIRECTORY_SEPARATOR.'after.php';
file_put_contents($viewAfter, <<<'PHP'
<?php
$__key = 'trigger';
$__value = null;
$before = isset($trigger);
$$__key = $$__key ?? $__value;
return [$before, isset($trigger), gettype($trigger)];
PHP);

class ViewReq2_Slot
{
    public $html = 'OpenMe';
}

$r1 = (static function () use ($viewBefore) {
    extract(['trigger' => new \tests\php\ViewReq2_Slot()], EXTR_SKIP);
    return require $viewBefore;
})();
if ($r1[0] !== true || $r1[1] !== true || $r1[2] !== 'OpenMe') {
    \Log::fatal('include 前应能看到 extract 对象: '.var_export($r1, true));
}

$r2 = (static function () use ($viewAfter) {
    extract(['trigger' => new \tests\php\ViewReq2_Slot()], EXTR_SKIP);
    return require $viewAfter;
})();
if ($r2[0] !== true) {
    \Log::fatal('?? 前 isset 应为 true: '.var_export($r2, true));
}
if ($r2[1] !== true || $r2[2] !== 'class') {
    \Log::fatal('?? 后不应丢掉对象: '.var_export($r2, true));
}

\Log::info('view_props_vv_wipe 测试通过');
