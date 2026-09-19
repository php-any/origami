<?php
    use Filament\Support\View\ComponentAttributeBag;

    $columns = $this->getColumns();
    $pollingInterval = $this->getPollingInterval();

    $heading = $this->getHeading();
    $description = $this->getDescription();
    $hasHeading = filled($heading);
    $hasDescription = filled($description);
?>

<?php if (isset($component)) { $__componentOriginaleba50eae279e8730a71fb5403943b02af405b4ca7174fe95cfd3cc564939da96 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginaleba50eae279e8730a71fb5403943b02af405b4ca7174fe95cfd3cc564939da96 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-widgets::components.widget','data' => ['attributes' => 
        (new ComponentAttributeBag)
            ->merge([
                'wire:poll.' . $pollingInterval => $pollingInterval ? true : null,
            ], escape: false)
            ->class([
                'fi-wi-stats-overview',
            ])
    ]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-widgets::widget'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
        (new ComponentAttributeBag)
            ->merge([
                'wire:poll.' . $pollingInterval => $pollingInterval ? true : null,
            ], escape: false)
            ->class([
                'fi-wi-stats-overview',
            ])
    )]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

    <?php echo e($this->content); ?>

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
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\widgets\resources\views/stats-overview-widget.blade.php ENDPATH**/ ?>