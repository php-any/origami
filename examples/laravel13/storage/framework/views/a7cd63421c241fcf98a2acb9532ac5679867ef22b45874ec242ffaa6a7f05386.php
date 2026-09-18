<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames((['method']));

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

foreach (array_filter((['method']), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
$type = match ($method) {
    'GET', 'OPTIONS', 'ANY' => 'default',
    'POST' => 'success',
    'PUT', 'PATCH' => 'primary',
    'DELETE' => 'error',
    default => 'default',
};
?>

<?php if (isset($component)) { $__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.badge','data' => ['type' => ''.e($type).'']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::badge'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['type' => ''.e($type).'']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

    <?php if (isset($component)) { $__componentOriginal902141d19c06a21e11042ec8558dbd68ad3a3a226308e0fec0f158c3f6906a03 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal902141d19c06a21e11042ec8558dbd68ad3a3a226308e0fec0f158c3f6906a03 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.globe','data' => ['class' => 'w-2.5 h-2.5']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.globe'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-2.5 h-2.5']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal902141d19c06a21e11042ec8558dbd68ad3a3a226308e0fec0f158c3f6906a03)): ?>
<?php $attributes = $__attributesOriginal902141d19c06a21e11042ec8558dbd68ad3a3a226308e0fec0f158c3f6906a03; ?>
<?php unset($__attributesOriginal902141d19c06a21e11042ec8558dbd68ad3a3a226308e0fec0f158c3f6906a03); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal902141d19c06a21e11042ec8558dbd68ad3a3a226308e0fec0f158c3f6906a03)): ?>
<?php $component = $__componentOriginal902141d19c06a21e11042ec8558dbd68ad3a3a226308e0fec0f158c3f6906a03; ?>
<?php unset($__componentOriginal902141d19c06a21e11042ec8558dbd68ad3a3a226308e0fec0f158c3f6906a03); ?>
<?php endif; ?>
    <?php echo e($method); ?>

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
<?php /**PATH d:\gitcode.com\origami\examples\laravel13\vendor\laravel\framework\src\illuminate\foundation\providers/../resources/exceptions/renderer/components/http-method.blade.php ENDPATH**/ ?>