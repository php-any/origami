<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'heading' => null,
    'subheading' => null,
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
    'heading' => null,
    'subheading' => null,
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Support\Facades\FilamentView;
    use Filament\Tables\Contracts\HasTable;
    use Filament\View\PanelsRenderHook;

    $heading ??= $this->getHeading();
    $subheading ??= $this->getSubHeading();
    $hasLogo = $this->hasLogo();
?>

<div <?php echo e($attributes->class(['fi-simple-page'])); ?>>
    <?php echo e(FilamentView::renderHook(PanelsRenderHook::SIMPLE_PAGE_START, scopes: $this->getRenderHookScopes())); ?>


    <div class="fi-simple-page-content">
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($heading) || $hasLogo || filled($subheading)): ?>
            <?php if (isset($component)) { $__componentOriginal531b9832e42ead0dd7db94ab3e25060b90d57385d97b91198147350dc1b92617 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal531b9832e42ead0dd7db94ab3e25060b90d57385d97b91198147350dc1b92617 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-panels::components.header.simple','data' => ['heading' => $heading,'logo' => $hasLogo,'subheading' => $subheading]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-panels::header.simple'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['heading' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($heading),'logo' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasLogo),'subheading' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($subheading)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal531b9832e42ead0dd7db94ab3e25060b90d57385d97b91198147350dc1b92617)): ?>
<?php $attributes = $__attributesOriginal531b9832e42ead0dd7db94ab3e25060b90d57385d97b91198147350dc1b92617; ?>
<?php unset($__attributesOriginal531b9832e42ead0dd7db94ab3e25060b90d57385d97b91198147350dc1b92617); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal531b9832e42ead0dd7db94ab3e25060b90d57385d97b91198147350dc1b92617)): ?>
<?php $component = $__componentOriginal531b9832e42ead0dd7db94ab3e25060b90d57385d97b91198147350dc1b92617; ?>
<?php unset($__componentOriginal531b9832e42ead0dd7db94ab3e25060b90d57385d97b91198147350dc1b92617); ?>
<?php endif; ?>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php echo e($slot); ?>

    </div>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $this instanceof HasTable): ?>
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
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php echo e(FilamentView::renderHook(PanelsRenderHook::SIMPLE_PAGE_END, scopes: $this->getRenderHookScopes())); ?>

</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/page/simple.blade.php ENDPATH**/ ?>