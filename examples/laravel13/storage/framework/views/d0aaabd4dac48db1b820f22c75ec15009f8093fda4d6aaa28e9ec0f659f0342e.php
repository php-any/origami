<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'fullHeight' => false,
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
    'fullHeight' => false,
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Pages\Enums\SubNavigationPosition;
    use Filament\Support\Facades\FilamentView;
    use Filament\Tables\Contracts\HasTable;
    use Filament\View\PanelsRenderHook;
    use Illuminate\Contracts\Support\Htmlable;

    $subNavigation = $this->getCachedSubNavigation();
    $subNavigationPosition = $this->getSubNavigationPosition();
    $widgetData = $this->getWidgetData();
?>

<div
    <?php echo e($attributes->class([
            'fi-page',
            'fi-height-full' => $fullHeight,
            'fi-page-has-sub-navigation' => $subNavigation,
            "fi-page-has-sub-navigation-{$subNavigationPosition->value}" => $subNavigation,
            ...$this->getPageClasses(),
        ])); ?>

>
    <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_START, scopes: $this->getRenderHookScopes())); ?>


    <div class="fi-page-header-main-ctn">
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($subNavigation): ?>
            <div
                class="fi-page-main-sub-navigation-mobile-menu-render-hook-ctn"
            >
                <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_SUB_NAVIGATION_MOBILE_MENU_BEFORE, scopes: $this->getRenderHookScopes())); ?>

            </div>

            <?php if (isset($component)) { $__componentOriginal63e7686290c02493ecdd24c7f8d5abed49f29cd9b7667ca32d1d8b70ed1c905f = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal63e7686290c02493ecdd24c7f8d5abed49f29cd9b7667ca32d1d8b70ed1c905f = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.page.sub-navigation.mobile-menu','data' => ['navigation' => $subNavigation]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::page.sub-navigation.mobile-menu'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['navigation' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($subNavigation)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal63e7686290c02493ecdd24c7f8d5abed49f29cd9b7667ca32d1d8b70ed1c905f)): ?>
<?php $attributes = $__attributesOriginal63e7686290c02493ecdd24c7f8d5abed49f29cd9b7667ca32d1d8b70ed1c905f; ?>
<?php unset($__attributesOriginal63e7686290c02493ecdd24c7f8d5abed49f29cd9b7667ca32d1d8b70ed1c905f); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal63e7686290c02493ecdd24c7f8d5abed49f29cd9b7667ca32d1d8b70ed1c905f)): ?>
<?php $component = $__componentOriginal63e7686290c02493ecdd24c7f8d5abed49f29cd9b7667ca32d1d8b70ed1c905f; ?>
<?php unset($__componentOriginal63e7686290c02493ecdd24c7f8d5abed49f29cd9b7667ca32d1d8b70ed1c905f); ?>
<?php endif; ?>

            <div
                class="fi-page-main-sub-navigation-mobile-menu-render-hook-ctn"
            >
                <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_SUB_NAVIGATION_MOBILE_MENU_AFTER, scopes: $this->getRenderHookScopes())); ?>

            </div>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($header = $this->getHeader()): ?>
            <?php echo e($header); ?>

        <?php else: ?>
            <?php
                $heading = $this->getHeading();
                $headerActions = $this->getCachedHeaderActions();
                $headerActionsAlignment = $this->getHeaderActionsAlignment();
                $breadcrumbs = filament()->hasBreadcrumbs() ? $this->getBreadcrumbs() : [];
                $subheading = $this->getSubheading();
            ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($headerActions) || $breadcrumbs || filled($heading) || filled($subheading)): ?>
                <?php if (isset($component)) { $__componentOriginal2cf9bbeacf3fad7469f9726e90855f5419361e0b0306976f31e57eb64b5eb340 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal2cf9bbeacf3fad7469f9726e90855f5419361e0b0306976f31e57eb64b5eb340 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.header.index','data' => ['actions' => $headerActions,'actionsAlignment' => $headerActionsAlignment,'breadcrumbs' => $breadcrumbs,'heading' => $heading,'subheading' => $subheading]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::header'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['actions' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($headerActions),'actions-alignment' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($headerActionsAlignment),'breadcrumbs' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($breadcrumbs),'heading' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($heading),'subheading' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($subheading)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($heading instanceof Htmlable): ?>
                         <?php $__env->slot('heading', null, []); ?> 
                            <?php echo e($heading); ?>

                         <?php $__env->endSlot(); ?>
                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($subheading instanceof Htmlable): ?>
                         <?php $__env->slot('subheading', null, []); ?> 
                            <?php echo e($subheading); ?>

                         <?php $__env->endSlot(); ?>
                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal2cf9bbeacf3fad7469f9726e90855f5419361e0b0306976f31e57eb64b5eb340)): ?>
<?php $attributes = $__attributesOriginal2cf9bbeacf3fad7469f9726e90855f5419361e0b0306976f31e57eb64b5eb340; ?>
<?php unset($__attributesOriginal2cf9bbeacf3fad7469f9726e90855f5419361e0b0306976f31e57eb64b5eb340); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal2cf9bbeacf3fad7469f9726e90855f5419361e0b0306976f31e57eb64b5eb340)): ?>
<?php $component = $__componentOriginal2cf9bbeacf3fad7469f9726e90855f5419361e0b0306976f31e57eb64b5eb340; ?>
<?php unset($__componentOriginal2cf9bbeacf3fad7469f9726e90855f5419361e0b0306976f31e57eb64b5eb340); ?>
<?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <div class="fi-page-main">
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($subNavigation): ?>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($subNavigationPosition === SubNavigationPosition::Start): ?>
                    <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_SUB_NAVIGATION_START_BEFORE, scopes: $this->getRenderHookScopes())); ?>


                    <?php if (isset($component)) { $__componentOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.page.sub-navigation.sidebar','data' => ['navigation' => $subNavigation]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::page.sub-navigation.sidebar'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['navigation' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($subNavigation)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771)): ?>
<?php $attributes = $__attributesOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771; ?>
<?php unset($__attributesOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771)): ?>
<?php $component = $__componentOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771; ?>
<?php unset($__componentOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771); ?>
<?php endif; ?>

                    <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_SUB_NAVIGATION_START_AFTER, scopes: $this->getRenderHookScopes())); ?>

                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($subNavigationPosition === SubNavigationPosition::Top): ?>
                    <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_SUB_NAVIGATION_TOP_BEFORE, scopes: $this->getRenderHookScopes())); ?>


                    <?php if (isset($component)) { $__componentOriginalfd4be6ec5c31f08491a3a6d5c86848b3247d3736a3ad76d3ca22fcd6c2ce7f33 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalfd4be6ec5c31f08491a3a6d5c86848b3247d3736a3ad76d3ca22fcd6c2ce7f33 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.page.sub-navigation.tabs','data' => ['navigation' => $subNavigation]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::page.sub-navigation.tabs'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['navigation' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($subNavigation)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalfd4be6ec5c31f08491a3a6d5c86848b3247d3736a3ad76d3ca22fcd6c2ce7f33)): ?>
<?php $attributes = $__attributesOriginalfd4be6ec5c31f08491a3a6d5c86848b3247d3736a3ad76d3ca22fcd6c2ce7f33; ?>
<?php unset($__attributesOriginalfd4be6ec5c31f08491a3a6d5c86848b3247d3736a3ad76d3ca22fcd6c2ce7f33); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalfd4be6ec5c31f08491a3a6d5c86848b3247d3736a3ad76d3ca22fcd6c2ce7f33)): ?>
<?php $component = $__componentOriginalfd4be6ec5c31f08491a3a6d5c86848b3247d3736a3ad76d3ca22fcd6c2ce7f33; ?>
<?php unset($__componentOriginalfd4be6ec5c31f08491a3a6d5c86848b3247d3736a3ad76d3ca22fcd6c2ce7f33); ?>
<?php endif; ?>

                    <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_SUB_NAVIGATION_TOP_AFTER, scopes: $this->getRenderHookScopes())); ?>

                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <div class="fi-page-content">
                <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_HEADER_WIDGETS_BEFORE, scopes: $this->getRenderHookScopes())); ?>


                <?php echo e($this->headerWidgets); ?>


                <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_HEADER_WIDGETS_AFTER, scopes: $this->getRenderHookScopes())); ?>


                <?php echo e($slot); ?>


                <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_FOOTER_WIDGETS_BEFORE, scopes: $this->getRenderHookScopes())); ?>


                <?php echo e($this->footerWidgets); ?>


                <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_FOOTER_WIDGETS_AFTER, scopes: $this->getRenderHookScopes())); ?>

            </div>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($subNavigation && $subNavigationPosition === SubNavigationPosition::End): ?>
                <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_SUB_NAVIGATION_END_BEFORE, scopes: $this->getRenderHookScopes())); ?>


                <?php if (isset($component)) { $__componentOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.page.sub-navigation.sidebar','data' => ['navigation' => $subNavigation]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::page.sub-navigation.sidebar'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['navigation' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($subNavigation)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771)): ?>
<?php $attributes = $__attributesOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771; ?>
<?php unset($__attributesOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771)): ?>
<?php $component = $__componentOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771; ?>
<?php unset($__componentOriginalc4266be4cc0812dcdd968af9a0860034a3951215d320ea0f92617e79a6954771); ?>
<?php endif; ?>

                <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_SUB_NAVIGATION_END_AFTER, scopes: $this->getRenderHookScopes())); ?>

            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        </div>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($footer = $this->getFooter()): ?>
            <?php echo e($footer); ?>

        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    </div>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! ($this instanceof HasTable)): ?>
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
    <?php elseif($this->isTableLoaded() && filled($this->defaultTableAction)): ?>
        <div
            wire:init="mountAction(<?php echo \Illuminate\Support\Js::from($this->defaultTableAction)->toHtml() ?> , <?php if(filled($this->defaultTableActionArguments)): ?> <?php echo \Illuminate\Support\Js::from($this->defaultTableActionArguments)->toHtml() ?> <?php else: ?> {} <?php endif; ?> , <?php echo \Illuminate\Support\Js::from($this->getDefaultTableActionUrlContext())->toHtml() ?>)"
        ></div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($this->defaultAction)): ?>
        <div
            wire:init="mountAction(<?php echo \Illuminate\Support\Js::from($this->defaultAction)->toHtml() ?> , <?php if(filled($this->defaultActionArguments)): ?> <?php echo \Illuminate\Support\Js::from($this->defaultActionArguments)->toHtml() ?> <?php else: ?> {} <?php endif; ?> , <?php echo \Illuminate\Support\Js::from($this->getDefaultActionUrlContext())->toHtml() ?>)"
        ></div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php echo e(FilamentView::renderHook(PanelsRenderHook::PAGE_END, scopes: $this->getRenderHookScopes())); ?>


    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(method_exists($this, 'hasUnsavedDataChangesAlert') && $this->hasUnsavedDataChangesAlert()): ?>
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(FilamentView::hasSpaMode()): ?>
                            <?php
                    $__scriptKey = '485869876-0';
                    ob_start();
                ?>
                <script>
                    setUpSpaModeUnsavedDataChangesAlert({
                        body: <?php echo \Illuminate\Support\Js::from(__('filament-panels::unsaved-changes-alert.body'))->toHtml() ?>,
                        resolveLivewireComponentUsing: () => window.Livewire.find('<?php echo e($_instance->getId()); ?>'),
                        $wire,
                    })
                </script>
                            <?php
                    $__output = ob_get_clean();

                    \Livewire\store($this)->push('scripts', $__output, $__scriptKey)
                ?>
        <?php else: ?>
                            <?php
                    $__scriptKey = '485869876-0';
                    ob_start();
                ?>
                <script>
                    setUpUnsavedDataChangesAlert({ $wire })
                </script>
                            <?php
                    $__output = ob_get_clean();

                    \Livewire\store($this)->push('scripts', $__output, $__scriptKey)
                ?>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! app()->hasDebugModeEnabled()): ?>
                        <?php
                    $__scriptKey = '485869876-0';
                    ob_start();
                ?>
            <script>
                window.filamentErrorNotifications = <?php echo \Illuminate\Support\Js::from($this->hasErrorNotifications() ? $this->getErrorNotifications() : null)->toHtml() ?>
            </script>
                        <?php
                    $__output = ob_get_clean();

                    \Livewire\store($this)->push('scripts', $__output, $__scriptKey)
                ?>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if (isset($component)) { $__componentOriginalfd434aaf7608fcd7e252434e816a8d09683ebe9c92136b442843b67f21a36120 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalfd434aaf7608fcd7e252434e816a8d09683ebe9c92136b442843b67f21a36120 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.unsaved-action-changes-alert','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::unsaved-action-changes-alert'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalfd434aaf7608fcd7e252434e816a8d09683ebe9c92136b442843b67f21a36120)): ?>
<?php $attributes = $__attributesOriginalfd434aaf7608fcd7e252434e816a8d09683ebe9c92136b442843b67f21a36120; ?>
<?php unset($__attributesOriginalfd434aaf7608fcd7e252434e816a8d09683ebe9c92136b442843b67f21a36120); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalfd434aaf7608fcd7e252434e816a8d09683ebe9c92136b442843b67f21a36120)): ?>
<?php $component = $__componentOriginalfd434aaf7608fcd7e252434e816a8d09683ebe9c92136b442843b67f21a36120; ?>
<?php unset($__componentOriginalfd434aaf7608fcd7e252434e816a8d09683ebe9c92136b442843b67f21a36120); ?>
<?php endif; ?>
</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/page/index.blade.php ENDPATH**/ ?>