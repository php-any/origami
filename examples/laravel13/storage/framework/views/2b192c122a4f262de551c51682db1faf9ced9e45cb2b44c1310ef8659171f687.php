<?php
    use Filament\Support\Facades\FilamentView;
    use Filament\Widgets\View\WidgetsRenderHook;
?>

<?php if (isset($component)) { $__componentOriginaleba50eae279e8730a71fb5403943b02af405b4ca7174fe95cfd3cc564939da96 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginaleba50eae279e8730a71fb5403943b02af405b4ca7174fe95cfd3cc564939da96 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-widgets::components.widget','data' => ['class' => 'fi-wi-table']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-widgets::widget'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'fi-wi-table']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

    <?php echo e(FilamentView::renderHook(WidgetsRenderHook::TABLE_WIDGET_START, scopes: static::class)); ?>


    <?php echo e($this->table ?? null); ?>


    <?php echo e(FilamentView::renderHook(WidgetsRenderHook::TABLE_WIDGET_END, scopes: static::class)); ?>

 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginaleba50eae279e8730a71fb5403943b02af405b4ca7174fe95cfd3cc564939da96)): ?>
<?php $attributes = $__attributesOriginaleba50eae279e8730a71fb5403943b02af405b4ca7174fe95cfd3cc564939da96; ?>
<?php unset($__attributesOriginaleba50eae279e8730a71fb5403943b02af405b4ca7174fe95cfd3cc564939da96); ?>
<?php endif; ?>
<?php if (isset($__componentOriginaleba50eae279e8730a71fb5403943b02af405b4ca7174fe95cfd3cc564939da96)): ?>
<?php $component = $__componentOriginaleba50eae279e8730a71fb5403943b02af405b4ca7174fe95cfd3cc564939da96; ?>
<?php unset($__componentOriginaleba50eae279e8730a71fb5403943b02af405b4ca7174fe95cfd3cc564939da96); ?>
<?php endif; ?>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\widgets\resources\views/table-widget.blade.php ENDPATH**/ ?>