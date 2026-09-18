<div>
    <?php
        use Filament\Enums\GlobalSearchPosition;

        $navigation = filament()->getNavigation();
        $isRtl = __('filament-panels::layout.direction') === 'rtl';
        $isSidebarCollapsibleOnDesktop = filament()->isSidebarCollapsibleOnDesktop();
        $isSidebarFullyCollapsibleOnDesktop = filament()->isSidebarFullyCollapsibleOnDesktop();
        $hasNavigation = filament()->hasNavigation();
        $hasTopbar = filament()->hasTopbar();
        $hasTenantMenu = filament()->hasTenancy() && filament()->hasTenantMenu();
        $hasGlobalSearchInSidebar = filament()->isGlobalSearchEnabled() && filament()->getGlobalSearchPosition() === GlobalSearchPosition::Sidebar;
    ?>

    
    <div
        x-data="{}"
        <?php if($isSidebarCollapsibleOnDesktop || $isSidebarFullyCollapsibleOnDesktop): ?>
            x-cloak
        <?php else: ?>
            x-cloak="-lg"
        <?php endif; ?>
        x-bind:class="{ 'fi-sidebar-open': $store.sidebar.isOpen }"
        id="fi-main-sidebar"
        class="fi-sidebar fi-main-sidebar"
    >
        <?php echo e(\Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::SIDEBAR_START)); ?>


        <div class="fi-sidebar-header-ctn">
            <header
                class="fi-sidebar-header"
            >
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if((! $hasTopbar) && $isSidebarCollapsibleOnDesktop): ?>
                    <?php if (isset($component)) { $__componentOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.icon-button','data' => ['color' => 'gray','icon' => $isRtl ? \Filament\Support\Icons\Heroicon::OutlinedChevronLeft : \Filament\Support\Icons\Heroicon::OutlinedChevronRight,'iconAlias' => 
                            $isRtl
                            ? [
                                \Filament\View\PanelsIconAlias::SIDEBAR_EXPAND_BUTTON_RTL,
                                \Filament\View\PanelsIconAlias::SIDEBAR_EXPAND_BUTTON,
                            ]
                            : \Filament\View\PanelsIconAlias::SIDEBAR_EXPAND_BUTTON
                        ,'iconSize' => 'lg','label' => __('filament-panels::layout.actions.sidebar.expand.label'),'xCloak' => true,'xData' => '{}','ariaControls' => 'fi-main-sidebar','xBind:ariaExpanded' => '$store.sidebar.isOpen','xOn:click' => '$store.sidebar.open()','xShow' => '! $store.sidebar.isOpen','class' => 'fi-sidebar-open-collapse-sidebar-btn']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::icon-button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'gray','icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isRtl ? \Filament\Support\Icons\Heroicon::OutlinedChevronLeft : \Filament\Support\Icons\Heroicon::OutlinedChevronRight),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
                            $isRtl
                            ? [
                                \Filament\View\PanelsIconAlias::SIDEBAR_EXPAND_BUTTON_RTL,
                                \Filament\View\PanelsIconAlias::SIDEBAR_EXPAND_BUTTON,
                            ]
                            : \Filament\View\PanelsIconAlias::SIDEBAR_EXPAND_BUTTON
                        ),'icon-size' => 'lg','label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament-panels::layout.actions.sidebar.expand.label')),'x-cloak' => true,'x-data' => '{}','aria-controls' => 'fi-main-sidebar','x-bind:aria-expanded' => '$store.sidebar.isOpen','x-on:click' => '$store.sidebar.open()','x-show' => '! $store.sidebar.isOpen','class' => 'fi-sidebar-open-collapse-sidebar-btn']); ?>
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

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if((! $hasTopbar) && ($isSidebarCollapsibleOnDesktop || $isSidebarFullyCollapsibleOnDesktop)): ?>
                    <?php if (isset($component)) { $__componentOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.icon-button','data' => ['color' => 'gray','icon' => $isRtl ? \Filament\Support\Icons\Heroicon::OutlinedChevronRight : \Filament\Support\Icons\Heroicon::OutlinedChevronLeft,'iconAlias' => 
                            $isRtl
                            ? [
                                \Filament\View\PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON_RTL,
                                \Filament\View\PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON,
                            ]
                            : \Filament\View\PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON
                        ,'iconSize' => 'lg','label' => __('filament-panels::layout.actions.sidebar.collapse.label'),'xCloak' => true,'xData' => '{}','ariaControls' => 'fi-main-sidebar','xBind:ariaExpanded' => '$store.sidebar.isOpen','xOn:click' => '$store.sidebar.close()','xShow' => '$store.sidebar.isOpen','class' => 'fi-sidebar-close-collapse-sidebar-btn']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::icon-button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'gray','icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isRtl ? \Filament\Support\Icons\Heroicon::OutlinedChevronRight : \Filament\Support\Icons\Heroicon::OutlinedChevronLeft),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
                            $isRtl
                            ? [
                                \Filament\View\PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON_RTL,
                                \Filament\View\PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON,
                            ]
                            : \Filament\View\PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON
                        ),'icon-size' => 'lg','label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament-panels::layout.actions.sidebar.collapse.label')),'x-cloak' => true,'x-data' => '{}','aria-controls' => 'fi-main-sidebar','x-bind:aria-expanded' => '$store.sidebar.isOpen','x-on:click' => '$store.sidebar.close()','x-show' => '$store.sidebar.isOpen','class' => 'fi-sidebar-close-collapse-sidebar-btn']); ?>
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

                <?php echo e(\Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::SIDEBAR_LOGO_BEFORE)); ?>


                <div
                    <?php if($isSidebarCollapsibleOnDesktop || $isSidebarFullyCollapsibleOnDesktop): ?>
                        x-show="$store.sidebar.isOpen"
                    <?php endif; ?>
                    class="fi-sidebar-header-logo-ctn"
                >
                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($homeUrl = filament()->getHomeUrl()): ?>
                        <a <?php echo e(\Filament\Support\generate_href_html($homeUrl)); ?>>
                            <?php if (isset($component)) { $__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.logo','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::logo'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8)): ?>
<?php $attributes = $__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8; ?>
<?php unset($__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8)): ?>
<?php $component = $__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8; ?>
<?php unset($__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8); ?>
<?php endif; ?>
                        </a>
                    <?php else: ?>
                        <?php if (isset($component)) { $__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.logo','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::logo'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8)): ?>
<?php $attributes = $__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8; ?>
<?php unset($__attributesOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8)): ?>
<?php $component = $__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8; ?>
<?php unset($__componentOriginalc3719a76dbf502761dcd34324c2546c67febeb13347d2c6a3e5f12c030657fe8); ?>
<?php endif; ?>
                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                </div>

                <?php echo e(\Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::SIDEBAR_LOGO_AFTER)); ?>

            </header>
        </div>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasTenantMenu || $hasGlobalSearchInSidebar): ?>
            <div
                <?php if((! $hasTenantMenu) && ($isSidebarCollapsibleOnDesktop || $isSidebarFullyCollapsibleOnDesktop)): ?>
                    x-show="$store.sidebar.isOpen"
                <?php endif; ?>
                class="fi-sidebar-header-controls"
            >
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasTenantMenu): ?>
                    <?php if (isset($component)) { $__componentOriginal346b160454f48b4f0880dbc7b2166144211a2adb32656c45a9b55ee5b2d5bf7e = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal346b160454f48b4f0880dbc7b2166144211a2adb32656c45a9b55ee5b2d5bf7e = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.tenant-menu','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::tenant-menu'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal346b160454f48b4f0880dbc7b2166144211a2adb32656c45a9b55ee5b2d5bf7e)): ?>
<?php $attributes = $__attributesOriginal346b160454f48b4f0880dbc7b2166144211a2adb32656c45a9b55ee5b2d5bf7e; ?>
<?php unset($__attributesOriginal346b160454f48b4f0880dbc7b2166144211a2adb32656c45a9b55ee5b2d5bf7e); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal346b160454f48b4f0880dbc7b2166144211a2adb32656c45a9b55ee5b2d5bf7e)): ?>
<?php $component = $__componentOriginal346b160454f48b4f0880dbc7b2166144211a2adb32656c45a9b55ee5b2d5bf7e; ?>
<?php unset($__componentOriginal346b160454f48b4f0880dbc7b2166144211a2adb32656c45a9b55ee5b2d5bf7e); ?>
<?php endif; ?>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasGlobalSearchInSidebar): ?>
                    <div
                        <?php if($isSidebarCollapsibleOnDesktop || $isSidebarFullyCollapsibleOnDesktop): ?>
                            x-show="$store.sidebar.isOpen"
                        <?php endif; ?>
                    >
                        <?php
$__split = function ($name, $params = []) {
    return [$name, $params];
};
[$__name, $__params] = $__split(Filament\Livewire\GlobalSearch::class);

$__keyOuter = $__key ?? null;

$__key = null;
$__componentSlots = [];

$__key ??= \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::generateKey('lw-3667863714-0', $__key);

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
                    </div>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
            </div>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <nav
            aria-label="<?php echo e(__('filament-panels::layout.navigation.label')); ?>"
            class="fi-sidebar-nav"
        >
            <?php echo e(\Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::SIDEBAR_NAV_START)); ?>


            <ul class="fi-sidebar-nav-groups">
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $navigation; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $group): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                    <?php
                        $isGroupActive = $group->isActive();
                        $isGroupCollapsible = $group->isCollapsible();
                        $groupIcon = $group->getIcon();
                        $groupItems = $group->getItems();
                        $groupLabel = $group->getLabel();
                        $groupExtraSidebarAttributeBag = $group->getExtraSidebarAttributeBag();
                    ?>

                    <?php if (isset($component)) { $__componentOriginal86b198484e8dd17672d4ca7c4cacdf366c2e4c078cd0e64389c4da880ac14819 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal86b198484e8dd17672d4ca7c4cacdf366c2e4c078cd0e64389c4da880ac14819 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.sidebar.group','data' => ['active' => $isGroupActive,'collapsible' => $isGroupCollapsible,'icon' => $groupIcon,'items' => $groupItems,'label' => $groupLabel,'attributes' => \Filament\Support\prepare_inherited_attributes($groupExtraSidebarAttributeBag)]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::sidebar.group'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['active' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isGroupActive),'collapsible' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isGroupCollapsible),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupIcon),'items' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupItems),'label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupLabel),'attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(\Filament\Support\prepare_inherited_attributes($groupExtraSidebarAttributeBag))]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal86b198484e8dd17672d4ca7c4cacdf366c2e4c078cd0e64389c4da880ac14819)): ?>
<?php $attributes = $__attributesOriginal86b198484e8dd17672d4ca7c4cacdf366c2e4c078cd0e64389c4da880ac14819; ?>
<?php unset($__attributesOriginal86b198484e8dd17672d4ca7c4cacdf366c2e4c078cd0e64389c4da880ac14819); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal86b198484e8dd17672d4ca7c4cacdf366c2e4c078cd0e64389c4da880ac14819)): ?>
<?php $component = $__componentOriginal86b198484e8dd17672d4ca7c4cacdf366c2e4c078cd0e64389c4da880ac14819; ?>
<?php unset($__componentOriginal86b198484e8dd17672d4ca7c4cacdf366c2e4c078cd0e64389c4da880ac14819); ?>
<?php endif; ?>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
            </ul>

            <script>
                var collapsedGroups = JSON.parse(
                    localStorage.getItem('collapsedGroups'),
                )

                if (collapsedGroups === null || collapsedGroups === 'null') {
                    localStorage.setItem(
                        'collapsedGroups',
                        JSON.stringify(<?php echo \Illuminate\Support\Js::from(
                        collect($navigation)
                            ->filter(fn (\Filament\Navigation\NavigationGroup $group): bool => $group->isCollapsed())
                            ->map(fn (\Filament\Navigation\NavigationGroup $group): string => $group->getLabel())
                            ->values()
                            ->all()
                    )->toHtml() ?>),
                    )
                }

                collapsedGroups = JSON.parse(
                    localStorage.getItem('collapsedGroups'),
                )

                document
                    .querySelectorAll('.fi-sidebar-group')
                    .forEach((group) => {
                        if (
                            !collapsedGroups.includes(group.dataset.groupLabel)
                        ) {
                            return
                        }

                        // Alpine.js loads too slow, so attempt to hide a
                        // collapsed sidebar group earlier.
                        group.querySelector(
                            '.fi-sidebar-group-items',
                        ).style.display = 'none'
                        group.classList.add('fi-collapsed')
                    })
            </script>

            <?php echo e(\Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::SIDEBAR_NAV_END)); ?>

        </nav>

        <?php
            $isAuthenticated = filament()->auth()->check();
            $hasDatabaseNotificationsInSidebar = filament()->hasDatabaseNotifications() && filament()->getDatabaseNotificationsPosition() === \Filament\Enums\DatabaseNotificationsPosition::Sidebar;
            $hasUserMenuInSidebar = filament()->hasUserMenu() && filament()->getUserMenuPosition() === \Filament\Enums\UserMenuPosition::Sidebar;
            $shouldRenderFooter = $isAuthenticated && ($hasDatabaseNotificationsInSidebar || $hasUserMenuInSidebar);
        ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($shouldRenderFooter): ?>
            <div class="fi-sidebar-footer">
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasDatabaseNotificationsInSidebar): ?>
                    <?php
$__split = function ($name, $params = []) {
    return [$name, $params];
};
[$__name, $__params] = $__split(filament()->getDatabaseNotificationsLivewireComponent(), [
                        'lazy' => filament()->hasLazyLoadedDatabaseNotifications(),
                    ]);

$__keyOuter = $__key ?? null;

$__key = null;
$__componentSlots = [];

$__key ??= \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::generateKey('lw-3667863714-0', $__key);

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
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasUserMenuInSidebar): ?>
                    <?php if (isset($component)) { $__componentOriginal35440ad68537cb637b8820876dc0165490d1dcdd28d4a73c8a0a5ee8f0a69896 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal35440ad68537cb637b8820876dc0165490d1dcdd28d4a73c8a0a5ee8f0a69896 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.user-menu','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::user-menu'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal35440ad68537cb637b8820876dc0165490d1dcdd28d4a73c8a0a5ee8f0a69896)): ?>
<?php $attributes = $__attributesOriginal35440ad68537cb637b8820876dc0165490d1dcdd28d4a73c8a0a5ee8f0a69896; ?>
<?php unset($__attributesOriginal35440ad68537cb637b8820876dc0165490d1dcdd28d4a73c8a0a5ee8f0a69896); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal35440ad68537cb637b8820876dc0165490d1dcdd28d4a73c8a0a5ee8f0a69896)): ?>
<?php $component = $__componentOriginal35440ad68537cb637b8820876dc0165490d1dcdd28d4a73c8a0a5ee8f0a69896; ?>
<?php unset($__componentOriginal35440ad68537cb637b8820876dc0165490d1dcdd28d4a73c8a0a5ee8f0a69896); ?>
<?php endif; ?>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
            </div>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php echo e(\Filament\Support\Facades\FilamentView::renderHook(\Filament\View\PanelsRenderHook::SIDEBAR_FOOTER)); ?>

    </div>
    

    <?php if (isset($component)) { $__componentOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-actions::components.modals','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-actions::modals'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15)): ?>
<?php $attributes = $__attributesOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15; ?>
<?php unset($__attributesOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15)): ?>
<?php $component = $__componentOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15; ?>
<?php unset($__componentOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15); ?>
<?php endif; ?>
</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/livewire/sidebar.blade.php ENDPATH**/ ?>