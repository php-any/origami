<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'user' => filament()->auth()->user(),
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
    'user' => filament()->auth()->user(),
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    $src = filament()->getUserAvatarUrl($user);
    $alt = __('filament-panels::layout.avatar.alt', ['name' => filament()->getUserName($user)]);
?>

<?php if (isset($component)) { $__componentOriginalc3796ce1b1220cc603cce7c0f5541bb06d8eb06b4254b24403f01218dfb57033 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalc3796ce1b1220cc603cce7c0f5541bb06d8eb06b4254b24403f01218dfb57033 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.avatar','data' => ['src' => $src,'alt' => $alt,'attributes' => 
        \Filament\Support\prepare_inherited_attributes($attributes)
            ->class(['fi-user-avatar'])
    ]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::avatar'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['src' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($src),'alt' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($alt),'attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
        \Filament\Support\prepare_inherited_attributes($attributes)
            ->class(['fi-user-avatar'])
    )]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalc3796ce1b1220cc603cce7c0f5541bb06d8eb06b4254b24403f01218dfb57033)): ?>
<?php $attributes = $__attributesOriginalc3796ce1b1220cc603cce7c0f5541bb06d8eb06b4254b24403f01218dfb57033; ?>
<?php unset($__attributesOriginalc3796ce1b1220cc603cce7c0f5541bb06d8eb06b4254b24403f01218dfb57033); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalc3796ce1b1220cc603cce7c0f5541bb06d8eb06b4254b24403f01218dfb57033)): ?>
<?php $component = $__componentOriginalc3796ce1b1220cc603cce7c0f5541bb06d8eb06b4254b24403f01218dfb57033; ?>
<?php unset($__componentOriginalc3796ce1b1220cc603cce7c0f5541bb06d8eb06b4254b24403f01218dfb57033); ?>
<?php endif; ?>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/components/avatar/user.blade.php ENDPATH**/ ?>