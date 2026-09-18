<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'active' => false,
    'activeChildItems' => false,
    'activeIcon' => null,
    'badge' => null,
    'badgeColor' => null,
    'badgeTooltip' => null,
    'childItems' => [],
    'first' => false,
    'grouped' => false,
    'icon' => null,
    'last' => false,
    'shouldOpenUrlInNewTab' => false,
    'sidebarCollapsible' => true,
    'subGrouped' => false,
    'subNavigation' => false,
    'url',
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
    'activeChildItems' => false,
    'activeIcon' => null,
    'badge' => null,
    'badgeColor' => null,
    'badgeTooltip' => null,
    'childItems' => [],
    'first' => false,
    'grouped' => false,
    'icon' => null,
    'last' => false,
    'shouldOpenUrlInNewTab' => false,
    'sidebarCollapsible' => true,
    'subGrouped' => false,
    'subNavigation' => false,
    'url',
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Support\Enums\IconSize;
    use Filament\Support\View\ComponentAttributeBag;

    $sidebarCollapsible = $sidebarCollapsible && filament()->isSidebarCollapsibleOnDesktop();
?>

<li
    <?php echo e($attributes->class([
            'fi-sidebar-item',
            'fi-active' => $active,
            'fi-sidebar-item-has-active-child-items' => $activeChildItems,
            'fi-sidebar-item-has-url' => filled($url),
        ])); ?>

>
    <a
        <?php echo e(\Filament\Support\generate_href_html($url, $shouldOpenUrlInNewTab)); ?>

        <?php if($active): ?>
            aria-current="page"
        <?php endif; ?>
        x-on:click="window.matchMedia(`(max-width: 1024px)`).matches && $store.sidebar.close()"
        <?php if($sidebarCollapsible && (! $subNavigation)): ?>
            x-bind:aria-label="$store.sidebar.isOpen ? null : <?php echo \Illuminate\Support\Js::from(trim(strip_tags($slot->toHtml())))->toHtml() ?>"
            x-data="{ tooltip: false }"
            x-effect="
                tooltip = $store.sidebar.isOpen
                    ? false
                    : {
                          content: <?php echo \Illuminate\Support\Js::from($slot->toHtml())->toHtml() ?>,
                          placement: document.dir === 'rtl' ? 'left' : 'right',
                          theme: $store.theme,
                      }
            "
            x-tooltip.html="tooltip"
        <?php endif; ?>
        class="fi-sidebar-item-btn"
    >
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($icon) && ((! $subGrouped) || ($sidebarCollapsible && (! $subNavigation)))): ?>
            <?php echo e(\Filament\Support\generate_icon_html(($active && $activeIcon) ? $activeIcon : $icon, attributes: (new ComponentAttributeBag([
                    'x-show' => ($subGrouped && $sidebarCollapsible) ? '! $store.sidebar.isOpen' : false,
                ]))->class(['fi-sidebar-item-icon']), size: IconSize::Large)); ?>

        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if((blank($icon) && $grouped) || $subGrouped): ?>
            <div
                <?php if(filled($icon) && $subGrouped && $sidebarCollapsible && (! $subNavigation)): ?>
                    x-show="$store.sidebar.isOpen"
                <?php endif; ?>
                class="fi-sidebar-item-grouped-border"
            >
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $first): ?>
                    <div
                        class="fi-sidebar-item-grouped-border-part-not-first"
                    ></div>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $last): ?>
                    <div
                        class="fi-sidebar-item-grouped-border-part-not-last"
                    ></div>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <div class="fi-sidebar-item-grouped-border-part"></div>
            </div>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <span
            <?php if($sidebarCollapsible && (! $subNavigation)): ?>
                x-show="$store.sidebar.isOpen"
                x-transition:enter="fi-transition-enter"
                x-transition:enter-start="fi-transition-enter-start"
                x-transition:enter-end="fi-transition-enter-end"
            <?php endif; ?>
            class="fi-sidebar-item-label"
        >
            <?php echo e($slot); ?>

        </span>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($badge)): ?>
            <span
                <?php if($sidebarCollapsible && (! $subNavigation)): ?>
                    x-show="$store.sidebar.isOpen"
                    x-transition:enter="fi-transition-enter"
                    x-transition:enter-start="fi-transition-enter-start"
                    x-transition:enter-end="fi-transition-enter-end"
                <?php endif; ?>
                class="fi-sidebar-item-badge-ctn"
            >
                <?php if (isset($component)) { $__componentOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.badge','data' => ['color' => $badgeColor,'tooltip' => $badgeTooltip]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::badge'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($badgeColor),'tooltip' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($badgeTooltip)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                    <?php echo e($badge); ?>

                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f)): ?>
<?php $attributes = $__attributesOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f; ?>
<?php unset($__attributesOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f)): ?>
<?php $component = $__componentOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f; ?>
<?php unset($__componentOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f); ?>
<?php endif; ?>
            </span>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    </a>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($childItems && (blank($url) || $active || $activeChildItems)): ?>
        <ul class="fi-sidebar-sub-group-items">
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $childItems; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $childItem): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                <?php
                    $isChildItemChildItemsActive = $childItem->isChildItemsActive();
                    $isChildActive = (! $isChildItemChildItemsActive) && $childItem->isActive();
                    $childItemActiveIcon = $childItem->getActiveIcon();
                    $childItemBadge = $childItem->getBadge();
                    $childItemBadgeColor = $childItem->getBadgeColor($childItemBadge);
                    $childItemBadgeTooltip = $childItem->getBadgeTooltip($childItemBadge);
                    $childItemIcon = $childItem->getIcon();
                    $shouldChildItemOpenUrlInNewTab = $childItem->shouldOpenUrlInNewTab();
                    $childItemUrl = $childItem->getUrl();
                    $childItemExtraAttributes = $childItem->getExtraAttributeBag();
                ?>

                <?php if (isset($component)) { $__componentOriginalef06dae617244e2c8d6e31778b7ba11215234367c29009e82c70b0dc82775c55 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalef06dae617244e2c8d6e31778b7ba11215234367c29009e82c70b0dc82775c55 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.sidebar.item','data' => ['active' => $isChildActive,'activeChildItems' => $isChildItemChildItemsActive,'activeIcon' => $childItemActiveIcon,'badge' => $childItemBadge,'badgeColor' => $childItemBadgeColor,'badgeTooltip' => $childItemBadgeTooltip,'first' => $loop->first,'grouped' => true,'icon' => $childItemIcon,'last' => $loop->last,'shouldOpenUrlInNewTab' => $shouldChildItemOpenUrlInNewTab,'subGrouped' => true,'subNavigation' => $subNavigation,'url' => $childItemUrl,'attributes' => \Filament\Support\prepare_inherited_attributes($childItemExtraAttributes)]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::sidebar.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['active' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isChildActive),'active-child-items' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isChildItemChildItemsActive),'active-icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($childItemActiveIcon),'badge' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($childItemBadge),'badge-color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($childItemBadgeColor),'badge-tooltip' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($childItemBadgeTooltip),'first' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($loop->first),'grouped' => true,'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($childItemIcon),'last' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($loop->last),'should-open-url-in-new-tab' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($shouldChildItemOpenUrlInNewTab),'sub-grouped' => true,'sub-navigation' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($subNavigation),'url' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($childItemUrl),'attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(\Filament\Support\prepare_inherited_attributes($childItemExtraAttributes))]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                    <?php echo e($childItem->getLabel()); ?>

                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalef06dae617244e2c8d6e31778b7ba11215234367c29009e82c70b0dc82775c55)): ?>
<?php $attributes = $__attributesOriginalef06dae617244e2c8d6e31778b7ba11215234367c29009e82c70b0dc82775c55; ?>
<?php unset($__attributesOriginalef06dae617244e2c8d6e31778b7ba11215234367c29009e82c70b0dc82775c55); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalef06dae617244e2c8d6e31778b7ba11215234367c29009e82c70b0dc82775c55)): ?>
<?php $component = $__componentOriginalef06dae617244e2c8d6e31778b7ba11215234367c29009e82c70b0dc82775c55; ?>
<?php unset($__componentOriginalef06dae617244e2c8d6e31778b7ba11215234367c29009e82c70b0dc82775c55); ?>
<?php endif; ?>
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
        </ul>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
</li>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/sidebar/item.blade.php ENDPATH**/ ?>