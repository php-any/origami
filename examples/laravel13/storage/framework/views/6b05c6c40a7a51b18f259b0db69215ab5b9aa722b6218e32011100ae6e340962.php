<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'alpineDisabled' => null,
    'alpineValid' => null,
    'disabled' => false,
    'inlinePrefix' => false,
    'inlineSuffix' => false,
    'prefix' => null,
    'prefixActions' => [],
    'prefixIcon' => null,
    'prefixIconColor' => 'gray',
    'prefixIconAlias' => null,
    'suffix' => null,
    'suffixActions' => [],
    'suffixIcon' => null,
    'suffixIconColor' => 'gray',
    'suffixIconAlias' => null,
    'valid' => true,
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
    'alpineDisabled' => null,
    'alpineValid' => null,
    'disabled' => false,
    'inlinePrefix' => false,
    'inlineSuffix' => false,
    'prefix' => null,
    'prefixActions' => [],
    'prefixIcon' => null,
    'prefixIconColor' => 'gray',
    'prefixIconAlias' => null,
    'suffix' => null,
    'suffixActions' => [],
    'suffixIcon' => null,
    'suffixIconColor' => 'gray',
    'suffixIconAlias' => null,
    'valid' => true,
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Actions\Action;
    use Filament\Support\View\ComponentAttributeBag;
    use Filament\Support\View\Components\InputComponent\WrapperComponent\IconComponent;
    use Illuminate\Support\Str;

    $prefixActions = array_filter(
        $prefixActions,
        fn (Action $prefixAction): bool => $prefixAction->isVisible(),
    );

    $suffixActions = array_filter(
        $suffixActions,
        fn (Action $suffixAction): bool => $suffixAction->isVisible(),
    );

    $hasAlpineDisabledClasses = filled($alpineDisabled);
    $hasAlpineValidClasses = filled($alpineValid);
    $hasAlpineClasses = $hasAlpineDisabledClasses || $hasAlpineValidClasses;

    $wireTarget = $attributes->whereStartsWith(['wire:target'])->first();

    $hasLoadingIndicator = filled($wireTarget);

    if ($hasLoadingIndicator) {
        $loadingIndicatorTarget = html_entity_decode($wireTarget, ENT_QUOTES);
    }

    $loadingDelay = ($prefixIcon || $prefixIconAlias || $suffixIcon || $suffixIconAlias || $hasLoadingIndicator)
        ? config('filament.livewire_loading_delay', 'default')
        : null;

    $prefixIconHtml = ($prefixIcon || $prefixIconAlias)
        ? \Filament\Support\generate_icon_html($prefixIcon, $prefixIconAlias, (new ComponentAttributeBag)
            ->merge([
                'wire:loading.remove.delay.' . $loadingDelay => $hasLoadingIndicator,
                'wire:target' => $hasLoadingIndicator ? $loadingIndicatorTarget : false,
            ], escape: false)
            ->color(IconComponent::class, $prefixIconColor))
        : null;

    $suffixIconHtml = ($suffixIcon || $suffixIconAlias)
        ? \Filament\Support\generate_icon_html($suffixIcon, $suffixIconAlias, (new ComponentAttributeBag)
            ->merge([
                'wire:loading.remove.delay.' . $loadingDelay => $hasLoadingIndicator,
                'wire:target' => $hasLoadingIndicator ? $loadingIndicatorTarget : false,
            ], escape: false)
            ->color(IconComponent::class, $suffixIconColor))
        : null;

    $hasPrefix = count($prefixActions) || ($prefixIconHtml !== null) || filled($prefix);
    $hasSuffix = count($suffixActions) || ($suffixIconHtml !== null) || filled($suffix);

    $hasFocusInputListener = $attributes->has('x-on:focus-input.stop');
    $canClickPrefixAffix = $hasFocusInputListener && (($prefixIconHtml !== null) || filled($prefix));
    $canClickSuffixAffix = $hasFocusInputListener && (($suffixIconHtml !== null) || filled($suffix));
?>

<div
    <?php if($hasAlpineClasses): ?>
        x-bind:class="{
            <?php echo e($hasAlpineDisabledClasses ? "'fi-disabled': {$alpineDisabled}," : null); ?>

            <?php echo e($hasAlpineValidClasses ? "'fi-invalid': ! ({$alpineValid})," : null); ?>

        }"
    <?php endif; ?>
    <?php echo e($attributes
            ->except(['wire:target', 'tabindex'])
            ->class([
                'fi-input-wrp',
                'fi-disabled' => (! $hasAlpineClasses) && $disabled,
                'fi-invalid' => (! $hasAlpineClasses) && (! $valid),
            ])); ?>

>
    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasPrefix || $hasLoadingIndicator): ?>
        <div
            <?php if(! $hasPrefix): ?>
                wire:loading.delay.<?php echo e($loadingDelay); ?>.flex
                wire:target="<?php echo e($loadingIndicatorTarget); ?>"
                <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e(Str::random()).''; ?>wire:key="<?php echo e(Str::random()); ?>" 
            <?php endif; ?>
            <?php if($canClickPrefixAffix): ?>
                x-on:click="$dispatch('focus-input')"
            <?php endif; ?>
            class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                'fi-input-wrp-prefix',
                'fi-input-wrp-prefix-has-content' => $hasPrefix,
                'fi-inline' => $inlinePrefix,
                'fi-input-wrp-prefix-has-label' => filled($prefix),
            ]); ?>"
        >
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(count($prefixActions)): ?>
                <div
                    class="<?php echo \Illuminate\Support\Arr::toCssClasses(['fi-input-wrp-actions']); ?>"
                    <?php if($canClickPrefixAffix): ?> x-on:click.stop <?php endif; ?>
                >
                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $prefixActions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $prefixAction): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                        <?php echo e($prefixAction); ?>

                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                </div>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php echo e($prefixIconHtml); ?>


            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasLoadingIndicator): ?>
                <?php echo e(\Filament\Support\generate_loading_indicator_html((new ComponentAttributeBag([
                        'wire:loading.delay.' . $loadingDelay => $hasPrefix,
                        'wire:target' => $hasPrefix ? $loadingIndicatorTarget : null,
                    ]))->color(IconComponent::class, 'gray'))); ?>

            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($prefix)): ?>
                <span class="fi-input-wrp-label">
                    <?php echo e($prefix); ?>

                </span>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        </div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <div
        <?php if($hasLoadingIndicator && (! $hasPrefix)): ?>
            <?php if($inlinePrefix): ?>
                wire:loading.delay.<?php echo e($loadingDelay); ?>.class.remove="ps-3"
            <?php endif; ?>

            wire:target="<?php echo e($loadingIndicatorTarget); ?>"
        <?php endif; ?>
        class="<?php echo \Illuminate\Support\Arr::toCssClasses([
            'fi-input-wrp-content-ctn',
            'fi-input-wrp-content-ctn-ps' => $hasLoadingIndicator && (! $hasPrefix) && $inlinePrefix,
        ]); ?>"
    >
        <?php echo e($slot); ?>

    </div>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasSuffix): ?>
        <div
            <?php if($canClickSuffixAffix): ?>
                x-on:click="$dispatch('focus-input')"
            <?php endif; ?>
            class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                'fi-input-wrp-suffix',
                'fi-inline' => $inlineSuffix,
                'fi-input-wrp-suffix-has-label' => filled($suffix),
            ]); ?>"
        >
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($suffix)): ?>
                <span class="fi-input-wrp-label">
                    <?php echo e($suffix); ?>

                </span>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php echo e($suffixIconHtml); ?>


            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(count($suffixActions)): ?>
                <div
                    class="<?php echo \Illuminate\Support\Arr::toCssClasses(['fi-input-wrp-actions']); ?>"
                    <?php if($canClickSuffixAffix): ?> x-on:click.stop <?php endif; ?>
                >
                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $suffixActions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $suffixAction): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                        <?php echo e($suffixAction); ?>

                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                </div>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        </div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\support\resources\views/components/input/wrapper.blade.php ENDPATH**/ ?>