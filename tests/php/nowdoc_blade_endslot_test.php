<?php

namespace tests\php;

/**
 * 预处理器不得把 nowdoc/heredoc 里的 Blade @endslot/@endcomponent 改成 }。
 * Livewire 全页组件布局模板依赖这些指令交给 Blade 编译。
 */
$tpl = <<<'HTML'
<?php $layout->viewContext->mergeIntoNewEnvironment($__env); ?>

@component($layout->view, $layout->params)
    @slot($layout->slotOrSection)
        {!! $content !!}
    @endslot
@endcomponent
HTML;

if (!str_contains($tpl, '@endslot') || !str_contains($tpl, '@endcomponent')) {
    Log::fatal('nowdoc 中的 @endslot/@endcomponent 被破坏: ' . var_export($tpl, true));
}
if (str_contains($tpl, '<?php } ?>')) {
    Log::fatal('nowdoc 中的 Blade @end 被误换成花括号: ' . var_export($tpl, true));
}

Log::info('nowdoc Blade @endslot 保留测试通过');
