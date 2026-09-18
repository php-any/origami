<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'icon',
    'theme',
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
    'icon',
    'theme',
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\View\PanelsIconAlias;

    $label = __("filament-panels::layout.actions.theme_switcher.{$theme}.label");
?>

<button
    aria-label="<?php echo e($label); ?>"
    type="button"
    x-on:click="(theme = <?php echo \Illuminate\Support\Js::from($theme)->toHtml() ?>) && close()"
    x-tooltip="{
        content: <?php echo \Illuminate\Support\Js::from($label)->toHtml() ?>,
        theme: $store.theme,
    }"
    x-bind:aria-pressed="theme === <?php echo \Illuminate\Support\Js::from($theme)->toHtml() ?> ? 'true' : 'false'"
    x-bind:class="{ 'fi-active': theme === <?php echo \Illuminate\Support\Js::from($theme)->toHtml() ?> }"
    class="fi-theme-switcher-btn"
>
    <?php echo e(\Filament\Support\generate_icon_html($icon, alias: match ($theme) {
            'light' => PanelsIconAlias::THEME_SWITCHER_LIGHT_BUTTON,
            'dark' => PanelsIconAlias::THEME_SWITCHER_DARK_BUTTON,
            'system' => PanelsIconAlias::THEME_SWITCHER_SYSTEM_BUTTON,
        })); ?>

</button>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/theme-switcher/button.blade.php ENDPATH**/ ?>