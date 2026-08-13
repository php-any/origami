<?php

namespace tests\php;

class StaticIn_Loader
{
    protected static $ns = 'Foo\\';

    public function load($alias)
    {
        if (static::$ns && str_starts_with($alias, static::$ns)) {
            return 'facade';
        }
        return 'ok';
    }
}

$l = new StaticIn_Loader();
$r = $l->load('Bar');
if ($r !== 'ok') {
    Log::fatal('直接调用失败: '.var_export($r, true));
}

spl_autoload_register([$l, 'load']);
// trigger autoload of unknown class path via class_exists
$ok = class_exists('Some\\Missing\\ClassXYZ', true);
Log::info('static_in_instance_method 测试通过 autoload_ran');
