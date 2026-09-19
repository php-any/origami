<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'breadcrumbs' => [],
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
    'breadcrumbs' => [],
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Support\Icons\Heroicon;
    use Filament\Support\View\ComponentAttributeBag as FilamentComponentAttributeBag;
    use Filament\Support\View\SupportIconAlias;

    use function Filament\Support\generate_icon_html;
?>

<nav
    <?php echo e($attributes
            ->merge(['aria-label' => __('filament::components/breadcrumbs.label')], escape: true)
            ->class(['fi-breadcrumbs'])); ?>

>
    <ol class="fi-breadcrumbs-list">
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $breadcrumbs; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $url => $label): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
            <li class="fi-breadcrumbs-item">
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $loop->first): ?>
                    <?php echo e(generate_icon_html(Heroicon::ChevronRight, alias: SupportIconAlias::BREADCRUMBS_SEPARATOR, attributes: (new FilamentComponentAttributeBag)->merge(['aria-hidden' => 'true'], escape: false)->class([
                            'fi-breadcrumbs-item-separator fi-ltr',
                        ]))); ?>


                    <?php echo e(generate_icon_html(Heroicon::ChevronLeft, alias: SupportIconAlias::BREADCRUMBS_SEPARATOR_RTL, attributes: (new FilamentComponentAttributeBag)->merge(['aria-hidden' => 'true'], escape: false)->class([
                            'fi-breadcrumbs-item-separator fi-rtl',
                        ]))); ?>

                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(is_int($url)): ?>
                    <span
                        <?php if($loop->last): ?> aria-current="page" <?php endif; ?>
                        class="fi-breadcrumbs-item-label"
                    >
                        <?php echo e($label); ?>

                    </span>
                <?php else: ?>
                    <a
                        <?php echo e(\Filament\Support\generate_href_html($url)); ?>

                        <?php if($loop->last): ?> aria-current="page" <?php endif; ?>
                        class="fi-breadcrumbs-item-label"
                    >
                        <?php echo e($label); ?>

                    </a>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
            </li>
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
    </ol>
</nav>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\support\resources\views/components/breadcrumbs.blade.php ENDPATH**/ ?>