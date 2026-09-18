<?php
    use Filament\Support\Icons\Heroicon;
?>

<div
    x-data="{ theme: null }"
    x-init="
        $watch('theme', () => {
            $dispatch('theme-changed', theme)
        })

        theme = localStorage.getItem('theme') || <?php echo \Illuminate\Support\Js::from(filament()->getDefaultThemeMode()->value)->toHtml() ?>
    "
    role="group"
    aria-label="<?php echo e(__('filament-panels::layout.actions.theme_switcher.label')); ?>"
    class="fi-theme-switcher"
>
    <?php if (isset($component)) { $__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.theme-switcher.button','data' => ['icon' => Heroicon::Sun,'theme' => 'light']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::theme-switcher.button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(Heroicon::Sun),'theme' => 'light']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d)): ?>
<?php $attributes = $__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d; ?>
<?php unset($__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d)): ?>
<?php $component = $__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d; ?>
<?php unset($__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d); ?>
<?php endif; ?>

    <?php if (isset($component)) { $__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.theme-switcher.button','data' => ['icon' => Heroicon::Moon,'theme' => 'dark']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::theme-switcher.button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(Heroicon::Moon),'theme' => 'dark']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d)): ?>
<?php $attributes = $__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d; ?>
<?php unset($__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d)): ?>
<?php $component = $__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d; ?>
<?php unset($__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d); ?>
<?php endif; ?>

    <?php if (isset($component)) { $__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.theme-switcher.button','data' => ['icon' => Heroicon::ComputerDesktop,'theme' => 'system']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::theme-switcher.button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(Heroicon::ComputerDesktop),'theme' => 'system']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d)): ?>
<?php $attributes = $__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d; ?>
<?php unset($__attributesOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d)): ?>
<?php $component = $__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d; ?>
<?php unset($__componentOriginal11a13dfd9dee3cf70f4a741dcfc539e74a89341a2c5b22f65448644942355f9d); ?>
<?php endif; ?>
</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/theme-switcher/index.blade.php ENDPATH**/ ?>