<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'alt' => '',
    'circular' => true,
    'size' => 'md',
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
    'alt' => '',
    'circular' => true,
    'size' => 'md',
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<img
    alt="<?php echo e($alt); ?>"
    <?php echo e($attributes
            ->class([
                'fi-avatar',
                'fi-circular' => $circular,
                match ($size) {
                    'sm', 'md', 'lg' => "fi-size-{$size}",
                    default => $size,
                },
            ])); ?>

/>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\support\resources\views/components/avatar.blade.php ENDPATH**/ ?>