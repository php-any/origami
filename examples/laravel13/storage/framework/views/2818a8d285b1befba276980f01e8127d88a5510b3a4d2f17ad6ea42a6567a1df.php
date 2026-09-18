<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames((['queries']));

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

foreach (array_filter((['queries']), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<div
    <?php echo e($attributes->merge(['class' => "flex flex-col gap-2.5 bg-neutral-50 dark:bg-white/1 border border-neutral-200 dark:border-neutral-800 rounded-xl p-2.5 shadow-xs"])); ?>

    x-data="{
        totalQueries: <?php echo e(min(count($queries), 100)); ?>,
        currentPage: 1,
        perPage: 10,
        get totalPages() {
            return Math.ceil(this.totalQueries / this.perPage);
        },
        get hasPrevious() {
            return this.currentPage > 1;
        },
        get hasNext() {
            return this.currentPage < this.totalPages;
        },
        goToPage(page) {
            if (page >= 1 && page <= this.totalPages) {
                this.currentPage = page;
            }
        },
        first() {
            this.currentPage = 1;
        },
        last() {
            this.currentPage = this.totalPages;
        },
        previous() {
            if (this.hasPrevious) {
                this.currentPage--;
            }
        },
        next() {
            if (this.hasNext) {
                this.currentPage++;
            }
        },
        get visiblePages() {
            const total = this.totalPages;
            const current = this.currentPage;
            const pages = [];

            if (total <= 7) {
                for (let i = 1; i <= total; i++) {
                    pages.push({ type: 'page', value: i });
                }
            } else {
                if (current <= 4) {
                    for (let i = 1; i <= 5; i++) {
                        pages.push({ type: 'page', value: i });
                    }
                    if (total > 6) {
                        pages.push({ type: 'ellipsis', value: '...', id: 'end' });
                        pages.push({ type: 'page', value: total });
                    }
                } else if (current > total - 4) {
                    pages.push({ type: 'page', value: 1 });
                    if (total > 6) {
                        pages.push({ type: 'ellipsis', value: '...', id: 'start' });
                    }
                    for (let i = Math.max(total - 4, 2); i <= total; i++) {
                        pages.push({ type: 'page', value: i });
                    }
                } else {
                    pages.push({ type: 'page', value: 1 });
                    pages.push({ type: 'ellipsis', value: '...', id: 'start' });
                    for (let i = current - 1; i <= current + 1; i++) {
                        pages.push({ type: 'page', value: i });
                    }
                    pages.push({ type: 'ellipsis', value: '...', id: 'end' });
                    pages.push({ type: 'page', value: total });
                }
            }
            return pages;
        }
    }"
>
    <div class="flex items-center justify-between p-2">
        <div class="flex items-center gap-2.5">
            <div class="bg-white dark:bg-neutral-800 border border-neutral-200 dark:border-white/5 rounded-md w-6 h-6 flex items-center justify-center p-1">
                <?php if (isset($component)) { $__componentOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.database','data' => ['class' => 'w-2.5 h-2.5 text-blue-500 dark:text-emerald-500']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.database'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-2.5 h-2.5 text-blue-500 dark:text-emerald-500']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5)): ?>
<?php $attributes = $__attributesOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5; ?>
<?php unset($__attributesOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5)): ?>
<?php $component = $__componentOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5; ?>
<?php unset($__componentOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5); ?>
<?php endif; ?>
            </div>
            <h3 class="text-base font-semibold">Queries</h3>
        </div>
        <div x-show="totalQueries > 0" class="text-sm text-neutral-500 dark:text-neutral-400 flex items-center gap-2">
            <span x-text="`${((currentPage - 1) * perPage) + 1}-${Math.min(currentPage * perPage, totalQueries)} of ${totalQueries}`"></span>
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(count($queries) > 100): ?>
                <?php if (isset($component)) { $__componentOriginalf8b629fa372ce55032251ea6839755640a8b889c551c243d9330e417082bce17 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalf8b629fa372ce55032251ea6839755640a8b889c551c243d9330e417082bce17 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.info','data' => ['class' => 'w-3 h-3 text-blue-500 dark:text-emerald-500','dataTippyContent' => 'Only the first 100 queries are shown']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.info'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-3 h-3 text-blue-500 dark:text-emerald-500','data-tippy-content' => 'Only the first 100 queries are shown']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalf8b629fa372ce55032251ea6839755640a8b889c551c243d9330e417082bce17)): ?>
<?php $attributes = $__attributesOriginalf8b629fa372ce55032251ea6839755640a8b889c551c243d9330e417082bce17; ?>
<?php unset($__attributesOriginalf8b629fa372ce55032251ea6839755640a8b889c551c243d9330e417082bce17); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalf8b629fa372ce55032251ea6839755640a8b889c551c243d9330e417082bce17)): ?>
<?php $component = $__componentOriginalf8b629fa372ce55032251ea6839755640a8b889c551c243d9330e417082bce17; ?>
<?php unset($__componentOriginalf8b629fa372ce55032251ea6839755640a8b889c551c243d9330e417082bce17); ?>
<?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        </div>
    </div>

    <div class="flex flex-col gap-1">
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__empty_1 = true; $__currentLoopData = array_slice($queries, 0, 100); $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $index => ['connectionName' => $connectionName, 'sql' => $sql, 'time' => $time]): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); $__empty_1 = false; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
        <div
            class="border border-neutral-200 dark:border-none bg-white dark:bg-white/[3%] rounded-md h-10 flex items-center justify-between gap-4 px-4 text-xs font-mono shadow-xs"
            x-show="Math.floor(<?php echo e($index); ?> / perPage) === (currentPage - 1)"
        >
            <div class="flex items-center gap-2 truncate">
                <div class="flex items-center gap-2">
                    <?php if (isset($component)) { $__componentOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.database','data' => ['class' => 'w-3 h-3 text-neutral-500 dark:text-neutral-400']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.database'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-3 h-3 text-neutral-500 dark:text-neutral-400']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5)): ?>
<?php $attributes = $__attributesOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5; ?>
<?php unset($__attributesOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5)): ?>
<?php $component = $__componentOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5; ?>
<?php unset($__componentOriginal573847b99661a6dabceba3eef70c9e9d91d2f1535b57e53b8f2e434a5d67c2e5); ?>
<?php endif; ?>
                    <span class="text-neutral-500 dark:text-neutral-400"><?php echo e($connectionName); ?></span>
                </div>
                <?php if (isset($component)) { $__componentOriginal60f51a3b3773d99d09e20b9ccfdbc33fd1b7917c98351e99d3b31ffcb65ce951 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal60f51a3b3773d99d09e20b9ccfdbc33fd1b7917c98351e99d3b31ffcb65ce951 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.syntax-highlight','data' => ['code' => $sql,'language' => 'sql','truncate' => true,'class' => 'min-w-0','dataTippyContent' => ''.e(nl2br($sql)).'']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::syntax-highlight'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['code' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($sql),'language' => 'sql','truncate' => true,'class' => 'min-w-0','data-tippy-content' => ''.e(nl2br($sql)).'']); ?>
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
            <div class="text-neutral-500 dark:text-neutral-200 text-right flex-shrink-0"><?php echo e($time); ?>ms</div>
        </div>
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); if ($__empty_1): ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
        <?php if (isset($component)) { $__componentOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal253ad76312637e58a35e9f310f446d8b3d72b69127f5c148b63a1ebac5476f56 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.empty-state','data' => ['message' => 'No queries executed']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::empty-state'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['message' => 'No queries executed']); ?>
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

    <!-- Pagination Controls -->
    <div x-cloak x-show="totalPages > 1" class="flex items-center justify-center gap-1 py-4 font-mono">
        <!-- First Button -->
        <button
            @click="first()"
            class="cursor-pointer flex items-center justify-center w-8 h-8 rounded-md transition-colors"
            :disabled="!hasPrevious"
            :class="hasPrevious ? 'text-neutral-500 dark:text-neutral-300 hover:bg-neutral-200 hover:dark:text-white hover:dark:bg-white/5' : 'text-neutral-600 cursor-not-allowed!'"
        >
            <?php if (isset($component)) { $__componentOriginalfd72ae4e99f1ca6fbe2f0c7a3d972487013c0993fb2090daf046fddcc913f7d4 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalfd72ae4e99f1ca6fbe2f0c7a3d972487013c0993fb2090daf046fddcc913f7d4 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.chevrons-left','data' => ['class' => 'w-3 h-3']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.chevrons-left'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-3 h-3']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalfd72ae4e99f1ca6fbe2f0c7a3d972487013c0993fb2090daf046fddcc913f7d4)): ?>
<?php $attributes = $__attributesOriginalfd72ae4e99f1ca6fbe2f0c7a3d972487013c0993fb2090daf046fddcc913f7d4; ?>
<?php unset($__attributesOriginalfd72ae4e99f1ca6fbe2f0c7a3d972487013c0993fb2090daf046fddcc913f7d4); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalfd72ae4e99f1ca6fbe2f0c7a3d972487013c0993fb2090daf046fddcc913f7d4)): ?>
<?php $component = $__componentOriginalfd72ae4e99f1ca6fbe2f0c7a3d972487013c0993fb2090daf046fddcc913f7d4; ?>
<?php unset($__componentOriginalfd72ae4e99f1ca6fbe2f0c7a3d972487013c0993fb2090daf046fddcc913f7d4); ?>
<?php endif; ?>
        </button>

        <!-- Previous Button -->
        <button
            @click="previous()"
            class="cursor-pointer flex items-center justify-center w-8 h-8 rounded-md transition-colors"
            :class="hasPrevious ? 'text-neutral-500 dark:text-neutral-300 hover:bg-neutral-200 hover:dark:text-white hover:dark:bg-white/5' : 'text-neutral-600 cursor-not-allowed!'"
            :disabled="!hasPrevious"
        >
            <?php if (isset($component)) { $__componentOriginalb93d5c5a8da092e663b028da7720c056884964a43470af349b33cf4186a1f01c = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalb93d5c5a8da092e663b028da7720c056884964a43470af349b33cf4186a1f01c = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.chevron-left','data' => ['class' => 'w-3 h-3']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.chevron-left'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-3 h-3']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalb93d5c5a8da092e663b028da7720c056884964a43470af349b33cf4186a1f01c)): ?>
<?php $attributes = $__attributesOriginalb93d5c5a8da092e663b028da7720c056884964a43470af349b33cf4186a1f01c; ?>
<?php unset($__attributesOriginalb93d5c5a8da092e663b028da7720c056884964a43470af349b33cf4186a1f01c); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalb93d5c5a8da092e663b028da7720c056884964a43470af349b33cf4186a1f01c)): ?>
<?php $component = $__componentOriginalb93d5c5a8da092e663b028da7720c056884964a43470af349b33cf4186a1f01c; ?>
<?php unset($__componentOriginalb93d5c5a8da092e663b028da7720c056884964a43470af349b33cf4186a1f01c); ?>
<?php endif; ?>
        </button>

        <!-- Page Numbers -->
        <template x-for="(page, index) in visiblePages" :key="`page-${page.type}-${page.value}-${page.id || index}`">
            <div>
                <template x-if="page.type === 'ellipsis'">
                    <span class="flex items-center justify-center w-8 h-8 text-neutral-500">...</span>
                </template>
                <template x-if="page.type === 'page'">
                    <button
                        @click="goToPage(page.value)"
                        class="cursor-pointer flex items-center justify-center w-8 h-8 rounded-md text-sm font-medium transition-colors"
                        :class="currentPage === page.value ? 'bg-blue-600 text-white' : 'text-neutral-500 dark:text-neutral-300 hover:bg-neutral-200 hover:dark:text-white hover:dark:bg-white/5'"
                        x-text="page.value"
                    ></button>
                </template>
            </div>
        </template>

        <!-- Next Button -->
        <button
            @click="next()"
            class="cursor-pointer flex items-center justify-center w-8 h-8 rounded-md transition-colors"
            :class="hasNext ? 'text-neutral-500 dark:text-neutral-300 hover:bg-neutral-200 hover:dark:text-white hover:dark:bg-white/5' : 'text-neutral-600 cursor-not-allowed!'"
            :disabled="!hasNext"
        >
            <?php if (isset($component)) { $__componentOriginal6b5f36e955da57c756e0aad56686c1d7b666be422794693fa0c7cabdda95b19d = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6b5f36e955da57c756e0aad56686c1d7b666be422794693fa0c7cabdda95b19d = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.chevron-right','data' => ['class' => 'w-3 h-3']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.chevron-right'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-3 h-3']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal6b5f36e955da57c756e0aad56686c1d7b666be422794693fa0c7cabdda95b19d)): ?>
<?php $attributes = $__attributesOriginal6b5f36e955da57c756e0aad56686c1d7b666be422794693fa0c7cabdda95b19d; ?>
<?php unset($__attributesOriginal6b5f36e955da57c756e0aad56686c1d7b666be422794693fa0c7cabdda95b19d); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal6b5f36e955da57c756e0aad56686c1d7b666be422794693fa0c7cabdda95b19d)): ?>
<?php $component = $__componentOriginal6b5f36e955da57c756e0aad56686c1d7b666be422794693fa0c7cabdda95b19d; ?>
<?php unset($__componentOriginal6b5f36e955da57c756e0aad56686c1d7b666be422794693fa0c7cabdda95b19d); ?>
<?php endif; ?>
        </button>

        <!-- Last Button -->
        <button
            @click="last()"
            class="cursor-pointer flex items-center justify-center w-8 h-8 rounded-md transition-colors"
            :class="hasNext ? 'text-neutral-500 dark:text-neutral-300 hover:bg-neutral-200 hover:dark:text-white hover:dark:bg-white/5' : 'text-neutral-600 cursor-not-allowed!'"
            :disabled="!hasNext"
        >
            <?php if (isset($component)) { $__componentOriginala24477544c1a5652ae6732734094e2ad8b3809e270f831c3afa696d12c1d8e73 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginala24477544c1a5652ae6732734094e2ad8b3809e270f831c3afa696d12c1d8e73 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.icons.chevrons-right','data' => ['class' => 'w-3 h-3']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::icons.chevrons-right'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'w-3 h-3']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginala24477544c1a5652ae6732734094e2ad8b3809e270f831c3afa696d12c1d8e73)): ?>
<?php $attributes = $__attributesOriginala24477544c1a5652ae6732734094e2ad8b3809e270f831c3afa696d12c1d8e73; ?>
<?php unset($__attributesOriginala24477544c1a5652ae6732734094e2ad8b3809e270f831c3afa696d12c1d8e73); ?>
<?php endif; ?>
<?php if (isset($__componentOriginala24477544c1a5652ae6732734094e2ad8b3809e270f831c3afa696d12c1d8e73)): ?>
<?php $component = $__componentOriginala24477544c1a5652ae6732734094e2ad8b3809e270f831c3afa696d12c1d8e73; ?>
<?php unset($__componentOriginala24477544c1a5652ae6732734094e2ad8b3809e270f831c3afa696d12c1d8e73); ?>
<?php endif; ?>
        </button>
    </div>
</div>
<?php /**PATH d:\gitcode.com\origami\examples\laravel13\vendor\laravel\framework\src\illuminate\foundation\providers/../resources/exceptions/renderer/components/query.blade.php ENDPATH**/ ?>