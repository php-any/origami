<?php

namespace tests\php;

/**
 * 模拟 Livewire Redirector::to 写入组件 store（不依赖 Filament）。
 */
class RedirectStore_Component
{
    public $store = [];

    public function redirect($url, $navigate = false)
    {
        $this->store['redirect'] = $url;
        if ($navigate) {
            $this->store['redirectUsingNavigate'] = true;
        }
    }
}

class RedirectStore_Redirector
{
    public $component;

    public function component($c)
    {
        $this->component = $c;
        return $this;
    }

    public function to($path)
    {
        $this->component->redirect($path);
        return $this;
    }
}

$comp = new RedirectStore_Component();
$redir = (new RedirectStore_Redirector())->component($comp);
$redir->to('/admin');
if (($comp->store['redirect'] ?? null) !== '/admin') {
    \Log::fatal('Redirector::to 未写入 store: '.var_export($comp->store, true));
}

// 闭包 !== null（EventBus finish 收集）
$fn = function ($x) { return $x; };
if (!($fn !== null)) {
    \Log::fatal('闭包 !== null 失败');
}

\Log::info('livewire_redirect_store 测试通过');
