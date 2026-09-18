<?php
/**
 * 在 Livewire update 路径上检查 redirect 绑定与 store。
 */
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Filament\Facades\Filament;
use Filament\Auth\Pages\Login;
use Illuminate\Contracts\Support\Responsable;
use Livewire\Mechanisms\HandleComponents\HandleComponents;
use Livewire\Features\SupportFileDownloads\SupportFileDownloads;
use Livewire\Features\SupportRedirects\Redirector as LwRedirector;

Filament::setCurrentPanel(Filament::getPanel('admin'));

// 模拟 hydrate 后的 boot 绑定
$html = app('livewire')->mount(Login::class);
echo "mounted\n";

// 从 livewire 拿到当前组件较难；直接 new + boot hooks
$component = app(Login::class);
$hook = new \Livewire\Features\SupportRedirects\SupportRedirects();
$hook->setComponent($component);
$hook->boot();

echo "redirect_after_boot=".get_class(app('redirect'))."\n";
echo "is_lw=".(app('redirect') instanceof LwRedirector ? 'yes':'no')."\n";

$component->form->fill(['email'=>'admin@example.com','password'=>'password','remember'=>false]);
$ret = $component->authenticate();
echo "ret_class=".get_class($ret)."\n";
echo "ret_responsable=".($ret instanceof Responsable ? 'yes':'no')."\n";

$fd = new SupportFileDownloads();
$fd->setComponent($component);
$finish = $fd->call();
$after = $finish($ret);
echo "after_finish_type=".gettype($after).(is_object($after)?(' '.get_class($after)):'')."\n";
echo "store=".var_export(\Livewire\store($component)->get('redirect'), true)."\n";
echo "redirect_now=".get_class(app('redirect'))."\n";

echo "DONE\n";
