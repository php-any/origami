<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'applyAction',
    'form',
    'headingTag' => 'h3',
    'resetActionPosition' => null,
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
    'applyAction',
    'form',
    'headingTag' => 'h3',
    'resetActionPosition' => null,
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Support\View\ComponentAttributeBag;
    use Filament\Tables\Enums\FiltersResetActionPosition;

    $resetActionPosition ??= FiltersResetActionPosition::Header;
?>

<div <?php echo e($attributes->class(['fi-ta-filters'])); ?>>
    <div class="fi-ta-filters-header">
        <<?php echo e($headingTag); ?> class="fi-ta-filters-heading">
            <?php echo e(__('filament-tables::table.filters.heading')); ?>

        </<?php echo e($headingTag); ?>>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($resetActionPosition === FiltersResetActionPosition::Header): ?>
            <div>
                <?php if (isset($component)) { $__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.link','data' => ['attributes' => 
                        \Filament\Support\prepare_inherited_attributes(
                            new ComponentAttributeBag([
                                'color' => 'danger',
                                'tag' => 'button',
                                'wire:click' => 'resetTableFiltersForm',
                                'wire:loading.remove.delay.' . config('filament.livewire_loading_delay', 'default') => '',
                                'wire:target' => 'resetTableFiltersForm',
                            ])
                        )
                    ]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::link'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
                        \Filament\Support\prepare_inherited_attributes(
                            new ComponentAttributeBag([
                                'color' => 'danger',
                                'tag' => 'button',
                                'wire:click' => 'resetTableFiltersForm',
                                'wire:loading.remove.delay.' . config('filament.livewire_loading_delay', 'default') => '',
                                'wire:target' => 'resetTableFiltersForm',
                            ])
                        )
                    )]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                    <?php echo e(__('filament-tables::table.filters.actions.reset.label')); ?>

                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe)): ?>
<?php $attributes = $__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe; ?>
<?php unset($__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe); ?>
<?php endif; ?>
<?php if (isset($__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe)): ?>
<?php $component = $__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe; ?>
<?php unset($__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe); ?>
<?php endif; ?>
            </div>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    </div>

    <?php echo e($form); ?>


    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($applyAction->isVisible() || $resetActionPosition === FiltersResetActionPosition::Footer): ?>
        <div class="fi-ta-filters-actions-ctn">
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($applyAction->isVisible()): ?>
                <?php echo e($applyAction); ?>

            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($resetActionPosition === FiltersResetActionPosition::Footer): ?>
                <?php if (isset($component)) { $__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.button.index','data' => ['color' => 'danger','wire:click' => 'resetTableFiltersForm']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'danger','wire:click' => 'resetTableFiltersForm']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                    <?php echo e(__('filament-tables::table.filters.actions.reset.label')); ?>

                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef)): ?>
<?php $attributes = $__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef; ?>
<?php unset($__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef)): ?>
<?php $component = $__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef; ?>
<?php unset($__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef); ?>
<?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        </div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\tables\resources\views/components/filters.blade.php ENDPATH**/ ?>