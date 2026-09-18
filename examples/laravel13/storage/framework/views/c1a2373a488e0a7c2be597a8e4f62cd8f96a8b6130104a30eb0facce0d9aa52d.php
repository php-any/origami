<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames((['exception']));

foreach ($attributes->all() as $__key => $__value) {
    if (in_array($__key, $__propNames)) {
        $$__key = $$__key ?? $__value;
    } else {
        $__newAttributes[$__key] = $__value;
    }
}

$attributes = new \Illuminate\View\ComponentAttributeBag($__newAttributes);

unset($__propNames);
unset($__newAttributes);

foreach (array_filter((['exception']), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<div class="flex flex-col pt-8 sm:pt-16 overflow-x-auto">
    <div class="flex flex-col gap-5 mb-8">
        <h1 class="text-3xl font-semibold text-neutral-950 dark:text-white"><?php echo e($exception->class()); ?></h1>
        <?php if (isset($component)) { $__componentOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.file-with-line','data' => ['frame' => $exception->frames()->first(),'class' => '-mt-3 text-xs']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::file-with-line'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['frame' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception->frames()->first()),'class' => '-mt-3 text-xs']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6)): ?>
<?php $attributes = $__attributesOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6; ?>
<?php unset($__attributesOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6)): ?>
<?php $component = $__componentOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6; ?>
<?php unset($__componentOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6); ?>
<?php endif; ?>
        <p class="text-xl font-light text-neutral-800 dark:text-neutral-300">
            <?php echo e($exception->message()); ?>

        </p>
    </div>

    <div class="flex items-start gap-2 mb-8 sm:mb-16">
        <div class="bg-white dark:bg-white/[3%] border border-neutral-200 dark:border-white/10 divide-x divide-neutral-200 dark:divide-white/10 rounded-md shadow-xs flex items-center gap-0.5">
            <div class="flex items-center gap-1.5 h-6 px-[6px] font-mono text-[13px]">
                <span class="text-neutral-400 dark:text-neutral-500">LARAVEL</span>
                <span class="text-neutral-500 dark:text-neutral-300"><?php echo e(app()->version()); ?></span>
            </div>
            <div class="flex items-center gap-1.5 h-6 px-[6px] font-mono text-[13px]">
                <span class="text-neutral-400 dark:text-neutral-500">PHP</span>
                <span class="text-neutral-500 dark:text-neutral-300"><?php echo e(PHP_VERSION); ?></span>
            </div>
        </div>
        <?php if (isset($component)) { $__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.badge','data' => ['type' => 'error']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::badge'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['type' => 'error']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

            <?php if (isset($component)) { $__componentOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.alert','data' => ['class' => 'w-2.5 h-2.5']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.alert'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-2.5 h-2.5']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a)): ?>
<?php $attributes = $__attributesOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a; ?>
<?php unset($__attributesOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a)): ?>
<?php $component = $__componentOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a; ?>
<?php unset($__componentOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a); ?>
<?php endif; ?>
            UNHANDLED
         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad)): ?>
<?php $attributes = $__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad; ?>
<?php unset($__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad)): ?>
<?php $component = $__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad; ?>
<?php unset($__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad); ?>
<?php endif; ?>
        <?php if (isset($component)) { $__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.badge','data' => ['type' => 'error','variant' => 'solid']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::badge'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['type' => 'error','variant' => 'solid']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

            CODE <?php echo e($exception->code()); ?>

         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad)): ?>
<?php $attributes = $__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad; ?>
<?php unset($__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad)): ?>
<?php $component = $__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad; ?>
<?php unset($__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad); ?>
<?php endif; ?>
    </div>

    <?php if (isset($component)) { $__componentOriginal4c21687d9ffa3e3331d7ea5566406ec712772d14abb2fe52ac150ea6b9a1cdb8 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal4c21687d9ffa3e3331d7ea5566406ec712772d14abb2fe52ac150ea6b9a1cdb8 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.request-url','data' => ['exception' => $exception,'request' => $exception->request(),'class' => 'relative z-50']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::request-url'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['exception' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception),'request' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception->request()),'class' => 'relative z-50']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal4c21687d9ffa3e3331d7ea5566406ec712772d14abb2fe52ac150ea6b9a1cdb8)): ?>
<?php $attributes = $__attributesOriginal4c21687d9ffa3e3331d7ea5566406ec712772d14abb2fe52ac150ea6b9a1cdb8; ?>
<?php unset($__attributesOriginal4c21687d9ffa3e3331d7ea5566406ec712772d14abb2fe52ac150ea6b9a1cdb8); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal4c21687d9ffa3e3331d7ea5566406ec712772d14abb2fe52ac150ea6b9a1cdb8)): ?>
<?php $component = $__componentOriginal4c21687d9ffa3e3331d7ea5566406ec712772d14abb2fe52ac150ea6b9a1cdb8; ?>
<?php unset($__componentOriginal4c21687d9ffa3e3331d7ea5566406ec712772d14abb2fe52ac150ea6b9a1cdb8); ?>
<?php endif; ?>
</div>
<?php /**PATH d:\gitcode.com\origami\examples\laravel13\vendor\laravel\framework\src\illuminate\foundation\providers/../resources/exceptions/renderer/components/header.blade.php ENDPATH**/ ?>