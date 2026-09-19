<?php

namespace tests\php;

/**
 * 对齐 Illuminate\View\Engines\PhpEngine：catch (Throwable) 后 throw $e 必须继续向外抛，
 * 不能落到 return ob_get_clean() 把空 HTML 交给 Livewire（RootTagMissing）。
 */

class PhpEngineRethrow_Engine
{
    public function evaluatePath()
    {
        $obLevel = ob_get_level();
        ob_start();
        echo 'PARTIAL';
        try {
            throw new \InvalidArgumentException('View [] not found.');
        } catch (\Throwable $e) {
            $this->handleViewException($e, $obLevel);
        }

        return ltrim(ob_get_clean());
    }

    protected function handleViewException($e, $obLevel)
    {
        while (ob_get_level() > $obLevel) {
            ob_end_clean();
        }
        throw $e;
    }
}

$engine = new PhpEngineRethrow_Engine();
$got = 'no-throw';
try {
    $out = $engine->evaluatePath();
    $got = 'returned:'.var_export($out, true);
} catch (\Throwable $e) {
    $got = 'caught:'.get_class($e).':'.$e->getMessage();
}

if (!str_starts_with($got, 'caught:') || !str_contains($got, 'View [] not found.')) {
    \Log::fatal('PhpEngine throw $e 应向外抛, 实际: '.$got);
}

\Log::info('phpengine_rethrow 测试通过');
