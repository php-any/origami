<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'actions' => [],
    'actionsAlignment' => null,
    'breadcrumbs' => [],
    'heading' => null,
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
    'actions' => [],
    'actionsAlignment' => null,
    'breadcrumbs' => [],
    'heading' => null,
    'subheading' => null,
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Support\Facades\FilamentView;
    use Filament\View\PanelsRenderHook;
?>

<header
    <?php echo e($attributes->class([
            'fi-header',
            'fi-header-has-breadcrumbs' => $breadcrumbs,
        ])); ?>

>
    <div>
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($breadcrumbs): ?>
            <?php if (isset($component)) { $__componentOriginal3bde4a126f9ab6b0cc7b023c79de07bc95314ff71e6aff1141b8031dd5b83935 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal3bde4a126f9ab6b0cc7b023c79de07bc95314ff71e6aff1141b8031dd5b83935 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.breadcrumbs','data' => ['breadcrumbs' => $breadcrumbs]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::breadcrumbs'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['breadcrumbs' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($breadcrumbs)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal3bde4a126f9ab6b0cc7b023c79de07bc95314ff71e6aff1141b8031dd5b83935)): ?>
<?php $attributes = $__attributesOriginal3bde4a126f9ab6b0cc7b023c79de07bc95314ff71e6aff1141b8031dd5b83935; ?>
<?php unset($__attributesOriginal3bde4a126f9ab6b0cc7b023c79de07bc95314ff71e6aff1141b8031dd5b83935); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal3bde4a126f9ab6b0cc7b023c79de07bc95314ff71e6aff1141b8031dd5b83935)): ?>
<?php $component = $__componentOriginal3bde4a126f9ab6b0cc7b023c79de07bc95314ff71e6aff1141b8031dd5b83935; ?>
<?php unset($__componentOriginal3bde4a126f9ab6b0cc7b023c79de07bc95314ff71e6aff1141b8031dd5b83935); ?>
<?php endif; ?>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_HEADER_HEADING_BEFORE, scopes: $this->getRenderHookScopes())); ?>


        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($heading)): ?>
            <h1 class="fi-header-heading">
                <?php echo e($heading); ?>

            </h1>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_HEADER_HEADING_AFTER, scopes: $this->getRenderHookScopes())); ?>


        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($subheading)): ?>
            <p class="fi-header-subheading">
                <?php echo e($subheading); ?>

            </p>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    </div>

    <?php
        $beforeActions = FilamentView::renderHook(PanelsRenderHook::PAGE_HEADER_ACTIONS_BEFORE, scopes: $this->getRenderHookScopes());
        $afterActions = FilamentView::renderHook(PanelsRenderHook::PAGE_HEADER_ACTIONS_AFTER, scopes: $this->getRenderHookScopes());
    ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($beforeActions) || $actions || filled($afterActions)): ?>
        <div class="fi-header-actions-ctn">
            <?php echo e($beforeActions); ?>


            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($actions): ?>
                <?php if (isset($component)) { $__componentOriginalcf73f629bf156a266b3c5bc1df145b84a7d3f146e64b645132cc737b7fa86b30 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalcf73f629bf156a266b3c5bc1df145b84a7d3f146e64b645132cc737b7fa86b30 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.actions','data' => ['actions' => $actions,'alignment' => $actionsAlignment]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::actions'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['actions' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($actions),'alignment' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($actionsAlignment)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalcf73f629bf156a266b3c5bc1df145b84a7d3f146e64b645132cc737b7fa86b30)): ?>
<?php $attributes = $__attributesOriginalcf73f629bf156a266b3c5bc1df145b84a7d3f146e64b645132cc737b7fa86b30; ?>
<?php unset($__attributesOriginalcf73f629bf156a266b3c5bc1df145b84a7d3f146e64b645132cc737b7fa86b30); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalcf73f629bf156a266b3c5bc1df145b84a7d3f146e64b645132cc737b7fa86b30)): ?>
<?php $component = $__componentOriginalcf73f629bf156a266b3c5bc1df145b84a7d3f146e64b645132cc737b7fa86b30; ?>
<?php unset($__componentOriginalcf73f629bf156a266b3c5bc1df145b84a7d3f146e64b645132cc737b7fa86b30); ?>
<?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php echo e($afterActions); ?>

        </div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
</header>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/header/index.blade.php ENDPATH**/ ?>