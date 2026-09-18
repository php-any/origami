<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames((['routeParameters']));

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

foreach (array_filter((['routeParameters']), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<div class="flex flex-col gap-3">
    <h2 class="text-lg font-semibold">Routing parameters</h2>
    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($routeParameters): ?>
    <div class="bg-white dark:bg-white/[2%] border border-neutral-200 dark:border-neutral-800 rounded-md overflow-x-auto p-5 text-sm font-mono shadow-xs">
        <?php if (isset($component)) { $__componentOriginal60f51a3b3773d99d09e20b9ccfdbc33fd1b7917c98351e99d3b31ffcb65ce951 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal60f51a3b3773d99d09e20b9ccfdbc33fd1b7917c98351e99d3b31ffcb65ce951 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.syntax-highlight','data' => ['code' => $routeParameters,'language' => 'json']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::syntax-highlight'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['code' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($routeParameters),'language' => 'json']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal60f51a3b3773d99d09e20b9ccfdbc33fd1b7917c98351e99d3b31ffcb65ce951)): ?>
<?php $attributes = $__attributesOriginal60f51a3b3773d99d09e20b9ccfdbc33fd1b7917c98351e99d3b31ffcb65ce951; ?>
<?php unset($__attributesOriginal60f51a3b3773d99d09e20b9ccfdbc33fd1b7917c98351e99d3b31ffcb65ce951); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal60f51a3b3773d99d09e20b9ccfdbc33fd1b7917c98351e99d3b31ffcb65ce951)): ?>
<?php $component = $__componentOriginal60f51a3b3773d99d09e20b9ccfdbc33fd1b7917c98351e99d3b31ffcb65ce951; ?>
<?php unset($__componentOriginal60f51a3b3773d99d09e20b9ccfdbc33fd1b7917c98351e99d3b31ffcb65ce951); ?>
<?php endif; ?>
    </div>
    <?php else: ?>
    <?php if (isset($component)) { $__componentOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.empty-state','data' => ['message' => 'No routing parameters']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::empty-state'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['message' => 'No routing parameters']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56)): ?>
<?php $attributes = $__attributesOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56; ?>
<?php unset($__attributesOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56)): ?>
<?php $component = $__componentOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56; ?>
<?php unset($__componentOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56); ?>
<?php endif; ?>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
</div>
<?php /**PATH d:\gitcode.com\origami\examples\laravel13\vendor\laravel\framework\src\illuminate\foundation\providers/../resources/exceptions/renderer/components/routing-parameter.blade.php ENDPATH**/ ?>