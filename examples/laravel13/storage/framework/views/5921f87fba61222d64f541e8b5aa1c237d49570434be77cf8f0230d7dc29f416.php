<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'availableHeight' => null,
    'availableWidth' => null,
    'flip' => true,
    'maxHeight' => null,
    'offset' => 8,
    'placement' => null,
    'shift' => false,
    'size' => false,
    'sizePadding' => 16,
    'teleport' => false,
    'trigger' => null,
    'width' => null,
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
    'availableHeight' => null,
    'availableWidth' => null,
    'flip' => true,
    'maxHeight' => null,
    'offset' => 8,
    'placement' => null,
    'shift' => false,
    'size' => false,
    'sizePadding' => 16,
    'teleport' => false,
    'trigger' => null,
    'width' => null,
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Support\Enums\Width;

    $sizeConfig = collect([
        'availableHeight' => $availableHeight,
        'availableWidth' => $availableWidth,
        'padding' => $sizePadding,
    ])->filter()->toJson();

    if (is_string($width)) {
        $width = Width::tryFrom($width) ?? $width;
    }
?>

<div
    x-data="filamentDropdown"
    <?php echo e($attributes->class(['fi-dropdown'])); ?>

>
    <div
        x-on:keyup.enter="toggle($event)"
        x-on:keyup.space="toggle($event)"
        x-on:mousedown="if ($event.button === 0) toggle($event)"
        <?php echo e($trigger->attributes->class(['fi-dropdown-trigger'])); ?>

    >
        <?php echo e($trigger); ?>

    </div>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! \Filament\Support\is_slot_empty($slot)): ?>
        <div
            x-cloak
            x-float<?php echo e($placement ? ".placement.{$placement}" : ''); ?><?php echo e($size ? '.size' : ''); ?><?php echo e($flip ? '.flip' : ''); ?><?php echo e($shift ? '.shift' : ''); ?><?php echo e($teleport ? '.teleport' : ''); ?><?php echo e($offset ? '.offset' : ''); ?>="{ offset: <?php echo e($offset); ?>, <?php echo e($size ? ('size: ' . $sizeConfig) : ''); ?> }"
            x-ref="panel"
            x-transition:enter-start="fi-opacity-0"
            x-transition:leave-end="fi-opacity-0"
            <?php if($attributes->has('wire:key')): ?>
                wire:ignore.self
                <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($attributes->get('wire:key')).'.panel'; ?>wire:key="<?php echo e($attributes->get('wire:key')); ?>.panel"
            <?php endif; ?>
            class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                'fi-dropdown-panel',
                ($width instanceof Width) ? "fi-width-{$width->value}" : (is_string($width) ? $width : ''),
                'fi-scrollable' => $maxHeight || $size,
            ]); ?>"
            style="<?php echo \Illuminate\Support\Arr::toCssStyles([
                ('max-height: ' . e($maxHeight)) => $maxHeight,
            ]) ?>"
        >
            <?php echo e($slot); ?>

        </div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\support\resources\views/components/dropdown/index.blade.php ENDPATH**/ ?>