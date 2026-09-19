<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'applyAction',
    'columns' => null,
    'hasReorderableColumns',
    'hasToggleableColumns',
    'headingTag' => 'h3',
    'reorderAnimationDuration' => 300,
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
    'columns' => null,
    'hasReorderableColumns',
    'hasToggleableColumns',
    'headingTag' => 'h3',
    'reorderAnimationDuration' => 300,
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
    use Filament\Support\View\ComponentAttributeBag as FilamentComponentAttributeBag;
    use Filament\Tables\Enums\ColumnManagerResetActionPosition;

    $resetActionPosition ??= ColumnManagerResetActionPosition::Header;
?>

<div
    x-data="filamentTableColumnManager({
                columns: $wire.entangle('tableColumns'),
                isLive: <?php echo e($applyAction->isVisible() ? 'false' : 'true'); ?>,
            })"
    class="fi-ta-col-manager"
>
    <div class="fi-ta-col-manager-header">
        <<?php echo e($headingTag); ?> class="fi-ta-col-manager-heading">
            <?php echo e(__('filament-tables::table.column_manager.heading')); ?>

        </<?php echo e($headingTag); ?>>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($resetActionPosition === ColumnManagerResetActionPosition::Header): ?>
            <div>
                <?php if (isset($component)) { $__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.link','data' => ['attributes' => 
                        \Filament\Support\prepare_inherited_attributes(
                            new FilamentComponentAttributeBag([
                                'color' => 'danger',
                                'tag' => 'button',
                                'wire:click' => 'resetTableColumnManager',
                                'wire:loading.remove.delay.' . config('filament.livewire_loading_delay', 'default') => '',
                                'wire:target' => 'resetTableColumnManager',
                                'x-on:click' => 'resetDeferredColumns',
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
                            new FilamentComponentAttributeBag([
                                'color' => 'danger',
                                'tag' => 'button',
                                'wire:click' => 'resetTableColumnManager',
                                'wire:loading.remove.delay.' . config('filament.livewire_loading_delay', 'default') => '',
                                'wire:target' => 'resetTableColumnManager',
                                'x-on:click' => 'resetDeferredColumns',
                            ])
                        )
                    )]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                    <?php echo e(__('filament-tables::table.column_manager.actions.reset.label')); ?>

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

    <?php if (isset($component)) { $__componentOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.column-manager.content','data' => ['columns' => $columns,'hasReorderableColumns' => $hasReorderableColumns,'hasToggleableColumns' => $hasToggleableColumns,'reorderAnimationDuration' => $reorderAnimationDuration]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::column-manager.content'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columns),'has-reorderable-columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasReorderableColumns),'has-toggleable-columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasToggleableColumns),'reorder-animation-duration' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($reorderAnimationDuration)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d)): ?>
<?php $attributes = $__attributesOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d; ?>
<?php unset($__attributesOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d); ?>
<?php endif; ?>
<?php if (isset($__componentOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d)): ?>
<?php $component = $__componentOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d; ?>
<?php unset($__componentOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d); ?>
<?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($applyAction->isVisible() || $resetActionPosition === ColumnManagerResetActionPosition::Footer): ?>
        <div class="fi-ta-col-manager-actions-ctn">
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($applyAction->isVisible()): ?>
                <?php echo e($applyAction); ?>

            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($resetActionPosition === ColumnManagerResetActionPosition::Footer): ?>
                <?php if (isset($component)) { $__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.button.index','data' => ['color' => 'danger','wire:click' => 'resetTableColumnManager','xOn:click' => 'resetDeferredColumns']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'danger','wire:click' => 'resetTableColumnManager','x-on:click' => 'resetDeferredColumns']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                    <?php echo e(__('filament-tables::table.column_manager.actions.reset.label')); ?>

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
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\tables\resources\views/components/column-manager/index.blade.php ENDPATH**/ ?>