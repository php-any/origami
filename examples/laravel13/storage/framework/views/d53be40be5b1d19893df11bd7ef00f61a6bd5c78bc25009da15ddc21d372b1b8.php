<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames((['exception', 'request']));

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

foreach (array_filter((['exception', 'request']), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<div
    x-data="{
        copied: false,
        async copyToClipboard() {
            try {
                await window.copyToClipboard('<?php echo e($request->fullUrl()); ?>');
                this.copied = true;
                setTimeout(() => { this.copied = false }, 3000);
            } catch (err) {
                console.error('Failed to copy the requestURL: ', err);
            }
        }
    }"
    <?php echo e($attributes->merge(['class' => "bg-white dark:bg-[#1a1a1a] border border-neutral-200 dark:border-white/10 rounded-lg flex items-center justify-between h-10 px-2 shadow-xs"])); ?>

>
    <div class="flex items-center gap-3 w-full">
        <?php if (isset($component)) { $__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.badge','data' => ['type' => 'error','variant' => 'solid']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::badge'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['type' => 'error','variant' => 'solid']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

            <?php if (isset($component)) { $__componentOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.alert','data' => ['class' => 'w-2.5 h-2.5']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.alert'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-2.5 h-2.5']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a)): ?>
<?php $attributes = $__attributesOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a; ?>
<?php unset($__attributesOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a)): ?>
<?php $component = $__componentOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a; ?>
<?php unset($__componentOriginal83c55a1f31ae13629ffaf99786ab1596c6576330747cfa30e650e33da161666a); ?>
<?php endif; ?>
            <?php echo e($exception->httpStatusCode()); ?>

         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad)): ?>
<?php $attributes = $__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad; ?>
<?php unset($__attributesOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad)): ?>
<?php $component = $__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad; ?>
<?php unset($__componentOriginal877872489ce5868e35a3f4f0b1eed7462a0d76b59febe9d6d430994440c9a1ad); ?>
<?php endif; ?>
        <?php if (isset($component)) { $__componentOriginal4630a041cc6330c0d3de6055dbd84a1b4d08e3377b4f6fec3a7e10b7de7e4c81 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal4630a041cc6330c0d3de6055dbd84a1b4d08e3377b4f6fec3a7e10b7de7e4c81 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.http-method','data' => ['method' => ''.e($request->method()).'']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::http-method'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['method' => ''.e($request->method()).'']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal4630a041cc6330c0d3de6055dbd84a1b4d08e3377b4f6fec3a7e10b7de7e4c81)): ?>
<?php $attributes = $__attributesOriginal4630a041cc6330c0d3de6055dbd84a1b4d08e3377b4f6fec3a7e10b7de7e4c81; ?>
<?php unset($__attributesOriginal4630a041cc6330c0d3de6055dbd84a1b4d08e3377b4f6fec3a7e10b7de7e4c81); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal4630a041cc6330c0d3de6055dbd84a1b4d08e3377b4f6fec3a7e10b7de7e4c81)): ?>
<?php $component = $__componentOriginal4630a041cc6330c0d3de6055dbd84a1b4d08e3377b4f6fec3a7e10b7de7e4c81; ?>
<?php unset($__componentOriginal4630a041cc6330c0d3de6055dbd84a1b4d08e3377b4f6fec3a7e10b7de7e4c81); ?>
<?php endif; ?>
        <div class="flex-1 text-sm font-light truncate text-neutral-950 dark:text-white">
            <span data-tippy-content="<?php echo e($request->fullUrl()); ?>">
                <?php echo e($request->fullUrl()); ?>

            </span>
        </div>
        <button
            x-cloak
            @click="copyToClipboard()"
            class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                "rounded-md w-6 h-6 flex flex-shrink-0 items-center justify-center cursor-pointer border transition-colors duration-200 ease-in-out",
                "bg-white/5 border-neutral-200 hover:bg-neutral-100 dark:bg-white/5 dark:border-white/10 dark:hover:bg-white/10",
            ]); ?>"
        >
            <?php if (isset($component)) { $__componentOriginal8ed599f137df9ecfe3925827c2e0e1455125f7665a815f0eaf8a27a2afc1e95a = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal8ed599f137df9ecfe3925827c2e0e1455125f7665a815f0eaf8a27a2afc1e95a = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.copy','data' => ['class' => 'w-3 h-3 text-neutral-400','xShow' => '!copied']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.copy'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-3 h-3 text-neutral-400','x-show' => '!copied']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal8ed599f137df9ecfe3925827c2e0e1455125f7665a815f0eaf8a27a2afc1e95a)): ?>
<?php $attributes = $__attributesOriginal8ed599f137df9ecfe3925827c2e0e1455125f7665a815f0eaf8a27a2afc1e95a; ?>
<?php unset($__attributesOriginal8ed599f137df9ecfe3925827c2e0e1455125f7665a815f0eaf8a27a2afc1e95a); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal8ed599f137df9ecfe3925827c2e0e1455125f7665a815f0eaf8a27a2afc1e95a)): ?>
<?php $component = $__componentOriginal8ed599f137df9ecfe3925827c2e0e1455125f7665a815f0eaf8a27a2afc1e95a; ?>
<?php unset($__componentOriginal8ed599f137df9ecfe3925827c2e0e1455125f7665a815f0eaf8a27a2afc1e95a); ?>
<?php endif; ?>
            <?php if (isset($component)) { $__componentOriginalff1497a8922457f6bd99a773b2dac243b39513e7fe68526735b594f3447a2704 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalff1497a8922457f6bd99a773b2dac243b39513e7fe68526735b594f3447a2704 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.check','data' => ['class' => 'w-3 h-3 text-emerald-500','xShow' => 'copied']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.check'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-3 h-3 text-emerald-500','x-show' => 'copied']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalff1497a8922457f6bd99a773b2dac243b39513e7fe68526735b594f3447a2704)): ?>
<?php $attributes = $__attributesOriginalff1497a8922457f6bd99a773b2dac243b39513e7fe68526735b594f3447a2704; ?>
<?php unset($__attributesOriginalff1497a8922457f6bd99a773b2dac243b39513e7fe68526735b594f3447a2704); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalff1497a8922457f6bd99a773b2dac243b39513e7fe68526735b594f3447a2704)): ?>
<?php $component = $__componentOriginalff1497a8922457f6bd99a773b2dac243b39513e7fe68526735b594f3447a2704; ?>
<?php unset($__componentOriginalff1497a8922457f6bd99a773b2dac243b39513e7fe68526735b594f3447a2704); ?>
<?php endif; ?>
        </button>
    </div>
</div>
<?php /**PATH d:\gitcode.com\origami\examples\laravel13\vendor\laravel\framework\src\illuminate\foundation\providers/../resources/exceptions/renderer/components/request-url.blade.php ENDPATH**/ ?>