<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'active' => false,
    'collapsible' => true,
    'icon' => null,
    'items' => [],
    'label' => null,
    'sidebarCollapsible' => true,
    'subNavigation' => false,
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
    'collapsible' => true,
    'icon' => null,
    'items' => [],
    'label' => null,
    'sidebarCollapsible' => true,
    'subNavigation' => false,
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
    use Filament\Support\Icons\Heroicon;
    use Filament\View\PanelsIconAlias;
    use Illuminate\Contracts\Support\Htmlable;
    use Illuminate\Support\Str;

    $sidebarCollapsible = $sidebarCollapsible && filament()->isSidebarCollapsibleOnDesktop();
    $hasDropdown = filled($label) && filled($icon) && $sidebarCollapsible;
    // A slug alone is not unique: non-Latin labels slug to an empty string and distinct labels can
    // share a slug, producing duplicate ids that break each disclosure button's `aria-controls`.
    // A short hash of the raw label keeps the id unique per label and stable across renders.
    $groupLabel = $subNavigation ? "sub_navigation_{$label}" : (string) $label;
    $groupItemsId = 'fi-sidebar-group-items-' . Str::slug($groupLabel) . '-' . substr(md5($groupLabel), 0, 8);
?>

<li
    x-data="{ label: <?php echo \Illuminate\Support\Js::from($subNavigation ? "sub_navigation_{$label}" : $label)->toHtml() ?> }"
    data-group-label="<?php echo e($subNavigation ? "sub_navigation_{$label}" : $label); ?>"
    x-bind:class="{ 'fi-collapsed': $store.sidebar.groupIsCollapsed(label) }"
    <?php echo e($attributes->class([
            'fi-sidebar-group',
            'fi-active' => $active,
            'fi-collapsible' => $collapsible,
        ])); ?>

>
    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($label): ?>
        <div
            <?php if($collapsible): ?>
                x-on:click="$store.sidebar.toggleCollapsedGroup(label)"
            <?php endif; ?>
            <?php if($sidebarCollapsible): ?>
                x-show="$store.sidebar.isOpen"
                x-transition:enter="fi-transition-enter"
                x-transition:enter-start="fi-transition-enter-start"
                x-transition:enter-end="fi-transition-enter-end"
            <?php endif; ?>
            class="fi-sidebar-group-btn"
        >
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($icon): ?>
                <?php echo e(\Filament\Support\generate_icon_html($icon, size: IconSize::Large)); ?>

            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <span class="fi-sidebar-group-label">
                <?php echo e($label); ?>

            </span>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($collapsible): ?>
                <?php if (isset($component)) { $__componentOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.icon-button','data' => ['color' => 'gray','icon' => Heroicon::ChevronUp,'iconAlias' => PanelsIconAlias::SIDEBAR_GROUP_COLLAPSE_BUTTON,'label' => $label,'ariaControls' => $groupItemsId,'xBind:ariaExpanded' => '! $store.sidebar.groupIsCollapsed(label)','xOn:click.stop' => '$store.sidebar.toggleCollapsedGroup(label)','class' => 'fi-sidebar-group-collapse-btn']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::icon-button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'gray','icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(Heroicon::ChevronUp),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(PanelsIconAlias::SIDEBAR_GROUP_COLLAPSE_BUTTON),'label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($label),'aria-controls' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupItemsId),'x-bind:aria-expanded' => '! $store.sidebar.groupIsCollapsed(label)','x-on:click.stop' => '$store.sidebar.toggleCollapsedGroup(label)','class' => 'fi-sidebar-group-collapse-btn']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad)): ?>
<?php $attributes = $__attributesOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad; ?>
<?php unset($__attributesOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad)): ?>
<?php $component = $__componentOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad; ?>
<?php unset($__componentOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad); ?>
<?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        </div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasDropdown): ?>
        <?php if (isset($component)) { $__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.dropdown.index','data' => ['placement' => (__('filament-panels::layout.direction') === 'rtl') ? 'left-start' : 'right-start','xShow' => '! $store.sidebar.isOpen']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::dropdown'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['placement' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute((__('filament-panels::layout.direction') === 'rtl') ? 'left-start' : 'right-start'),'x-show' => '! $store.sidebar.isOpen']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

             <?php $__env->slot('trigger', null, []); ?> 
                <button
                    aria-label="<?php echo e($label); ?>"
                    x-data="{ tooltip: false }"
                    x-effect="
                        tooltip = $store.sidebar.isOpen
                            ? false
                            : {
                                  content: <?php echo \Illuminate\Support\Js::from($label)->toHtml() ?>,
                                  placement: document.dir === 'rtl' ? 'left' : 'right',
                                  theme: $store.theme,
                              }
                    "
                    x-tooltip.html="tooltip"
                    class="fi-sidebar-group-dropdown-trigger-btn"
                >
                    <?php echo e(\Filament\Support\generate_icon_html($icon, size: IconSize::Large)); ?>

                </button>
             <?php $__env->endSlot(); ?>

            <?php
                $lists = [];

                foreach ($items as $item) {
                    if ($childItems = $item->getChildItems()) {
                        $lists[] = [
                            $item,
                            ...$childItems,
                        ];
                        $lists[] = [];

                        continue;
                    }

                    if (empty($lists)) {
                        $lists[] = [$item];

                        continue;
                    }

                    $lists[count($lists) - 1][] = $item;
                }

                if (empty($lists[count($lists) - 1])) {
                    array_pop($lists);
                }
            ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($label)): ?>
                <?php if (isset($component)) { $__componentOriginal2a4b984efad10636c8eb2e2d0666cfa6ddcf02cf2cbee91a407c355bce805d2b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal2a4b984efad10636c8eb2e2d0666cfa6ddcf02cf2cbee91a407c355bce805d2b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.dropdown.header','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::dropdown.header'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                    <?php echo e($label); ?>

                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal2a4b984efad10636c8eb2e2d0666cfa6ddcf02cf2cbee91a407c355bce805d2b)): ?>
<?php $attributes = $__attributesOriginal2a4b984efad10636c8eb2e2d0666cfa6ddcf02cf2cbee91a407c355bce805d2b; ?>
<?php unset($__attributesOriginal2a4b984efad10636c8eb2e2d0666cfa6ddcf02cf2cbee91a407c355bce805d2b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal2a4b984efad10636c8eb2e2d0666cfa6ddcf02cf2cbee91a407c355bce805d2b)): ?>
<?php $component = $__componentOriginal2a4b984efad10636c8eb2e2d0666cfa6ddcf02cf2cbee91a407c355bce805d2b; ?>
<?php unset($__componentOriginal2a4b984efad10636c8eb2e2d0666cfa6ddcf02cf2cbee91a407c355bce805d2b); ?>
<?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $lists; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $list): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                <?php if (isset($component)) { $__componentOriginal926af3c227e7d331239bd4293e587c334e9e4aaff4fabf2c30bbd11744b0354e = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal926af3c227e7d331239bd4293e587c334e9e4aaff4fabf2c30bbd11744b0354e = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.dropdown.list.index','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::dropdown.list'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $list; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $item): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                        <?php
                            $itemIsActive = $item->isActive();
                            $itemBadge = $item->getBadge();
                            $itemBadgeColor = $item->getBadgeColor($itemBadge);
                            $itemBadgeTooltip = $item->getBadgeTooltip($itemBadge);
                            $itemUrl = $item->getUrl();
                            $itemIcon = $itemIsActive ? ($item->getActiveIcon() ?? $item->getIcon()) : $item->getIcon();
                            $shouldItemOpenUrlInNewTab = $item->shouldOpenUrlInNewTab();
                            $itemExtraAttributes = $item->getExtraAttributeBag();
                        ?>

                        <?php if (isset($component)) { $__componentOriginalc67084d1e4215442fec660f0d8cd4ec60e9c2f12bc4932156a12c753ba32c47e = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalc67084d1e4215442fec660f0d8cd4ec60e9c2f12bc4932156a12c753ba32c47e = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.dropdown.list.item','data' => ['badge' => $itemBadge,'badgeColor' => $itemBadgeColor,'badgeTooltip' => $itemBadgeTooltip,'color' => $itemIsActive ? 'primary' : 'gray','href' => $itemUrl,'icon' => $itemIcon,'tag' => 'a','target' => $shouldItemOpenUrlInNewTab ? '_blank' : null,'ariaCurrent' => $itemIsActive ? 'page' : null,'attributes' => \Filament\Support\prepare_inherited_attributes($itemExtraAttributes)]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::dropdown.list.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['badge' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadge),'badge-color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadgeColor),'badge-tooltip' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadgeTooltip),'color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemIsActive ? 'primary' : 'gray'),'href' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemUrl),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemIcon),'tag' => 'a','target' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($shouldItemOpenUrlInNewTab ? '_blank' : null),'aria-current' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemIsActive ? 'page' : null),'attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(\Filament\Support\prepare_inherited_attributes($itemExtraAttributes))]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                            <?php echo e($item->getLabel()); ?>

                         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalc67084d1e4215442fec660f0d8cd4ec60e9c2f12bc4932156a12c753ba32c47e)): ?>
<?php $attributes = $__attributesOriginalc67084d1e4215442fec660f0d8cd4ec60e9c2f12bc4932156a12c753ba32c47e; ?>
<?php unset($__attributesOriginalc67084d1e4215442fec660f0d8cd4ec60e9c2f12bc4932156a12c753ba32c47e); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalc67084d1e4215442fec660f0d8cd4ec60e9c2f12bc4932156a12c753ba32c47e)): ?>
<?php $component = $__componentOriginalc67084d1e4215442fec660f0d8cd4ec60e9c2f12bc4932156a12c753ba32c47e; ?>
<?php unset($__componentOriginalc67084d1e4215442fec660f0d8cd4ec60e9c2f12bc4932156a12c753ba32c47e); ?>
<?php endif; ?>
                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal926af3c227e7d331239bd4293e587c334e9e4aaff4fabf2c30bbd11744b0354e)): ?>
<?php $attributes = $__attributesOriginal926af3c227e7d331239bd4293e587c334e9e4aaff4fabf2c30bbd11744b0354e; ?>
<?php unset($__attributesOriginal926af3c227e7d331239bd4293e587c334e9e4aaff4fabf2c30bbd11744b0354e); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal926af3c227e7d331239bd4293e587c334e9e4aaff4fabf2c30bbd11744b0354e)): ?>
<?php $component = $__componentOriginal926af3c227e7d331239bd4293e587c334e9e4aaff4fabf2c30bbd11744b0354e; ?>
<?php unset($__componentOriginal926af3c227e7d331239bd4293e587c334e9e4aaff4fabf2c30bbd11744b0354e); ?>
<?php endif; ?>
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b)): ?>
<?php $attributes = $__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b; ?>
<?php unset($__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b)): ?>
<?php $component = $__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b; ?>
<?php unset($__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b); ?>
<?php endif; ?>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <ul
        <?php if(filled($label)): ?>
            id="<?php echo e($groupItemsId); ?>"

            <?php if($sidebarCollapsible): ?>
                x-show="$store.sidebar.isOpen ? ! $store.sidebar.groupIsCollapsed(label) : ! <?php echo \Illuminate\Support\Js::from($hasDropdown)->toHtml() ?>"
            <?php else: ?>
                x-show="! $store.sidebar.groupIsCollapsed(label)"
            <?php endif; ?>
            x-collapse.duration.200ms
        <?php endif; ?>
        <?php if($sidebarCollapsible): ?>
            x-transition:enter="fi-transition-enter"
            x-transition:enter-start="fi-transition-enter-start"
            x-transition:enter-end="fi-transition-enter-end"
        <?php endif; ?>
        class="fi-sidebar-group-items"
    >
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $items; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $item): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
            <?php
                $isItemChildItemsActive = $item->isChildItemsActive();
                $isItemActive = (! $isItemChildItemsActive) && $item->isActive();
                $itemActiveIcon = $item->getActiveIcon();
                $itemBadge = $item->getBadge();
                $itemBadgeColor = $item->getBadgeColor($itemBadge);
                $itemBadgeTooltip = $item->getBadgeTooltip($itemBadge);
                $itemChildItems = $item->getChildItems();
                $itemIcon = $item->getIcon();
                $shouldItemOpenUrlInNewTab = $item->shouldOpenUrlInNewTab();
                $itemUrl = $item->getUrl();
                $itemExtraAttributes = $item->getExtraAttributeBag();

                if ($icon) {
                    if ($hasDropdown || (blank($itemIcon) && blank($itemActiveIcon))) {
                        $itemIcon = null;
                        $itemActiveIcon = null;
                    } else {
                        throw new Exception('Navigation group [' . $label . '] has an icon but one or more of its items also have icons. Either the group or its items can have icons, but not both. This is to ensure a proper user experience.');
                    }
                }
            ?>

            <?php if (isset($component)) { $__componentOriginalef06dae617244e2c8d6e31778b7ba11215234367c29009e82c70b0dc82775c55 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalef06dae617244e2c8d6e31778b7ba11215234367c29009e82c70b0dc82775c55 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.sidebar.item','data' => ['active' => $isItemActive,'activeChildItems' => $isItemChildItemsActive,'activeIcon' => $itemActiveIcon,'badge' => $itemBadge,'badgeColor' => $itemBadgeColor,'badgeTooltip' => $itemBadgeTooltip,'childItems' => $itemChildItems,'first' => $loop->first,'grouped' => filled($label),'icon' => $itemIcon,'last' => $loop->last,'shouldOpenUrlInNewTab' => $shouldItemOpenUrlInNewTab,'sidebarCollapsible' => $sidebarCollapsible,'subNavigation' => $subNavigation,'url' => $itemUrl,'attributes' => \Filament\Support\prepare_inherited_attributes($itemExtraAttributes)]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::sidebar.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['active' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isItemActive),'active-child-items' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isItemChildItemsActive),'active-icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemActiveIcon),'badge' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadge),'badge-color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadgeColor),'badge-tooltip' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadgeTooltip),'child-items' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemChildItems),'first' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($loop->first),'grouped' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(filled($label)),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemIcon),'last' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($loop->last),'should-open-url-in-new-tab' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($shouldItemOpenUrlInNewTab),'sidebar-collapsible' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($sidebarCollapsible),'sub-navigation' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($subNavigation),'url' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemUrl),'attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(\Filament\Support\prepare_inherited_attributes($itemExtraAttributes))]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                <?php echo e($item->getLabel()); ?>


                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($itemIcon instanceof Htmlable): ?>
                     <?php $__env->slot('icon', null, []); ?> 
                        <?php echo e($itemIcon); ?>

                     <?php $__env->endSlot(); ?>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($itemActiveIcon instanceof Htmlable): ?>
                     <?php $__env->slot('activeIcon', null, []); ?> 
                        <?php echo e($itemActiveIcon); ?>

                     <?php $__env->endSlot(); ?>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
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
</li>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/sidebar/group.blade.php ENDPATH**/ ?>