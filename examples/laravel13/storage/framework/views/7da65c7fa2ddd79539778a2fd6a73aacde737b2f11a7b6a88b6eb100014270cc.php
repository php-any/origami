<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'livewire' => null,
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
    'livewire' => null,
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Livewire\Notifications;
    use Filament\Support\Facades\FilamentView;
    use Filament\View\PanelsRenderHook;

    $renderHookScopes = $livewire?->getRenderHookScopes();
?>

<!DOCTYPE html>
<html
    lang="<?php echo e(str_replace('_', '-', app()->getLocale())); ?>"
    dir="<?php echo e(__('filament-panels::layout.direction') ?? 'ltr'); ?>"
    class="<?php echo \Illuminate\Support\Arr::toCssClasses([
        'fi',
        'dark' => filament()->hasDarkMode() && filament()->hasDarkModeForced(),
    ]); ?>"
>
    <head>
        <?php echo e(FilamentView::renderHook(PanelsRenderHook::HEAD_START, scopes: $renderHookScopes)); ?>


        <meta charset="utf-8" />
        <meta name="csrf-token" content="<?php echo e(csrf_token()); ?>" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($favicon = filament()->getFavicon()): ?>
            <link rel="icon" href="<?php echo e($favicon); ?>" />
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php
            $title = trim(strip_tags($livewire?->getTitle() ?? ''));
            $brandName = trim(strip_tags(filament()->getBrandName()));
        ?>

        <title>
            <?php echo e(filled($title) ? $title : null); ?>

            <?php echo e(filled($brandName) && filled($title) ? ' - ' : null); ?>

            <?php echo e(filled($brandName) ? $brandName : null); ?>

        </title>

        <?php echo e(FilamentView::renderHook(PanelsRenderHook::STYLES_BEFORE, scopes: $renderHookScopes)); ?>


        <style>
            [x-cloak=''],
            [x-cloak='x-cloak'],
            [x-cloak='1'] {
                display: none !important;
            }

            [x-cloak='inline-flex'] {
                display: inline-flex !important;
            }

            @media (max-width: 1023px) {
                [x-cloak='-lg'] {
                    display: none !important;
                }
            }

            @media (min-width: 1024px) {
                [x-cloak='lg'] {
                    display: none !important;
                }
            }
        </style>

        <?php echo \Filament\Support\Facades\FilamentAsset::renderStyles() ?>

        <?php echo e(filament()->getTheme()->getHtml()); ?>

        <?php echo e(filament()->getFontPreloadHtml()); ?>

        <?php echo e(filament()->getMonoFontPreloadHtml()); ?>

        <?php echo e(filament()->getSerifFontPreloadHtml()); ?>

        <?php echo e(filament()->getFontHtml()); ?>

        <?php echo e(filament()->getMonoFontHtml()); ?>

        <?php echo e(filament()->getSerifFontHtml()); ?>


        <style>
            :root {
                --font-family: '<?php echo filament()->getFontFamily(); ?>';
                --mono-font-family: '<?php echo filament()->getMonoFontFamily(); ?>';
                --serif-font-family: '<?php echo filament()->getSerifFontFamily(); ?>';
                --sidebar-width: <?php echo e(filament()->getSidebarWidth()); ?>;
                --collapsed-sidebar-width: <?php echo e(filament()->getCollapsedSidebarWidth()); ?>;
                --default-theme-mode: <?php echo e(filament()->getDefaultThemeMode()->value); ?>;
            }

            html.fi {
                --livewire-progress-bar-color: var(--primary-500);
            }
        </style>

        <?php echo $__env->yieldPushContent('styles'); ?>

        <?php echo e(FilamentView::renderHook(PanelsRenderHook::STYLES_AFTER, scopes: $renderHookScopes)); ?>


        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! filament()->hasDarkMode()): ?>
            <script>
                localStorage.setItem('theme', 'light')
            </script>
        <?php elseif(filament()->hasDarkModeForced()): ?>
            <script>
                localStorage.setItem('theme', 'dark')
            </script>
        <?php else: ?>
            <script>
                const loadDarkMode = () => {
                    window.theme = localStorage.getItem('theme') ?? <?php echo \Illuminate\Support\Js::from(filament()->getDefaultThemeMode()->value)->toHtml() ?>

                    if (
                        window.theme === 'dark' ||
                        (window.theme === 'system' &&
                            window.matchMedia('(prefers-color-scheme: dark)')
                                .matches)
                    ) {
                        document.documentElement.classList.add('dark')
                    }
                }

                loadDarkMode()

                document.addEventListener('livewire:navigated', loadDarkMode)
            </script>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php echo e(FilamentView::renderHook(PanelsRenderHook::HEAD_END, scopes: $renderHookScopes)); ?>

    </head>

    <body
        <?php echo e($attributes
                ->merge($livewire?->getExtraBodyAttributes() ?? [], escape: false)
                ->class([
                    'fi-body',
                    'fi-panel-' . filament()->getId(),
                ])); ?>

    >
        <?php echo e(FilamentView::renderHook(PanelsRenderHook::BODY_START, scopes: $renderHookScopes)); ?>


        <?php echo e($slot); ?>


        <?php
$__split = function ($name, $params = []) {
    return [$name, $params];
};
[$__name, $__params] = $__split(Notifications::class);

$__keyOuter = $__key ?? null;

$__key = null;
$__componentSlots = [];

$__key ??= \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::generateKey('lw-3369476815-0', $__key);

$__html = app('livewire')->mount($__name, $__params, $__key, $__componentSlots);

echo $__html;

unset($__html);
unset($__key);
$__key = $__keyOuter;
unset($__keyOuter);
unset($__name);
unset($__params);
unset($__componentSlots);
unset($__split);
?>

        <?php echo e(FilamentView::renderHook(PanelsRenderHook::SCRIPTS_BEFORE, scopes: $renderHookScopes)); ?>


        <?php echo \Filament\Support\Facades\FilamentAsset::renderScripts(withCore: true) ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filament()->hasBroadcasting() && config('filament.broadcasting.echo')): ?>
            <script data-navigate-once>
                window.Echo = new window.EchoFactory(<?php echo \Illuminate\Support\Js::from(config('filament.broadcasting.echo'))->toHtml() ?>)

                window.dispatchEvent(new CustomEvent('EchoLoaded'))
            </script>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filament()->hasDarkMode() && (! filament()->hasDarkModeForced())): ?>
            <script>
                loadDarkMode()
            </script>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php echo $__env->yieldPushContent('scripts'); ?>

        <?php echo e(FilamentView::renderHook(PanelsRenderHook::SCRIPTS_AFTER, scopes: $renderHookScopes)); ?>


        <?php echo e(FilamentView::renderHook(PanelsRenderHook::BODY_END, scopes: $renderHookScopes)); ?>

    </body>
</html>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/layout/base.blade.php ENDPATH**/ ?>