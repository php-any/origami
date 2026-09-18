<?php if (isset($component)) { $__componentOriginal54cb83548d5b95790cb921ec15fffcee14a0ef30f5ad250dadde563d65b3cb93 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal54cb83548d5b95790cb921ec15fffcee14a0ef30f5ad250dadde563d65b3cb93 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.page.index','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::page'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

    <?php echo e($this->content); ?>

 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal54cb83548d5b95790cb921ec15fffcee14a0ef30f5ad250dadde563d65b3cb93)): ?>
<?php $attributes = $__attributesOriginal54cb83548d5b95790cb921ec15fffcee14a0ef30f5ad250dadde563d65b3cb93; ?>
<?php unset($__attributesOriginal54cb83548d5b95790cb921ec15fffcee14a0ef30f5ad250dadde563d65b3cb93); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal54cb83548d5b95790cb921ec15fffcee14a0ef30f5ad250dadde563d65b3cb93)): ?>
<?php $component = $__componentOriginal54cb83548d5b95790cb921ec15fffcee14a0ef30f5ad250dadde563d65b3cb93; ?>
<?php unset($__componentOriginal54cb83548d5b95790cb921ec15fffcee14a0ef30f5ad250dadde563d65b3cb93); ?>
<?php endif; ?>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/pages/page.blade.php ENDPATH**/ ?>