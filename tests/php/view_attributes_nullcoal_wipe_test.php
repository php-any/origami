<?php

namespace tests\php;

/**
 * 复现：$attributes ??= 之后 $$trigger ?? null 是否丢掉 extract 的 $trigger。
 */
$dir = sys_get_temp_dir().DIRECTORY_SEPARATOR.'origami_view_require3';
@mkdir($dir, 0777, true);
$view = $dir.DIRECTORY_SEPARATOR.'v.php';
file_put_contents($view, <<<'PHP'
<?php
$attributes ??= new stdClass();
$before = [isset($trigger), is_object($trigger)];
$__key = 'trigger';
$__value = null;
$$__key = $$__key ?? $__value;
return [
  'before' => $before,
  'after_isset' => isset($trigger),
  'after_type' => gettype($trigger),
  'after_obj' => is_object($trigger),
];
PHP);

class ViewReq3_Slot
{
    public $html = 'OpenMe';
}

$r = (static function () use ($view) {
    extract([
        'trigger' => new \tests\php\ViewReq3_Slot(),
        'attributes' => new \stdClass(),
    ], EXTR_SKIP);
    return require $view;
})();

if ($r['before'][0] !== true || $r['before'][1] !== true) {
    \Log::fatal('??=attributes 前 trigger 应在: '.var_export($r, true));
}
if ($r['after_isset'] !== true || $r['after_obj'] !== true) {
    \Log::fatal('?? 后 trigger 丢失: '.var_export($r, true));
}

\Log::info('view_attributes_nullcoal_wipe 测试通过');
