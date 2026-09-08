<?php

namespace tests\php;

/**
 * method_exists 对继承/子类方法应返回 true（Livewire ComponentHook::callBoot 依赖）。
 */
abstract class MethodExists_HookBase
{
    function callBoot()
    {
        return method_exists($this, 'boot');
    }
}

class MethodExists_SupportRedirects extends MethodExists_HookBase
{
    public function boot()
    {
        return 'booted';
    }
}

$h = new MethodExists_SupportRedirects();
if ($h->callBoot() !== true) {
    Log::fatal('method_exists($this, boot) 在子类上应为 true, got='.var_export($h->callBoot(), true));
}

if (!method_exists($h, 'boot')) {
    Log::fatal('method_exists 直接调用失败');
}

if (!method_exists(MethodExists_SupportRedirects::class, 'boot')) {
    Log::fatal('method_exists(class, boot) 失败');
}

Log::info('method_exists_boot 测试通过');
