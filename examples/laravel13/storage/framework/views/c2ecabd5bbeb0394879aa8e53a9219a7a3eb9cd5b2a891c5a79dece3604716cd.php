<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames((['routing']));

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

foreach (array_filter((['routing']), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<div class="flex flex-col gap-3">
    <h2 class="text-lg font-semibold">Routing</h2>
    <div class="flex flex-col">
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__empty_1 = true; $__currentLoopData = $routing; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $key => $value): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); $__empty_1 = false; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
        <div class="flex max-w-full items-baseline gap-2 h-10 text-sm font-mono">
            <div class="uppercase text-neutral-500 dark:text-neutral-400 shrink-0"><?php echo e($key); ?></div>
            <div class="min-w-6 grow h-3 border-b-2 border-dotted border-neutral-300 dark:border-white/20"></div>
            <div class="truncate text-neutral-900 dark:text-white">
                <span data-tippy-content="<?php echo e($value); ?>">
                    <?php echo e($value); ?>

                </span>
            </div>
        </div>
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); if ($__empty_1): ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
        <?php if (isset($component)) { $__componentOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.empty-state','data' => ['message' => 'No routing context']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::empty-state'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['message' => 'No routing context']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56)): ?>
<?php $attributes = $__attributesOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56; ?>
<?php unset($__attributesOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56)): ?>
<?php $component = $__componentOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56; ?>
<?php unset($__componentOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56); ?>
<?php endif; ?>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    </div>
</div>
<?php /**PATH d:\gitcode.com\origami\examples\laravel13\vendor\laravel\framework\src\illuminate\foundation\providers/../resources/exceptions/renderer/components/routing.blade.php ENDPATH**/ ?>