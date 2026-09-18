<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames((['frame']));

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

foreach (array_filter((['frame']), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<div class="grid gap-3 p-4 bg-neutral-50 dark:bg-transparent overflow-x-auto rounded-lg">
    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($frame->previous()): ?>
        <div class="flex">
            <?php if (isset($component)) { $__componentOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.formatted-source','data' => ['frame' => $frame,'className' => 'text-xs']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::formatted-source'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['frame' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($frame),'className' => 'text-xs']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536)): ?>
<?php $attributes = $__attributesOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536; ?>
<?php unset($__attributesOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536)): ?>
<?php $component = $__componentOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536; ?>
<?php unset($__componentOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536); ?>
<?php endif; ?>
        </div>
    <?php else: ?>
        <span class="font-mono text-xs leading-3 text-neutral-500">Entrypoint</span>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if (isset($component)) { $__componentOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.file-with-line','data' => ['frame' => $frame,'class' => 'text-xs']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::file-with-line'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['frame' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($frame),'class' => 'text-xs']); ?>
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
</div>
<?php /**PATH d:\gitcode.com\origami\examples\laravel13\vendor\laravel\framework\src\illuminate\foundation\providers/../resources/exceptions/renderer/components/vendor-frame.blade.php ENDPATH**/ ?>