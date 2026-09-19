<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'debounce' => '500ms',
    'onBlur' => false,
    'placeholder' => __('filament-tables::table.fields.search.placeholder'),
    'wireModel' => 'tableSearch',
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
    'debounce' => '500ms',
    'onBlur' => false,
    'placeholder' => __('filament-tables::table.fields.search.placeholder'),
    'wireModel' => 'tableSearch',
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Support\Icons\Heroicon;
    use Filament\Support\View\ComponentAttributeBag as FilamentComponentAttributeBag;
    use Filament\Tables\View\TablesIconAlias;

    $wireModelAttribute = $onBlur ? 'wire:model.live.blur' : "wire:model.live.debounce.{$debounce}";
?>

<div
    x-id="['input']"
    <?php echo e($attributes->class(['fi-ta-search-field'])); ?>

>
    <label x-bind:for="$id('input')" class="fi-sr-only">
        <?php echo e(__('filament-tables::table.fields.search.label')); ?>

    </label>

    <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => ['inlinePrefix' => true,'prefixIcon' => Heroicon::MagnifyingGlass,'prefixIconAlias' => TablesIconAlias::SEARCH_FIELD,'wire:target' => $wireModel]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['inline-prefix' => true,'prefix-icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(Heroicon::MagnifyingGlass),'prefix-icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(TablesIconAlias::SEARCH_FIELD),'wire:target' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($wireModel)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

        <?php if (isset($component)) { $__componentOriginal7c1d2ac545d23b624bb13365b4ddb78282f11c7aa6b4d158c17151cbbc7e684d = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal7c1d2ac545d23b624bb13365b4ddb78282f11c7aa6b4d158c17151cbbc7e684d = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.index','data' => ['attributes' => 
                (new FilamentComponentAttributeBag)->merge([
                    'autocomplete' => 'off',
                    'inlinePrefix' => true,
                    'maxlength' => 1000,
                    'placeholder' => $placeholder,
                    'type' => 'search',
                    'wire:key' => $this->getId() . '.table.' . $wireModel . '.field.input',
                    $wireModelAttribute => $wireModel,
                    'x-bind:id' => '$id(\'input\')',
                    'x-on:keyup' => 'if ($event.key === \'Enter\') { $wire.$refresh() }',
                ], escape: false)
            ]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
                (new FilamentComponentAttributeBag)->merge([
                    'autocomplete' => 'off',
                    'inlinePrefix' => true,
                    'maxlength' => 1000,
                    'placeholder' => $placeholder,
                    'type' => 'search',
                    'wire:key' => $this->getId() . '.table.' . $wireModel . '.field.input',
                    $wireModelAttribute => $wireModel,
                    'x-bind:id' => '$id(\'input\')',
                    'x-on:keyup' => 'if ($event.key === \'Enter\') { $wire.$refresh() }',
                ], escape: false)
            )]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal7c1d2ac545d23b624bb13365b4ddb78282f11c7aa6b4d158c17151cbbc7e684d)): ?>
<?php $attributes = $__attributesOriginal7c1d2ac545d23b624bb13365b4ddb78282f11c7aa6b4d158c17151cbbc7e684d; ?>
<?php unset($__attributesOriginal7c1d2ac545d23b624bb13365b4ddb78282f11c7aa6b4d158c17151cbbc7e684d); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal7c1d2ac545d23b624bb13365b4ddb78282f11c7aa6b4d158c17151cbbc7e684d)): ?>
<?php $component = $__componentOriginal7c1d2ac545d23b624bb13365b4ddb78282f11c7aa6b4d158c17151cbbc7e684d; ?>
<?php unset($__componentOriginal7c1d2ac545d23b624bb13365b4ddb78282f11c7aa6b4d158c17151cbbc7e684d); ?>
<?php endif; ?>
     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $attributes = $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $component = $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\tables\resources\views/components/search-field.blade.php ENDPATH**/ ?>