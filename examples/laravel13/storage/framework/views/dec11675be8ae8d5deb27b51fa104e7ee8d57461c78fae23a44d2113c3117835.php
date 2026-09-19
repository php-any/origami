<?php if (isset($component)) { $__componentOriginal8b442205e0895db478ee2f1b6137499c89aa2eec9307ae7f732056114c161ba9 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal8b442205e0895db478ee2f1b6137499c89aa2eec9307ae7f732056114c161ba9 = $attributes; } ?>
<?php $component = Filament\View\LegacyComponents\PageComponent::resolve([] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::page'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Filament\View\LegacyComponents\PageComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

    <?php echo e($this->form); ?>

 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal8b442205e0895db478ee2f1b6137499c89aa2eec9307ae7f732056114c161ba9)): ?>
<?php $attributes = $__attributesOriginal8b442205e0895db478ee2f1b6137499c89aa2eec9307ae7f732056114c161ba9; ?>
<?php unset($__attributesOriginal8b442205e0895db478ee2f1b6137499c89aa2eec9307ae7f732056114c161ba9); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal8b442205e0895db478ee2f1b6137499c89aa2eec9307ae7f732056114c161ba9)): ?>
<?php $component = $__componentOriginal8b442205e0895db478ee2f1b6137499c89aa2eec9307ae7f732056114c161ba9; ?>
<?php unset($__componentOriginal8b442205e0895db478ee2f1b6137499c89aa2eec9307ae7f732056114c161ba9); ?>
<?php endif; ?>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\resources\views/filament/pages/manage-settings.blade.php ENDPATH**/ ?>