<?php
    use Filament\Support\Icons\Heroicon;
    use Filament\View\PanelsIconAlias;
    use Illuminate\Support\Number;
?>

<?php if (isset($component)) { $__componentOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6a0945a74540aaea7ddf99e4114c1a1563f1560db549e9045a0d0f751d39b5ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.icon-button','data' => ['badge' => $unreadNotificationsCount ?: null,'color' => 'gray','icon' => Heroicon::OutlinedBell,'iconAlias' => PanelsIconAlias::TOPBAR_OPEN_DATABASE_NOTIFICATIONS_BUTTON,'iconSize' => 'lg','label' => 
        $unreadNotificationsCount
        ? trans_choice('filament-panels::layout.actions.open_database_notifications.label_with_unread_count', $unreadNotificationsCount, ['count' => Number::format($unreadNotificationsCount, locale: app()->getLocale())])
        : __('filament-panels::layout.actions.open_database_notifications.label')
    ,'class' => 'fi-topbar-database-notifications-btn']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::icon-button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['badge' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($unreadNotificationsCount ?: null),'color' => 'gray','icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(Heroicon::OutlinedBell),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(PanelsIconAlias::TOPBAR_OPEN_DATABASE_NOTIFICATIONS_BUTTON),'icon-size' => 'lg','label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
        $unreadNotificationsCount
        ? trans_choice('filament-panels::layout.actions.open_database_notifications.label_with_unread_count', $unreadNotificationsCount, ['count' => Number::format($unreadNotificationsCount, locale: app()->getLocale())])
        : __('filament-panels::layout.actions.open_database_notifications.label')
    ),'class' => 'fi-topbar-database-notifications-btn']); ?>
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
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/topbar/database-notifications-trigger.blade.php ENDPATH**/ ?>