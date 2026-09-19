<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'badge' => null,
    'badgeColor' => 'primary',
    'badgeSize' => null,
    'color' => 'primary',
    'disabled' => false,
    'form' => null,
    'formId' => null,
    'href' => null,
    'icon' => null,
    'iconAlias' => null,
    'iconPosition' => null,
    'iconSize' => null,
    'keyBindings' => null,
    'labelSrOnly' => false,
    'loadingIndicator' => true,
    'size' => null,
    'spaMode' => null,
    'tag' => 'a',
    'target' => null,
    'tooltip' => null,
    'type' => 'button',
    'weight' => null,
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
    'badge' => null,
    'badgeColor' => 'primary',
    'badgeSize' => null,
    'color' => 'primary',
    'disabled' => false,
    'form' => null,
    'formId' => null,
    'href' => null,
    'icon' => null,
    'iconAlias' => null,
    'iconPosition' => null,
    'iconSize' => null,
    'keyBindings' => null,
    'labelSrOnly' => false,
    'loadingIndicator' => true,
    'size' => null,
    'spaMode' => null,
    'tag' => 'a',
    'target' => null,
    'tooltip' => null,
    'type' => 'button',
    'weight' => null,
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Support\Enums\FontWeight;
    use Filament\Support\Enums\IconPosition;
    use Filament\Support\Enums\IconSize;
    use Filament\Support\Enums\Size;
    use Filament\Support\View\ComponentAttributeBag as FilamentComponentAttributeBag;
    use Filament\Support\View\Components\BadgeComponent;
    use Filament\Support\View\Components\LinkComponent;
    use Illuminate\Contracts\Support\Htmlable;
    use Illuminate\View\ComponentAttributeBag;
    use Illuminate\View\ComponentSlot;

    $badgeSize ??= Size::ExtraSmall;
    $iconPosition ??= IconPosition::Before;
    $size ??= Size::Medium;

    if (! $iconPosition instanceof IconPosition) {
        $iconPosition = filled($iconPosition) ? (IconPosition::tryFrom($iconPosition) ?? $iconPosition) : null;
    }

    if (! $badgeSize instanceof Size) {
        $badgeSize = filled($badgeSize) ? (Size::tryFrom($badgeSize) ?? $badgeSize) : null;
    }

    if (! $size instanceof Size) {
        $size = filled($size) ? (Size::tryFrom($size) ?? $size) : null;
    }

    if (filled($iconSize) && (! $iconSize instanceof IconSize)) {
        $iconSize = IconSize::tryFrom($iconSize) ?? $iconSize;
    }

    $iconSize ??= match ($size) {
        Size::ExtraSmall, Size::Small => IconSize::Small,
        default => null,
    };

    if (! $weight instanceof FontWeight) {
        $weight = filled($weight) ? (FontWeight::tryFrom($weight) ?? $weight) : null;
    }

    $wireTarget = $loadingIndicator ? $attributes->whereStartsWith(['wire:target', 'wire:click'])->filter(fn ($value): bool => filled($value))->first() : null;

    $hasLoadingIndicator = filled($wireTarget) || ($type === 'submit' && filled($form));

    if ($hasLoadingIndicator) {
        $loadingIndicatorTarget = html_entity_decode($wireTarget ?: $form, ENT_QUOTES);
    }

    $hasTooltip = filled($tooltip);

    $loadingDelay = ($icon || $hasLoadingIndicator)
        ? config('filament.livewire_loading_delay', 'default')
        : null;
?>

<<?php echo e($tag); ?>

    <?php if(($tag === 'a') && (! ($disabled && $hasTooltip))): ?>
        <?php echo e(\Filament\Support\generate_href_html($href, $target === '_blank', $spaMode)); ?>

    <?php endif; ?>
    <?php if($keyBindings): ?>
        x-bind:id="$id('key-bindings')"
        x-mousetrap.global.<?php echo e(collect($keyBindings)->map(fn (string $keyBinding): string => str_replace('+', '-', $keyBinding))->implode('.')); ?>="document.getElementById($el.id)?.click()"
    <?php endif; ?>
    <?php if($hasTooltip): ?>
        x-tooltip="{
            content: <?php echo \Illuminate\Support\Js::from($tooltip)->toHtml() ?>,
            theme: $store.theme,
            allowHTML: <?php echo \Illuminate\Support\Js::from($tooltip instanceof Htmlable)->toHtml() ?>,
        }"
    <?php endif; ?>
    <?php echo e($attributes
            ->merge([
                'aria-disabled' => $disabled ? 'true' : null,
                // Security: These attributes are rendered without escaping, so the `aria-label` must be escaped here, otherwise an `Htmlable` label could break out of the attribute. `doubleEncode: false` preserves entities that Blade has already escaped in the slot.
                'aria-label' => $labelSrOnly
                ? e(trim(strip_tags($slot->toHtml())), doubleEncode: false)
                : null,
                'disabled' => $disabled && blank($tooltip),
                'form' => $formId,
                'type' => $tag === 'button' ? $type : null,
                'wire:loading.attr' => $tag === 'button' ? 'disabled' : null,
                'wire:target' => ($hasLoadingIndicator && $loadingIndicatorTarget) ? $loadingIndicatorTarget : null,
            ], escape: false)
            ->when(
                $disabled && $hasTooltip,
                fn (ComponentAttributeBag $attributes) => $attributes->filter(
                    fn (mixed $value, string $key): bool => ! str($key)->startsWith(['href', 'x-on:', 'wire:click']),
                ),
            )
            ->class([
                'fi-link',
                'fi-disabled' => $disabled,
                ($size instanceof Size) ? "fi-size-{$size->value}" : (is_string($size) ? $size : ''),
                ($weight instanceof FontWeight) ? "fi-font-{$weight->value}" : (is_string($weight) ? $weight : ''),
            ])
            ->color(LinkComponent::class, $color)); ?>

>
    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($iconPosition === IconPosition::Before): ?>
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($icon || $iconAlias): ?>
            <?php echo e(\Filament\Support\generate_icon_html($icon, $iconAlias, (new Filament\Support\View\ComponentAttributeBag([
                    'wire:loading.remove.delay.' . $loadingDelay => $hasLoadingIndicator,
                    'wire:target' => $hasLoadingIndicator ? $loadingIndicatorTarget : false,
                ])), size: $iconSize)); ?>

        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasLoadingIndicator): ?>
            <?php echo e(\Filament\Support\generate_loading_indicator_html((new Filament\Support\View\ComponentAttributeBag([
                    'wire:loading.delay.' . $loadingDelay => '',
                    'wire:target' => $loadingIndicatorTarget,
                ])), size: $iconSize)); ?>

        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $labelSrOnly): ?>
        <span class="fi-link-label">
            <?php echo e($slot); ?>

        </span>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($iconPosition === IconPosition::After): ?>
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($icon || $iconAlias): ?>
            <?php echo e(\Filament\Support\generate_icon_html($icon, $iconAlias, (new Filament\Support\View\ComponentAttributeBag([
                    'wire:loading.remove.delay.' . $loadingDelay => $hasLoadingIndicator,
                    'wire:target' => $hasLoadingIndicator ? $loadingIndicatorTarget : false,
                ])), size: $iconSize)); ?>

        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasLoadingIndicator): ?>
            <?php echo e(\Filament\Support\generate_loading_indicator_html((new Filament\Support\View\ComponentAttributeBag([
                    'wire:loading.delay.' . $loadingDelay => '',
                    'wire:target' => $loadingIndicatorTarget,
                ])), size: $iconSize)); ?>

        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($badge)): ?>
        <div class="fi-link-badge-ctn">
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($badge instanceof ComponentSlot): ?>
                <?php echo e($badge); ?>

            <?php else: ?>
                <span
                    <?php echo e((new FilamentComponentAttributeBag)->color(BadgeComponent::class, $badgeColor)->class([
                            'fi-badge',
                            ($badgeSize instanceof Size) ? "fi-size-{$badgeSize->value}" : (is_string($badgeSize) ? $badgeSize : ''),
                        ])); ?>

                >
                    <?php echo e($badge); ?>

                </span>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        </div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
</<?php echo e($tag); ?>><?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\support\resources\views/components/link.blade.php ENDPATH**/ ?>