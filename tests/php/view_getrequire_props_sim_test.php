<?php

namespace tests\php;

/**
 * 模拟 Laravel Filesystem::getRequire：闭包内 extract + require 视图。
 */
$dir = sys_get_temp_dir().DIRECTORY_SEPARATOR.'origami_view_require';
@mkdir($dir, 0777, true);

$view = $dir.DIRECTORY_SEPARATOR.'drop.php';
file_put_contents($view, <<<'PHP'
<?php
$attributes ??= new stdClass();
$__key = 'trigger';
$__value = null;
$$__key = $$__key ?? $__value;
return [
  'isset' => isset($trigger),
  'type' => gettype($trigger),
  'is_object' => is_object($trigger),
  'class' => is_object($trigger) ? get_class($trigger) : null,
  'html' => (is_object($trigger) && isset($trigger->html)) ? $trigger->html : null,
];
PHP);

class ViewReq_Slot
{
    public $html = 'OpenMe';
    public $attributes;
    public function __construct()
    {
        $this->attributes = (object) ['class' => 'x'];
    }
}

$result = (static function () use ($view) {
    $__data = [
        'trigger' => new \tests\php\ViewReq_Slot(),
        'attributes' => new \stdClass(),
    ];
    extract($__data, EXTR_SKIP);
    return require $view;
})();

if (!is_array($result)) {
    \Log::fatal('require 返回非数组: '.var_export($result, true));
}
if (!$result['isset'] || !$result['is_object']) {
    \Log::fatal('视图内 $trigger 丢失: '.var_export($result, true));
}
if ($result['html'] !== 'OpenMe') {
    \Log::fatal('视图内 html 错误: '.var_export($result, true));
}

\Log::info('view_getrequire_props_sim 测试通过');
