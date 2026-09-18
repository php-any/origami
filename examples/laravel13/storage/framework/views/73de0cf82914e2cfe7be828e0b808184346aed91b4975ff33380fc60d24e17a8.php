<?php
    use Illuminate\Contracts\Support\Htmlable;
    use Illuminate\Support\Arr;

    $brandName = filament()->getBrandName();
    $brandLogo = filament()->getBrandLogo();
    $brandLogoHeight = filament()->getBrandLogoHeight() ?? '1.5rem';
    $darkModeBrandLogo = filament()->getDarkModeBrandLogo();
    $hasDarkModeBrandLogo = filled($darkModeBrandLogo);

    $getLogoClasses = fn (bool $isDarkMode): string => Arr::toCssClasses([
        'fi-logo',
        'fi-logo-light' => $hasDarkModeBrandLogo && (! $isDarkMode),
        'fi-logo-dark' => $isDarkMode,
    ]);

    $logoStyles = 'height: ' . e($brandLogoHeight);
?>


            <?php $content = (function ($args) {
                return function ($logo, $isDarkMode = false) use ($args) {
                    extract($args, EXTR_SKIP);
                    ob_start(); ?>
        
    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($logo instanceof Htmlable): ?>
        <div
            <?php echo e($attributes
                    ->class([$getLogoClasses($isDarkMode)])
                    ->style([$logoStyles])); ?>

        >
            <?php echo e($logo); ?>

        </div>
    <?php elseif(filled($logo)): ?>
        <img
            alt="<?php echo e(__('filament-panels::layout.logo.alt', ['name' => $brandName])); ?>"
            src="<?php echo e($logo); ?>"
            <?php echo e($attributes
                    ->class([$getLogoClasses($isDarkMode)])
                    ->style([$logoStyles])); ?>

        />
    <?php else: ?>
        <div
            <?php echo e($attributes->class([
                    $getLogoClasses($isDarkMode),
                ])); ?>

        >
            <?php echo e($brandName); ?>

        </div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php return new \Illuminate\Support\HtmlString(ob_get_clean()); };
                })(get_defined_vars()); ?>
        

<?php echo e($content($brandLogo)); ?>


<?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasDarkModeBrandLogo): ?>
    <?php echo e($content($darkModeBrandLogo, isDarkMode: true)); ?>

<?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/logo.blade.php ENDPATH**/ ?>