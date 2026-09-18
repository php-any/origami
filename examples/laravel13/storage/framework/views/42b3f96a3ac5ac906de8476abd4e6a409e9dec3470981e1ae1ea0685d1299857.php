<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'position' => null,
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
    'position' => null,
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
    use Filament\Enums\UserMenuPosition;
    use Filament\Support\Facades\FilamentView;
    use Filament\Support\Icons\Heroicon;
    use Filament\Support\View\ComponentAttributeBag;
    use Filament\View\PanelsIconAlias;
    use Filament\View\PanelsRenderHook;
    use Illuminate\Support\Arr;

    $user = filament()->auth()->user();

    $userName = filament()->getUserName($user);

    $items = $this->getUserMenuItems();

    $itemsBeforeAndAfterThemeSwitcher = collect($items)
        ->groupBy(fn (Action $item): bool => $item->getSort() < 0, preserveKeys: true)
        ->all();
    $itemsBeforeThemeSwitcher = $itemsBeforeAndAfterThemeSwitcher[true] ?? collect();
    $itemsAfterThemeSwitcher = $itemsBeforeAndAfterThemeSwitcher[false] ?? collect();

    $hasProfileHeader = $itemsBeforeThemeSwitcher->has('profile') &&
        blank(($item = Arr::first($itemsBeforeThemeSwitcher))->getUrl()) &&
        (! $item->hasAction());

    if ($itemsBeforeThemeSwitcher->has('profile')) {
        $itemsBeforeThemeSwitcher = $itemsBeforeThemeSwitcher->prepend($itemsBeforeThemeSwitcher->pull('profile'), 'profile');
    }

    $multiGroupAfterTheme = $this->hasMultipleUserMenuItemGroups();
    $afterThemeItemGroups = $multiGroupAfterTheme ? $this->getUserMenuItemGroupsAfterTheme() : [];

    $position ??= filament()->getUserMenuPosition();

    $isSidebarCollapsibleOnDesktop = filament()->isSidebarCollapsibleOnDesktop();
?>

<?php echo e(FilamentView::renderHook(PanelsRenderHook::USER_MENU_BEFORE)); ?>


<?php if (isset($component)) { $__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.dropdown.index','data' => ['placement' => ($position === UserMenuPosition::Topbar) ? 'bottom-end' : 'top-end','teleport' => $position === UserMenuPosition::Topbar,'attributes' => 
        \Filament\Support\prepare_inherited_attributes($attributes)
            ->class(['fi-user-menu'])
    ]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::dropdown'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['placement' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(($position === UserMenuPosition::Topbar) ? 'bottom-end' : 'top-end'),'teleport' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($position === UserMenuPosition::Topbar),'attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
        \Filament\Support\prepare_inherited_attributes($attributes)
            ->class(['fi-user-menu'])
    )]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

     <?php $__env->slot('trigger', null, []); ?> 
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($position === UserMenuPosition::Topbar): ?>
            <button
                aria-label="<?php echo e(__('filament-panels::layout.actions.open_user_menu.label')); ?>"
                type="button"
                class="fi-user-menu-trigger"
            >
                <?php if (isset($component)) { $__componentOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.avatar.user','data' => ['user' => $user,'loading' => 'lazy']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::avatar.user'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['user' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($user),'loading' => 'lazy']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4)): ?>
<?php $attributes = $__attributesOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4; ?>
<?php unset($__attributesOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4)): ?>
<?php $component = $__componentOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4; ?>
<?php unset($__componentOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4); ?>
<?php endif; ?>
            </button>
        <?php else: ?>
            <button
                aria-label="<?php echo e(filled($userName) ? $userName : __('filament-panels::layout.actions.open_user_menu.label')); ?>"
                type="button"
                class="fi-user-menu-trigger"
            >
                <?php if (isset($component)) { $__componentOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.avatar.user','data' => ['user' => $user,'loading' => 'lazy']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::avatar.user'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['user' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($user),'loading' => 'lazy']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4)): ?>
<?php $attributes = $__attributesOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4; ?>
<?php unset($__attributesOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4)): ?>
<?php $component = $__componentOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4; ?>
<?php unset($__componentOriginal1f46c54f2e4c4b4a5b9b9e5319d56b4855a492614eca21ea2afa61912dea45c4); ?>
<?php endif; ?>

                <span
                    <?php if($isSidebarCollapsibleOnDesktop): ?>
                        x-show="$store.sidebar.isOpen"
                    <?php endif; ?>
                    class="fi-user-menu-trigger-text"
                >
                    <?php echo e($userName); ?>

                </span>

                <?php echo e(\Filament\Support\generate_icon_html(Heroicon::ChevronUp, alias: PanelsIconAlias::USER_MENU_TOGGLE_BUTTON, attributes: new ComponentAttributeBag([
                        'x-show' => $isSidebarCollapsibleOnDesktop ? '$store.sidebar.isOpen' : null,
                    ]))); ?>

            </button>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
     <?php $__env->endSlot(); ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasProfileHeader): ?>
        <?php
            $item = $itemsBeforeThemeSwitcher['profile'];
            $itemColor = $item->getColor();
            $itemIcon = $item->getIcon();

            unset($itemsBeforeThemeSwitcher['profile']);
        ?>

        <?php echo e(FilamentView::renderHook(PanelsRenderHook::USER_MENU_PROFILE_BEFORE)); ?>


        <?php if (isset($component)) { $__componentOriginal2a4b984efad10636c8eb2e2d0666cfa6ddcf02cf2cbee91a407c355bce805d2b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal2a4b984efad10636c8eb2e2d0666cfa6ddcf02cf2cbee91a407c355bce805d2b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.dropdown.header','data' => ['color' => $itemColor,'icon' => $itemIcon]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::dropdown.header'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemColor),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($itemIcon)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

            <?php echo e($item->getLabel()); ?>

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

        <?php echo e(FilamentView::renderHook(PanelsRenderHook::USER_MENU_PROFILE_AFTER)); ?>

    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($itemsBeforeThemeSwitcher->isNotEmpty()): ?>
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

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $itemsBeforeThemeSwitcher; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $key => $item): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($key === 'profile'): ?>
                    <?php echo e(FilamentView::renderHook(PanelsRenderHook::USER_MENU_PROFILE_BEFORE)); ?>


                    <?php echo e($item); ?>


                    <?php echo e(FilamentView::renderHook(PanelsRenderHook::USER_MENU_PROFILE_AFTER)); ?>

                <?php else: ?>
                    <?php echo e($item); ?>

                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
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
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filament()->hasDarkMode() && (! filament()->hasDarkModeForced()) && filament()->hasThemeSwitcher()): ?>
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

            <?php if (isset($component)) { $__componentOriginale8dc36b8c734989d8a374c8dfa0e9b8ba8152889557e26f23ef59b3421cb06ac = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginale8dc36b8c734989d8a374c8dfa0e9b8ba8152889557e26f23ef59b3421cb06ac = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.theme-switcher.index','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::theme-switcher'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginale8dc36b8c734989d8a374c8dfa0e9b8ba8152889557e26f23ef59b3421cb06ac)): ?>
<?php $attributes = $__attributesOriginale8dc36b8c734989d8a374c8dfa0e9b8ba8152889557e26f23ef59b3421cb06ac; ?>
<?php unset($__attributesOriginale8dc36b8c734989d8a374c8dfa0e9b8ba8152889557e26f23ef59b3421cb06ac); ?>
<?php endif; ?>
<?php if (isset($__componentOriginale8dc36b8c734989d8a374c8dfa0e9b8ba8152889557e26f23ef59b3421cb06ac)): ?>
<?php $component = $__componentOriginale8dc36b8c734989d8a374c8dfa0e9b8ba8152889557e26f23ef59b3421cb06ac; ?>
<?php unset($__componentOriginale8dc36b8c734989d8a374c8dfa0e9b8ba8152889557e26f23ef59b3421cb06ac); ?>
<?php endif; ?>
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
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($multiGroupAfterTheme && $afterThemeItemGroups !== []): ?>
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $afterThemeItemGroups; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $afterThemeGroup): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
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

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $afterThemeGroup; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $key => $item): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($key === 'profile'): ?>
                        <?php echo e(FilamentView::renderHook(PanelsRenderHook::USER_MENU_PROFILE_BEFORE)); ?>


                        <?php echo e($item); ?>


                        <?php echo e(FilamentView::renderHook(PanelsRenderHook::USER_MENU_PROFILE_AFTER)); ?>

                    <?php else: ?>
                        <?php echo e($item); ?>

                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
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
    <?php elseif($itemsAfterThemeSwitcher->isNotEmpty()): ?>
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

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $itemsAfterThemeSwitcher; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $key => $item): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($key === 'profile'): ?>
                    <?php echo e(FilamentView::renderHook(PanelsRenderHook::USER_MENU_PROFILE_BEFORE)); ?>


                    <?php echo e($item); ?>


                    <?php echo e(FilamentView::renderHook(PanelsRenderHook::USER_MENU_PROFILE_AFTER)); ?>

                <?php else: ?>
                    <?php echo e($item); ?>

                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
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
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
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

<?php echo e(FilamentView::renderHook(PanelsRenderHook::USER_MENU_AFTER)); ?>

<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/user-menu.blade.php ENDPATH**/ ?>