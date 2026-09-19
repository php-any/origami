<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'active' => false,
    'ariaLabel' => null,
    'disabled' => false,
    'icon' => null,
    'iconAlias' => null,
    'label' => null,
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
    'active' => false,
    'ariaLabel' => null,
    'disabled' => false,
    'icon' => null,
    'iconAlias' => null,
    'label' => null,
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
    use Illuminate\Support\Number;
?>

<li
    <?php echo e($attributes->class([
            'fi-pagination-item',
            'fi-disabled' => $disabled,
            'fi-active' => $active,
        ])); ?>

>
    <button
        <?php if(filled($ariaLabel)): ?>
            aria-label="<?php echo e($ariaLabel); ?>"
        <?php endif; ?>
        <?php if($active): ?>
            aria-current="page"
        <?php endif; ?>
        <?php if($disabled): ?>
            aria-hidden="true"
        <?php endif; ?>
        <?php if($disabled): echo 'disabled'; endif; ?>
        type="button"
        class="fi-pagination-item-btn"
    >
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($icon || $iconAlias): ?>
            <?php echo e(\Filament\Support\generate_icon_html($icon, $iconAlias, attributes: (new ComponentAttributeBag)->merge(['aria-hidden' => 'true'], escape: false)->class([
                    'fi-pagination-item-icon',
                ]))); ?>

        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($label)): ?>
            <span class="fi-pagination-item-label">
                <?php echo e(is_numeric($label) ? Number::format($label) : ($label ?? '...')); ?>

            </span>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    </button>
</li>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\support\resources\views/components/pagination/item.blade.php ENDPATH**/ ?>