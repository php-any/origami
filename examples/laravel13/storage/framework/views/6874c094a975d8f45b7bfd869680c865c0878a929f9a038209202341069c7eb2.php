<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames((['frame']));

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

foreach (array_filter((['frame']), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<div
    x-data="{
        expanded: <?php echo e($frame->isMain() ? 'true' : 'false'); ?>,
        hasCode: <?php echo e($frame->snippet() ? 'true' : 'false'); ?>

    }"
    class="group rounded-lg border border-neutral-200 dark:border-white/10 overflow-hidden shadow-xs"
    :class="{ 'dark:border-white/5': expanded }"
>
    <div
        class="flex h-11 items-center gap-3 bg-white pr-2.5 pl-4 overflow-x-auto dark:bg-white/3"
        :class="{
            'cursor-pointer hover:bg-white/50 dark:hover:bg-white/5 hover:[&_svg]:stroke-emerald-500': hasCode,
            'dark:bg-white/5 rounded-t-lg': expanded,
            'dark:bg-white/3 rounded-lg': !expanded
        }"
        @click="hasCode && (expanded = !expanded)"
    >
        
        <div class="flex size-3 items-center justify-center flex-shrink-0">
          <div
          class="size-2 rounded-full"
          :class="{
            'bg-rose-500 dark:bg-neutral-400': expanded,
            'bg-rose-200 dark:bg-neutral-700': !expanded
          }"
          ></div>
        </div>

        <div class="flex flex-1 items-center justify-between gap-6 min-w-0">
            <?php if (isset($component)) { $__componentOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.formatted-source','data' => ['frame' => $frame]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::formatted-source'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['frame' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($frame)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536)): ?>
<?php $attributes = $__attributesOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536; ?>
<?php unset($__attributesOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536)): ?>
<?php $component = $__componentOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536; ?>
<?php unset($__componentOriginal948595ed2658d3fcc312521d04de3dc0004c77099b5ce29166c0be6d38b78536); ?>
<?php endif; ?>
            <?php if (isset($component)) { $__componentOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.file-with-line','data' => ['frame' => $frame,'direction' => 'rtl']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::file-with-line'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['frame' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($frame),'direction' => 'rtl']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6)): ?>
<?php $attributes = $__attributesOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6; ?>
<?php unset($__attributesOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6)): ?>
<?php $component = $__componentOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6; ?>
<?php unset($__componentOriginalee29f095cb3b4b91dcd2323359b981aefdfd9fc5452702c066ec644d44f683a6); ?>
<?php endif; ?>
        </div>

        <div class="flex-shrink-0">
            <button
                x-cloak
                type="button"
                class="flex h-6 w-6 cursor-pointer items-center justify-center rounded-md dark:border dark:border-white/8 group-hover:text-blue-500 group-hover:dark:text-emerald-500"
                :class="{
                    'text-blue-500 dark:text-emerald-500 dark:bg-white/5': expanded,
                    'text-neutral-500 dark:text-neutral-500 dark:bg-white/3': !expanded,
                }"
            >
                <?php if (isset($component)) { $__componentOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.chevrons-down-up','data' => ['xShow' => 'expanded']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.chevrons-down-up'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-show' => 'expanded']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964)): ?>
<?php $attributes = $__attributesOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964; ?>
<?php unset($__attributesOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964)): ?>
<?php $component = $__componentOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964; ?>
<?php unset($__componentOriginal5708536dfe075f01ce7dc0fc14d8e85b02462403a9e2ebf760c8e661c4f07964); ?>
<?php endif; ?>
                <?php if (isset($component)) { $__componentOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.chevrons-up-down','data' => ['xShow' => '!expanded']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.chevrons-up-down'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-show' => '!expanded']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a)): ?>
<?php $attributes = $__attributesOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a; ?>
<?php unset($__attributesOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a)): ?>
<?php $component = $__componentOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a; ?>
<?php unset($__componentOriginalbc7a17fec9881fc0f91d2773f0ec7e2c1998c8fb3553d4758937523663095f9a); ?>
<?php endif; ?>
            </button>
        </div>
    </div>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($snippet = $frame->snippet()): ?>
        <?php if (isset($component)) { $__componentOriginalba47bda4e3e10bd57c46dc9ee9c8e310c198962b93d02887e2dc98fad88aa2a5 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalba47bda4e3e10bd57c46dc9ee9c8e310c198962b93d02887e2dc98fad88aa2a5 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.frame-code','data' => ['code' => $snippet,'highlightedLine' => $frame->line(),'xShow' => 'expanded','xCloak' => !$frame->isMain()]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::frame-code'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['code' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($snippet),'highlightedLine' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($frame->line()),'x-show' => 'expanded','x-cloak' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(!$frame->isMain())]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalba47bda4e3e10bd57c46dc9ee9c8e310c198962b93d02887e2dc98fad88aa2a5)): ?>
<?php $attributes = $__attributesOriginalba47bda4e3e10bd57c46dc9ee9c8e310c198962b93d02887e2dc98fad88aa2a5; ?>
<?php unset($__attributesOriginalba47bda4e3e10bd57c46dc9ee9c8e310c198962b93d02887e2dc98fad88aa2a5); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalba47bda4e3e10bd57c46dc9ee9c8e310c198962b93d02887e2dc98fad88aa2a5)): ?>
<?php $component = $__componentOriginalba47bda4e3e10bd57c46dc9ee9c8e310c198962b93d02887e2dc98fad88aa2a5; ?>
<?php unset($__componentOriginalba47bda4e3e10bd57c46dc9ee9c8e310c198962b93d02887e2dc98fad88aa2a5); ?>
<?php endif; ?>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
</div>
<?php /**PATH d:\gitcode.com\origami\examples\laravel13\vendor\laravel\framework\src\illuminate\foundation\providers/../resources/exceptions/renderer/components/frame.blade.php ENDPATH**/ ?>