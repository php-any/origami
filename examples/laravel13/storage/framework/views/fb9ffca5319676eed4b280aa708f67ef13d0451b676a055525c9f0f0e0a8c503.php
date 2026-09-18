<?php
    use Filament\Enums\DatabaseNotificationsPosition;
    use Filament\Enums\GlobalSearchPosition;
    use Filament\Enums\UserMenuPosition;
    use Filament\Livewire\GlobalSearch;
    use Filament\Support\Facades\FilamentView;
    use Filament\Support\Icons\Heroicon;
    use Filament\View\PanelsIconAlias;
    use Filament\View\PanelsRenderHook;
?>

<div class="fi-topbar-ctn">
    <?php
        $isRtl = __('filament-panels::layout.direction') === 'rtl';
        $isSidebarCollapsibleOnDesktop = filament()->isSidebarCollapsibleOnDesktop();
        $isSidebarFullyCollapsibleOnDesktop = filament()->isSidebarFullyCollapsibleOnDesktop();
        $hasTopNavigation = filament()->hasTopNavigation();
        $hasNavigation = filament()->hasNavigation();
        $hasTenancy = filament()->hasTenancy();
    ?>

    <nav
        aria-label="<?php echo e(__('filament-panels::layout.topbar.label')); ?>"
        class="fi-topbar"
    >
        <?php echo e(FilamentView::renderHook(PanelsRenderHook::TOPBAR_START)); ?>


        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasNavigation): ?>
            <?php if (isset($component)) { $__componentOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.icon-button','data' => ['color' => 'gray','icon' => Heroicon::OutlinedBars3,'iconAlias' => PanelsIconAlias::TOPBAR_OPEN_SIDEBAR_BUTTON,'iconSize' => 'lg','label' => __('filament-panels::layout.actions.sidebar.expand.label'),'xCloak' => true,'xData' => '{}','ariaControls' => 'fi-main-sidebar','xBind:ariaExpanded' => '$store.sidebar.isOpen','xOn:click' => '$store.sidebar.open()','xShow' => '! $store.sidebar.isOpen','class' => 'fi-topbar-open-sidebar-btn']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::icon-button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'gray','icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(Heroicon::OutlinedBars3),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(PanelsIconAlias::TOPBAR_OPEN_SIDEBAR_BUTTON),'icon-size' => 'lg','label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament-panels::layout.actions.sidebar.expand.label')),'x-cloak' => true,'x-data' => '{}','aria-controls' => 'fi-main-sidebar','x-bind:aria-expanded' => '$store.sidebar.isOpen','x-on:click' => '$store.sidebar.open()','x-show' => '! $store.sidebar.isOpen','class' => 'fi-topbar-open-sidebar-btn']); ?>
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

            <?php if (isset($component)) { $__componentOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.icon-button','data' => ['color' => 'gray','icon' => Heroicon::OutlinedXMark,'iconAlias' => PanelsIconAlias::TOPBAR_CLOSE_SIDEBAR_BUTTON,'iconSize' => 'lg','label' => __('filament-panels::layout.actions.sidebar.collapse.label'),'xCloak' => true,'xData' => '{}','ariaControls' => 'fi-main-sidebar','xBind:ariaExpanded' => '$store.sidebar.isOpen','xOn:click' => '$store.sidebar.close()','xShow' => '$store.sidebar.isOpen','class' => 'fi-topbar-close-sidebar-btn']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::icon-button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'gray','icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(Heroicon::OutlinedXMark),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(PanelsIconAlias::TOPBAR_CLOSE_SIDEBAR_BUTTON),'icon-size' => 'lg','label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament-panels::layout.actions.sidebar.collapse.label')),'x-cloak' => true,'x-data' => '{}','aria-controls' => 'fi-main-sidebar','x-bind:aria-expanded' => '$store.sidebar.isOpen','x-on:click' => '$store.sidebar.close()','x-show' => '$store.sidebar.isOpen','class' => 'fi-topbar-close-sidebar-btn']); ?>
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

        <div class="fi-topbar-start">
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSidebarCollapsibleOnDesktop || $isSidebarFullyCollapsibleOnDesktop): ?>
                <div
                    x-show="$store.sidebar.isOpen || <?php echo \Illuminate\Support\Js::from($isSidebarCollapsibleOnDesktop)->toHtml() ?>"
                    class="fi-topbar-collapse-sidebar-btn-ctn"
                >
                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSidebarCollapsibleOnDesktop): ?>
                        <?php if (isset($component)) { $__componentOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.icon-button','data' => ['color' => 'gray','icon' => $isRtl ? Heroicon::OutlinedChevronLeft : Heroicon::OutlinedChevronRight,'iconAlias' => 
                                $isRtl
                                ? [
                                    PanelsIconAlias::SIDEBAR_EXPAND_BUTTON_RTL,
                                    PanelsIconAlias::SIDEBAR_EXPAND_BUTTON,
                                ]
                                : PanelsIconAlias::SIDEBAR_EXPAND_BUTTON
                            ,'iconSize' => 'lg','label' => __('filament-panels::layout.actions.sidebar.expand.label'),'xCloak' => true,'xData' => '{}','ariaControls' => 'fi-main-sidebar','xBind:ariaExpanded' => '$store.sidebar.isOpen','xOn:click' => '$store.sidebar.open()','xShow' => '! $store.sidebar.isOpen','class' => 'fi-topbar-open-collapse-sidebar-btn']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::icon-button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'gray','icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isRtl ? Heroicon::OutlinedChevronLeft : Heroicon::OutlinedChevronRight),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
                                $isRtl
                                ? [
                                    PanelsIconAlias::SIDEBAR_EXPAND_BUTTON_RTL,
                                    PanelsIconAlias::SIDEBAR_EXPAND_BUTTON,
                                ]
                                : PanelsIconAlias::SIDEBAR_EXPAND_BUTTON
                            ),'icon-size' => 'lg','label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament-panels::layout.actions.sidebar.expand.label')),'x-cloak' => true,'x-data' => '{}','aria-controls' => 'fi-main-sidebar','x-bind:aria-expanded' => '$store.sidebar.isOpen','x-on:click' => '$store.sidebar.open()','x-show' => '! $store.sidebar.isOpen','class' => 'fi-topbar-open-collapse-sidebar-btn']); ?>
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

                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSidebarCollapsibleOnDesktop || $isSidebarFullyCollapsibleOnDesktop): ?>
                        <?php if (isset($component)) { $__componentOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.icon-button','data' => ['color' => 'gray','icon' => $isRtl ? Heroicon::OutlinedChevronRight : Heroicon::OutlinedChevronLeft,'iconAlias' => 
                                $isRtl
                                ? [
                                    PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON_RTL,
                                    PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON,
                                ]
                                : PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON
                            ,'iconSize' => 'lg','label' => __('filament-panels::layout.actions.sidebar.collapse.label'),'xCloak' => true,'xData' => '{}','ariaControls' => 'fi-main-sidebar','xBind:ariaExpanded' => '$store.sidebar.isOpen','xOn:click' => '$store.sidebar.close()','xShow' => '$store.sidebar.isOpen','class' => 'fi-topbar-close-collapse-sidebar-btn']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::icon-button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'gray','icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isRtl ? Heroicon::OutlinedChevronRight : Heroicon::OutlinedChevronLeft),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
                                $isRtl
                                ? [
                                    PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON_RTL,
                                    PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON,
                                ]
                                : PanelsIconAlias::SIDEBAR_COLLAPSE_BUTTON
                            ),'icon-size' => 'lg','label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament-panels::layout.actions.sidebar.collapse.label')),'x-cloak' => true,'x-data' => '{}','aria-controls' => 'fi-main-sidebar','x-bind:aria-expanded' => '$store.sidebar.isOpen','x-on:click' => '$store.sidebar.close()','x-show' => '$store.sidebar.isOpen','class' => 'fi-topbar-close-collapse-sidebar-btn']); ?>
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

            <?php echo e(FilamentView::renderHook(PanelsRenderHook::TOPBAR_LOGO_BEFORE)); ?>


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

            <?php echo e(FilamentView::renderHook(PanelsRenderHook::TOPBAR_LOGO_AFTER)); ?>

        </div>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasTopNavigation || (! $hasNavigation)): ?>
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasTenancy && filament()->hasTenantMenu()): ?>
                <?php if (isset($component)) { $__componentOriginal346b160454f48b4f0880dbc7b2166144211a2adb32656c45a9b55ee5b2d5bf7e = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal346b160454f48b4f0880dbc7b2166144211a2adb32656c45a9b55ee5b2d5bf7e = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.tenant-menu','data' => ['teleport' => true]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::tenant-menu'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['teleport' => true]); ?>
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

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasNavigation): ?>
                <?php
                    $navigation = filament()->getNavigation();
                ?>

                <ul class="fi-topbar-nav-groups">
                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $navigation; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $group): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                        <?php
                            $groupLabel = $group->getLabel();
                            $groupExtraTopbarAttributeBag = $group->getExtraTopbarAttributeBag();
                            $isGroupActive = $group->isActive();
                            $groupIcon = $group->getIcon();
                        ?>

                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($groupLabel): ?>
                            <?php if (isset($component)) { $__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.dropdown.index','data' => ['placement' => 'bottom-start','teleport' => true,'attributes' => \Filament\Support\prepare_inherited_attributes($groupExtraTopbarAttributeBag)]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::dropdown'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['placement' => 'bottom-start','teleport' => true,'attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(\Filament\Support\prepare_inherited_attributes($groupExtraTopbarAttributeBag))]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                 <?php $__env->slot('trigger', null, []); ?> 
                                    <?php if (isset($component)) { $__componentOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.topbar.item','data' => ['active' => $isGroupActive,'icon' => $groupIcon]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::topbar.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['active' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isGroupActive),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupIcon)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                        <?php echo e($groupLabel); ?>

                                     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e)): ?>
<?php $attributes = $__attributesOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e; ?>
<?php unset($__attributesOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e)): ?>
<?php $component = $__componentOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e; ?>
<?php unset($__componentOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e); ?>
<?php endif; ?>
                                 <?php $__env->endSlot(); ?>

                                <?php
                                    $lists = [];

                                    foreach ($group->getItems() as $item) {
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
                                                $isItemActive = $item->isActive();
                                                $itemBadge = $item->getBadge();
                                                $itemBadgeColor = $item->getBadgeColor($itemBadge);
                                                $itemBadgeTooltip = $item->getBadgeTooltip($itemBadge);
                                                $itemUrl = $item->getUrl();
                                                $itemIcon = $isItemActive ? ($item->getActiveIcon() ?? $item->getIcon()) : $item->getIcon();
                                                $shouldItemOpenUrlInNewTab = $item->shouldOpenUrlInNewTab();
                                                $itemExtraAttributes = $item->getExtraAttributeBag();
                                            ?>

                                            <?php if (isset($component)) { $__componentOriginalc67084d1e4215442fec660f0d8cd4ec60e9c2f12bc4932156a12c753ba32c47e = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalc67084d1e4215442fec660f0d8cd4ec60e9c2f12bc4932156a12c753ba32c47e = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.dropdown.list.item','data' => ['badge' => $itemBadge,'badgeColor' => $itemBadgeColor,'badgeTooltip' => $itemBadgeTooltip,'color' => $isItemActive ? 'primary' : 'gray','href' => $itemUrl,'icon' => $itemIcon,'tag' => 'a','target' => $shouldItemOpenUrlInNewTab ? '_blank' : null,'ariaCurrent' => $isItemActive ? 'page' : null,'attributes' => \Filament\Support\prepare_inherited_attributes($itemExtraAttributes)]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::dropdown.list.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['badge' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadge),'badge-color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadgeColor),'badge-tooltip' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadgeTooltip),'color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isItemActive ? 'primary' : 'gray'),'href' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemUrl),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemIcon),'tag' => 'a','target' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($shouldItemOpenUrlInNewTab ? '_blank' : null),'aria-current' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isItemActive ? 'page' : null),'attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(\Filament\Support\prepare_inherited_attributes($itemExtraAttributes))]); ?>
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
                        <?php else: ?>
                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $group->getItems(); $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $item): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                <?php
                                    $isItemActive = $item->isActive();
                                    $itemActiveIcon = $item->getActiveIcon();
                                    $itemBadge = $item->getBadge();
                                    $itemBadgeColor = $item->getBadgeColor($itemBadge);
                                    $itemBadgeTooltip = $item->getBadgeTooltip($itemBadge);
                                    $itemIcon = $item->getIcon();
                                    $shouldItemOpenUrlInNewTab = $item->shouldOpenUrlInNewTab();
                                    $itemUrl = $item->getUrl();
                                    $itemExtraAttributes = $item->getExtraAttributeBag();
                                ?>

                                <?php if (isset($component)) { $__componentOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.topbar.item','data' => ['active' => $isItemActive,'activeIcon' => $itemActiveIcon,'badge' => $itemBadge,'badgeColor' => $itemBadgeColor,'badgeTooltip' => $itemBadgeTooltip,'icon' => $itemIcon,'shouldOpenUrlInNewTab' => $shouldItemOpenUrlInNewTab,'url' => $itemUrl,'attributes' => \Filament\Support\prepare_inherited_attributes($itemExtraAttributes)]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::topbar.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['active' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isItemActive),'active-icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemActiveIcon),'badge' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadge),'badge-color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadgeColor),'badge-tooltip' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemBadgeTooltip),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemIcon),'should-open-url-in-new-tab' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($shouldItemOpenUrlInNewTab),'url' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemUrl),'attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(\Filament\Support\prepare_inherited_attributes($itemExtraAttributes))]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                    <?php echo e($item->getLabel()); ?>

                                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e)): ?>
<?php $attributes = $__attributesOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e; ?>
<?php unset($__attributesOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e)): ?>
<?php $component = $__componentOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e; ?>
<?php unset($__componentOriginal04d1c735f4adb7a89e18ebc624027c15b88b7402e1acff4d49427100240be74e); ?>
<?php endif; ?>
                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                </ul>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <div
            <?php if($hasTenancy): ?>
                x-persist="topbar.end.panel-<?php echo e(filament()->getId()); ?>.tenant-<?php echo e(filament()->getTenant()?->getKey()); ?>"
            <?php else: ?>
                x-persist="topbar.end.panel-<?php echo e(filament()->getId()); ?>"
            <?php endif; ?>
            class="fi-topbar-end"
        >
            <?php echo e(FilamentView::renderHook(PanelsRenderHook::GLOBAL_SEARCH_BEFORE)); ?>


            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filament()->isGlobalSearchEnabled() && filament()->getGlobalSearchPosition() === GlobalSearchPosition::Topbar): ?>
                <?php
$__split = function ($name, $params = []) {
    return [$name, $params];
};
[$__name, $__params] = $__split(GlobalSearch::class);

$__keyOuter = $__key ?? null;

$__key = null;
$__componentSlots = [];

$__key ??= \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::generateKey('lw-1411410434-0', $__key);

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

            <?php echo e(FilamentView::renderHook(PanelsRenderHook::GLOBAL_SEARCH_AFTER)); ?>


            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filament()->auth()->check()): ?>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filament()->hasDatabaseNotifications() && filament()->getDatabaseNotificationsPosition() === DatabaseNotificationsPosition::Topbar): ?>
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

$__key ??= \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::generateKey('lw-1411410434-0', $__key);

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

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filament()->hasUserMenu() && filament()->getUserMenuPosition() === UserMenuPosition::Topbar): ?>
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
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        </div>

        <?php echo e(FilamentView::renderHook(PanelsRenderHook::TOPBAR_END)); ?>

    </nav>

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
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/livewire/topbar.blade.php ENDPATH**/ ?>