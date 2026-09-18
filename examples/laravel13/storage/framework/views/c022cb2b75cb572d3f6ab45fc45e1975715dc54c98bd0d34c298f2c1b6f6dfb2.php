<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'heading' => null,
    'logo' => true,
    'subheading' => null,
]));

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

foreach (array_filter(([
    'heading' => null,
    'logo' => true,
    'subheading' => null,
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<header class="fi-simple-header">
    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($logo): ?>
        <?php if (isset($component)) { $__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.logo','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::logo'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8)): ?>
<?php $attributes = $__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8; ?>
<?php unset($__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8)): ?>
<?php $component = $__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8; ?>
<?php unset($__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8); ?>
<?php endif; ?>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($heading)): ?>
        <h1 class="fi-simple-header-heading">
            <?php echo e($heading); ?>

        </h1>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($subheading)): ?>
        <p class="fi-simple-header-subheading">
            <?php echo e($subheading); ?>

        </p>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
</header>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/header/simple.blade.php ENDPATH**/ ?>