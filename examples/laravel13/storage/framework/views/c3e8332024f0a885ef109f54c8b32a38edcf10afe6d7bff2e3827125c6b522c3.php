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

<div class="flex flex-col gap-2.5 bg-neutral-50 dark:bg-white/1 border border-neutral-200 dark:border-neutral-800 rounded-xl p-2.5 shadow-xs">
    <div class="flex items-center gap-2.5 p-2">
        <div class="bg-white dark:bg-neutral-800 border border-neutral-200 dark:border-white/5 rounded-md w-6 h-6 flex items-center justify-center p-1">
            <?php if (isset($component)) { $__componentOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.alert','data' => ['class' => 'w-2.5 h-2.5 text-blue-500 dark:text-emerald-500']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.alert'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-2.5 h-2.5 text-blue-500 dark:text-emerald-500']); ?>
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
        </div>
        <h3 class="text-base font-semibold text-neutral-900 dark:text-white">Exception trace</h3>
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($exception->previousExceptions()->isNotEmpty()): ?>
            <a href="#previous-exceptions" class="ml-auto text-sm text-neutral-500 dark:text-neutral-400 hover:text-blue-500 dark:hover:text-emerald-500 transition-colors">
                <?php echo e($exception->previousExceptions()->count()); ?> previous <?php echo e(\Illuminate\Support\Str::plural('exception', $exception->previousExceptions()->count())); ?>

            </a>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    </div>

    <div class="flex flex-col gap-1.5">
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $exception->frameGroups(); $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $group): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($group['is_vendor']): ?>
                <?php if (isset($component)) { $__componentOriginal64a76438cd033bb8ffab5197e88cc5ca78ab4d7ad1fa53727edfbefd12cf84f7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal64a76438cd033bb8ffab5197e88cc5ca78ab4d7ad1fa53727edfbefd12cf84f7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.vendor-frames','data' => ['frames' => $group['frames']]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::vendor-frames'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['frames' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($group['frames'])]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal64a76438cd033bb8ffab5197e88cc5ca78ab4d7ad1fa53727edfbefd12cf84f7)): ?>
<?php $attributes = $__attributesOriginal64a76438cd033bb8ffab5197e88cc5ca78ab4d7ad1fa53727edfbefd12cf84f7; ?>
<?php unset($__attributesOriginal64a76438cd033bb8ffab5197e88cc5ca78ab4d7ad1fa53727edfbefd12cf84f7); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal64a76438cd033bb8ffab5197e88cc5ca78ab4d7ad1fa53727edfbefd12cf84f7)): ?>
<?php $component = $__componentOriginal64a76438cd033bb8ffab5197e88cc5ca78ab4d7ad1fa53727edfbefd12cf84f7; ?>
<?php unset($__componentOriginal64a76438cd033bb8ffab5197e88cc5ca78ab4d7ad1fa53727edfbefd12cf84f7); ?>
<?php endif; ?>
            <?php else: ?>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $group['frames']; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $frame): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                    <?php if (isset($component)) { $__componentOriginalbb43b3e1ba62815657d223f89b0fc83d81cd40ba4053133151a492968015c08f = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalbb43b3e1ba62815657d223f89b0fc83d81cd40ba4053133151a492968015c08f = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.frame','data' => ['frame' => $frame]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::frame'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['frame' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($frame)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalbb43b3e1ba62815657d223f89b0fc83d81cd40ba4053133151a492968015c08f)): ?>
<?php $attributes = $__attributesOriginalbb43b3e1ba62815657d223f89b0fc83d81cd40ba4053133151a492968015c08f; ?>
<?php unset($__attributesOriginalbb43b3e1ba62815657d223f89b0fc83d81cd40ba4053133151a492968015c08f); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalbb43b3e1ba62815657d223f89b0fc83d81cd40ba4053133151a492968015c08f)): ?>
<?php $component = $__componentOriginalbb43b3e1ba62815657d223f89b0fc83d81cd40ba4053133151a492968015c08f; ?>
<?php unset($__componentOriginalbb43b3e1ba62815657d223f89b0fc83d81cd40ba4053133151a492968015c08f); ?>
<?php endif; ?>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
    </div>
</div>
<?php /**PATH d:\gitcode.com\origami\examples\laravel13\vendor\laravel\framework\src\illuminate\foundation\providers/../resources/exceptions/renderer/components/trace.blade.php ENDPATH**/ ?>