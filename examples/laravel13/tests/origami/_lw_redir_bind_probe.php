<?php

require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->bootstrapWith([
    Illuminate\Foundation\Bootstrap\LoadEnvironmentVariables::class,
    Illuminate\Foundation\Bootstrap\LoadConfiguration::class,
    Illuminate\Foundation\Bootstrap\HandleExceptions::class,
    Illuminate\Foundation\Bootstrap\RegisterFacades::class,
    Illuminate\Foundation\Bootstrap\SetRequestForConsole::class,
    Illuminate\Foundation\Bootstrap\RegisterProviders::class,
    Illuminate\Foundation\Bootstrap\BootProviders::class,
]);

function p($m) { file_put_contents(__DIR__.'/_login_progress.txt', $m."\n", FILE_APPEND); }
@unlink(__DIR__.'/_login_progress.txt');

p('1');
$fakeComp = new \App\Livewire\Admin\Login();
p('2');
$hook = new \Livewire\Features\SupportRedirects\SupportRedirects();
p('3');
$hook->setComponent($fakeComp);
p('4');
$hook->callBoot();
p('5 after_boot='.get_class(app('redirect')));
p('6 is_lw='.(app('redirect') instanceof \Livewire\Features\SupportRedirects\Redirector ? 'yes' : 'no'));

$redir = app('redirect');
p('7');
try {
    $ret = $redir->route('admin.dashboard');
    p('8 ret='.get_class($ret));
    p('9 store='.var_export(\Livewire\store($fakeComp)->get('redirect'), true));
} catch (Throwable $e) {
    p('8 ERR '.$e->getMessage());
}
p('done');
