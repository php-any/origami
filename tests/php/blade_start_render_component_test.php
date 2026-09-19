<?php

namespace tests\php;

/**
 * 对齐 Factory::startComponent / renderComponent：if (ob_start()) 后入栈，
 * render 时 array_pop 必须拿到同一 $view。
 */
class BladeComp_Env
{
    public $componentStack = [];
    public $componentData = [];
    public $slots = [];

    public function startComponent($view, array $data = [])
    {
        if (ob_start()) {
            $this->componentStack[] = $view;
            $this->componentData[$this->currentComponent()] = $data;
            $this->slots[$this->currentComponent()] = [];
        }
    }

    public function renderComponent()
    {
        $view = array_pop($this->componentStack);
        $defaultSlot = trim(ob_get_clean());
        $idx = count($this->componentStack);
        if (!array_key_exists($idx, $this->slots)) {
            \Log::fatal('slots 缺少 key '.$idx.' keys='.json_encode(array_keys($this->slots)).' view='.var_export($view, true));
        }
        $data = $this->componentData[$idx];

        return [$view, $data, $defaultSlot];
    }

    protected function currentComponent()
    {
        return count($this->componentStack) - 1;
    }
}

$e = new BladeComp_Env();
$e->startComponent('outer', ['o' => 1]);
echo 'OUTER_START';
$e->startComponent('inner', ['i' => 2]);
echo 'INNER';
[$innerView, $innerData, $innerSlot] = $e->renderComponent();
if ($innerView !== 'inner' || $innerSlot !== 'INNER') {
    \Log::fatal('inner 渲染错误 view='.var_export($innerView, true).' slot='.var_export($innerSlot, true));
}
echo 'BETWEEN';
[$outerView, $outerData, $outerSlot] = $e->renderComponent();
if ($outerView !== 'outer') {
    \Log::fatal('outer view 应为 outer, 实际: '.var_export($outerView, true));
}
if (!str_contains($outerSlot, 'BETWEEN')) {
    \Log::fatal('outer slot 缺少 BETWEEN: '.var_export($outerSlot, true));
}
if ($e->componentStack !== []) {
    \Log::fatal('结束后 stack 应空: '.var_export($e->componentStack, true));
}

$e2 = new BladeComp_Env();
for ($i = 0; $i < 12; $i++) {
    $e2->startComponent('v'.$i, ['n' => $i]);
    echo 'L'.$i;
}
for ($i = 11; $i >= 0; $i--) {
    [$v, $d, $s] = $e2->renderComponent();
    if ($v !== 'v'.$i) {
        \Log::fatal('深度 '.$i.' view='.var_export($v, true));
    }
}

\Log::info('blade_start_render_component 测试通过');
