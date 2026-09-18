<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames((['frames']));

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

foreach (array_filter((['frames']), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php use \Illuminate\Support\Str; ?>

<div
    x-data="{ expanded: false }"
    class="group rounded-lg border border-neutral-200 dark:border-white/5"
    :class="{
        'bg-white dark:bg-white/5 shadow-xs': expanded,
        'border-dashed border-neutral-300 bg-neutral-50 opacity-90 dark:border-white/10 dark:bg-white/1': !expanded,
    }"
>
    <div
        class="flex h-11 cursor-pointer items-center gap-3 rounded-lg pr-2.5 pl-4 hover:bg-white/50 dark:hover:bg-white/2"
        @click="expanded = !expanded"
    >
        <?php if (isset($component)) { $__componentOriginal53b4a8a998443432c7c8448d360cd9769ae63aa00bd95d96f1a842095c7fe410 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal53b4a8a998443432c7c8448d360cd9769ae63aa00bd95d96f1a842095c7fe410 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.folder','data' => ['class' => 'w-3 h-3 text-neutral-400','xShow' => '!expanded','xCloak' => true]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.folder'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-3 h-3 text-neutral-400','x-show' => '!expanded','x-cloak' => true]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal53b4a8a998443432c7c8448d360cd9769ae63aa00bd95d96f1a842095c7fe410)): ?>
<?php $attributes = $__attributesOriginal53b4a8a998443432c7c8448d360cd9769ae63aa00bd95d96f1a842095c7fe410; ?>
<?php unset($__attributesOriginal53b4a8a998443432c7c8448d360cd9769ae63aa00bd95d96f1a842095c7fe410); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal53b4a8a998443432c7c8448d360cd9769ae63aa00bd95d96f1a842095c7fe410)): ?>
<?php $component = $__componentOriginal53b4a8a998443432c7c8448d360cd9769ae63aa00bd95d96f1a842095c7fe410; ?>
<?php unset($__componentOriginal53b4a8a998443432c7c8448d360cd9769ae63aa00bd95d96f1a842095c7fe410); ?>
<?php endif; ?>
        <?php if (isset($component)) { $__componentOriginal31a83bb50c4fc35c18d6abaefe5a76ab04b75276eb3a6b59d39d77ca9087920a = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal31a83bb50c4fc35c18d6abaefe5a76ab04b75276eb3a6b59d39d77ca9087920a = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.folder-open','data' => ['class' => 'w-3 h-3 text-blue-500 dark:text-emerald-500','xShow' => 'expanded']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.folder-open'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-3 h-3 text-blue-500 dark:text-emerald-500','x-show' => 'expanded']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal31a83bb50c4fc35c18d6abaefe5a76ab04b75276eb3a6b59d39d77ca9087920a)): ?>
<?php $attributes = $__attributesOriginal31a83bb50c4fc35c18d6abaefe5a76ab04b75276eb3a6b59d39d77ca9087920a; ?>
<?php unset($__attributesOriginal31a83bb50c4fc35c18d6abaefe5a76ab04b75276eb3a6b59d39d77ca9087920a); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal31a83bb50c4fc35c18d6abaefe5a76ab04b75276eb3a6b59d39d77ca9087920a)): ?>
<?php $component = $__componentOriginal31a83bb50c4fc35c18d6abaefe5a76ab04b75276eb3a6b59d39d77ca9087920a; ?>
<?php unset($__componentOriginal31a83bb50c4fc35c18d6abaefe5a76ab04b75276eb3a6b59d39d77ca9087920a); ?>
<?php endif; ?>

        <div class="flex-1 font-mono text-xs leading-3 text-neutral-900 dark:text-neutral-400">
            <?php echo e(count($frames)); ?> vendor <?php echo e(Str::plural('frame', count($frames))); ?>

        </div>

        <button
            x-cloak
            type="button"
            class="flex h-6 w-6 cursor-pointer items-center justify-center rounded-md dark:border dark:border-white/8 group-hover:text-blue-500 group-hover:dark:text-emerald-500"
            :class="{
                'text-blue-500 dark:text-emerald-500 dark:bg-white/5': expanded,
                'text-neutral-500 dark:text-neutral-500 dark:bg-white/3': !expanded,
            }"
        >
            <?php if (isset($component)) { $__componentOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.chevrons-down-up','data' => ['xShow' => 'expanded']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.chevrons-down-up'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-show' => 'expanded']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964)): ?>
<?php $attributes = $__attributesOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964; ?>
<?php unset($__attributesOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964)): ?>
<?php $component = $__componentOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964; ?>
<?php unset($__componentOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964); ?>
<?php endif; ?>
            <?php if (isset($component)) { $__componentOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.chevrons-up-down','data' => ['xShow' => '!expanded']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.chevrons-up-down'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-show' => '!expanded']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a)): ?>
<?php $attributes = $__attributesOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a; ?>
<?php unset($__attributesOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a)): ?>
<?php $component = $__componentOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a; ?>
<?php unset($__componentOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a); ?>
<?php endif; ?>
        </button>
    </div>

    <div x-cloak class="flex flex-col rounded-b-lg divide-y divide-neutral-200 border-t border-neutral-200 dark:divide-white/5 dark:border-white/5" x-show="expanded">
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $frames; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $frame): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
            <div class="flex flex-col divide-y divide-neutral-200 dark:divide-white/5">
                <?php if (isset($component)) { $__componentOriginalfa71dcd6057e75ef54fe50611fc63e1701ea84fa9498e7a1860470c30163e2e2 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalfa71dcd6057e75ef54fe50611fc63e1701ea84fa9498e7a1860470c30163e2e2 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.vendor-frame','data' => ['frame' => $frame]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::vendor-frame'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['frame' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($frame)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalfa71dcd6057e75ef54fe50611fc63e1701ea84fa9498e7a1860470c30163e2e2)): ?>
<?php $attributes = $__attributesOriginalfa71dcd6057e75ef54fe50611fc63e1701ea84fa9498e7a1860470c30163e2e2; ?>
<?php unset($__attributesOriginalfa71dcd6057e75ef54fe50611fc63e1701ea84fa9498e7a1860470c30163e2e2); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalfa71dcd6057e75ef54fe50611fc63e1701ea84fa9498e7a1860470c30163e2e2)): ?>
<?php $component = $__componentOriginalfa71dcd6057e75ef54fe50611fc63e1701ea84fa9498e7a1860470c30163e2e2; ?>
<?php unset($__componentOriginalfa71dcd6057e75ef54fe50611fc63e1701ea84fa9498e7a1860470c30163e2e2); ?>
<?php endif; ?>
            </div>
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
    </div>
</div>
<?php /**PATH d:\gitcode.com\origami\examples\laravel13\vendor\laravel\framework\src\illuminate\foundation\providers/../resources/exceptions/renderer/components/vendor-frames.blade.php ENDPATH**/ ?>