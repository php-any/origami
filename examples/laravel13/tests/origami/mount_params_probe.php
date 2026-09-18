<?php
/**
 * 复现 Livewire __mountParamsContainer 注册与 new。
 * 匿名类放在 bootstrap 之后 include，避免解析期找不到 Livewire\Component。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

echo "livewire=".get_class(app('livewire'))."\n";
echo "component_loaded=".(class_exists(\Livewire\Component::class)?'y':'n')."\n";

require __DIR__.'/mount_params_probe_body.php';
