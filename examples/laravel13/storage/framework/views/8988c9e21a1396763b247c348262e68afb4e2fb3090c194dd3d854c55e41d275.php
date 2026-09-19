<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'currentPageOptionProperty' => 'tableRecordsPerPage',
    'extremeLinks' => false,
    'paginator',
    'pageOptions' => [],
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
    'currentPageOptionProperty' => 'tableRecordsPerPage',
    'extremeLinks' => false,
    'paginator',
    'pageOptions' => [],
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
    use Filament\Support\View\SupportIconAlias;
    use Illuminate\Contracts\Pagination\CursorPaginator;
    use Illuminate\Pagination\LengthAwarePaginator;
    use Illuminate\Support\Number;

    $isRtl = __('filament-panels::layout.direction') === 'rtl';
    $isSimple = ! $paginator instanceof LengthAwarePaginator;
?>

<nav
    aria-label="<?php echo e(__('filament::components/pagination.label')); ?>"
    <?php echo e($attributes->class([
            'fi-pagination',
            'fi-simple' => $isSimple,
        ])); ?>

>
    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $paginator->onFirstPage()): ?>
        <?php
            if ($paginator instanceof CursorPaginator) {
                $wireClickAction = "setPage('{$paginator->previousCursor()->encode()}', '{$paginator->getCursorName()}')";
            } else {
                $wireClickAction = "previousPage('{$paginator->getPageName()}')";
            }
        ?>

        <?php if (isset($component)) { $__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.button.index','data' => ['color' => 'gray','rel' => 'prev','wire:click' => $wireClickAction,'wire:key' => $this->getId() . '.pagination.previous','class' => 'fi-pagination-previous-btn']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'gray','rel' => 'prev','wire:click' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($wireClickAction),'wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.pagination.previous'),'class' => 'fi-pagination-previous-btn']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

            <?php echo e(__('filament::components/pagination.actions.previous.label')); ?>

         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef)): ?>
<?php $attributes = $__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef; ?>
<?php unset($__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef)): ?>
<?php $component = $__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef; ?>
<?php unset($__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef); ?>
<?php endif; ?>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $isSimple): ?>
        <span class="fi-pagination-overview">
            <?php echo e(trans_choice(
                    'filament::components/pagination.overview',
                    $paginator->total(),
                    [
                        'first' => Number::format($paginator->firstItem() ?? 0),
                        'last' => Number::format($paginator->lastItem() ?? 0),
                        'total' => Number::format($paginator->total()),
                    ],
                )); ?>

        </span>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(count($pageOptions) > 1): ?>
        <div class="fi-pagination-records-per-page-select-ctn">
            <label class="fi-pagination-records-per-page-select fi-compact">
                <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                    <?php if (isset($component)) { $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.select','data' => ['wire:model.live' => $currentPageOptionProperty]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.select'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['wire:model.live' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($currentPageOptionProperty)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $pageOptions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $option): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                            <option value="<?php echo e($option); ?>">
                                <?php echo e($option === 'all' ? __('filament::components/pagination.fields.records_per_page.options.all') : $option); ?>

                            </option>
                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $attributes = $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $component = $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
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

                <span class="fi-sr-only">
                    <?php echo e(__('filament::components/pagination.fields.records_per_page.label')); ?>

                </span>
            </label>

            <label class="fi-pagination-records-per-page-select">
                <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => ['prefix' => __('filament::components/pagination.fields.records_per_page.label')]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['prefix' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament::components/pagination.fields.records_per_page.label'))]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                    <?php if (isset($component)) { $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.select','data' => ['wire:model.live' => $currentPageOptionProperty]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.select'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['wire:model.live' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($currentPageOptionProperty)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $pageOptions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $option): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                            <option value="<?php echo e($option); ?>">
                                <?php echo e($option === 'all' ? __('filament::components/pagination.fields.records_per_page.options.all') : $option); ?>

                            </option>
                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $attributes = $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $component = $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
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
            </label>
        </div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($paginator->hasMorePages()): ?>
        <?php
            if ($paginator instanceof CursorPaginator) {
                $wireClickAction = "setPage('{$paginator->nextCursor()->encode()}', '{$paginator->getCursorName()}')";
            } else {
                $wireClickAction = "nextPage('{$paginator->getPageName()}')";
            }
        ?>

        <?php if (isset($component)) { $__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.button.index','data' => ['color' => 'gray','rel' => 'next','wire:click' => $wireClickAction,'wire:key' => $this->getId() . '.pagination.next','class' => 'fi-pagination-next-btn']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::button'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'gray','rel' => 'next','wire:click' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($wireClickAction),'wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.pagination.next'),'class' => 'fi-pagination-next-btn']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

            <?php echo e(__('filament::components/pagination.actions.next.label')); ?>

         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef)): ?>
<?php $attributes = $__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef; ?>
<?php unset($__attributesOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef)): ?>
<?php $component = $__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef; ?>
<?php unset($__componentOriginal1cd039e01206aed786ce85171c75c0fa68208d7acbe1d87b78f326fa6a0386ef); ?>
<?php endif; ?>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if((! $isSimple) && $paginator->hasPages()): ?>
        <ol class="fi-pagination-items">
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $paginator->onFirstPage()): ?>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($extremeLinks): ?>
                    <?php if (isset($component)) { $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.pagination.item','data' => ['ariaLabel' => __('filament::components/pagination.actions.first.label'),'icon' => $isRtl ? Heroicon::ChevronDoubleRight : Heroicon::ChevronDoubleLeft,'iconAlias' => 
                            $isRtl
                            ? SupportIconAlias::PAGINATION_FIRST_BUTTON_RTL
                            : SupportIconAlias::PAGINATION_FIRST_BUTTON
                        ,'rel' => 'first','wire:click' => 'gotoPage(1, \'' . $paginator->getPageName() . '\')','wire:key' => $this->getId() . '.pagination.first']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::pagination.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['aria-label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament::components/pagination.actions.first.label')),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isRtl ? Heroicon::ChevronDoubleRight : Heroicon::ChevronDoubleLeft),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
                            $isRtl
                            ? SupportIconAlias::PAGINATION_FIRST_BUTTON_RTL
                            : SupportIconAlias::PAGINATION_FIRST_BUTTON
                        ),'rel' => 'first','wire:click' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute('gotoPage(1, \'' . $paginator->getPageName() . '\')'),'wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.pagination.first')]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $attributes = $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $component = $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <?php if (isset($component)) { $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.pagination.item','data' => ['ariaLabel' => __('filament::components/pagination.actions.previous.label'),'icon' => $isRtl ? Heroicon::ChevronRight : Heroicon::ChevronLeft,'iconAlias' => 
                        $isRtl
                        ? [
                            SupportIconAlias::PAGINATION_PREVIOUS_BUTTON_RTL,
                            SupportIconAlias::PAGINATION_PREVIOUS_BUTTON,
                        ]
                        : SupportIconAlias::PAGINATION_PREVIOUS_BUTTON
                    ,'rel' => 'prev','wire:click' => 'previousPage(\'' . $paginator->getPageName() . '\')','wire:key' => $this->getId() . '.pagination.previous']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::pagination.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['aria-label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament::components/pagination.actions.previous.label')),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isRtl ? Heroicon::ChevronRight : Heroicon::ChevronLeft),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
                        $isRtl
                        ? [
                            SupportIconAlias::PAGINATION_PREVIOUS_BUTTON_RTL,
                            SupportIconAlias::PAGINATION_PREVIOUS_BUTTON,
                        ]
                        : SupportIconAlias::PAGINATION_PREVIOUS_BUTTON
                    ),'rel' => 'prev','wire:click' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute('previousPage(\'' . $paginator->getPageName() . '\')'),'wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.pagination.previous')]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $attributes = $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $component = $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $paginator->render()->offsetGet('elements'); $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $element): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(is_string($element)): ?>
                    <?php if (isset($component)) { $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.pagination.item','data' => ['disabled' => true,'label' => $element]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::pagination.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['disabled' => true,'label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($element)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $attributes = $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $component = $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(is_array($element)): ?>
                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $element; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $page => $url): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                        <?php if (isset($component)) { $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.pagination.item','data' => ['active' => $page === $paginator->currentPage(),'ariaLabel' => trans_choice('filament::components/pagination.actions.go_to_page.label', $page, ['page' => Number::format($page)]),'label' => Number::format($page),'wire:click' => 'gotoPage(' . $page . ', \'' . $paginator->getPageName() . '\')','wire:key' => $this->getId() . '.pagination.' . $paginator->getPageName() . '.' . $page]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::pagination.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['active' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($page === $paginator->currentPage()),'aria-label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(trans_choice('filament::components/pagination.actions.go_to_page.label', $page, ['page' => Number::format($page)])),'label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(Number::format($page)),'wire:click' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute('gotoPage(' . $page . ', \'' . $paginator->getPageName() . '\')'),'wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.pagination.' . $paginator->getPageName() . '.' . $page)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $attributes = $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $component = $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>
                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($paginator->hasMorePages()): ?>
                <?php if (isset($component)) { $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.pagination.item','data' => ['ariaLabel' => __('filament::components/pagination.actions.next.label'),'icon' => $isRtl ? Heroicon::ChevronLeft : Heroicon::ChevronRight,'iconAlias' => 
                        $isRtl
                        ? [
                            SupportIconAlias::PAGINATION_NEXT_BUTTON_RTL,
                            SupportIconAlias::PAGINATION_NEXT_BUTTON,
                        ]
                        : SupportIconAlias::PAGINATION_NEXT_BUTTON
                    ,'rel' => 'next','wire:click' => 'nextPage(\'' . $paginator->getPageName() . '\')','wire:key' => $this->getId() . '.pagination.next']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::pagination.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['aria-label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament::components/pagination.actions.next.label')),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isRtl ? Heroicon::ChevronLeft : Heroicon::ChevronRight),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
                        $isRtl
                        ? [
                            SupportIconAlias::PAGINATION_NEXT_BUTTON_RTL,
                            SupportIconAlias::PAGINATION_NEXT_BUTTON,
                        ]
                        : SupportIconAlias::PAGINATION_NEXT_BUTTON
                    ),'rel' => 'next','wire:click' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute('nextPage(\'' . $paginator->getPageName() . '\')'),'wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.pagination.next')]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $attributes = $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $component = $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($extremeLinks): ?>
                    <?php if (isset($component)) { $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.pagination.item','data' => ['ariaLabel' => __('filament::components/pagination.actions.last.label'),'icon' => $isRtl ? Heroicon::ChevronDoubleLeft : Heroicon::ChevronDoubleRight,'iconAlias' => 
                            $isRtl
                            ? SupportIconAlias::PAGINATION_LAST_BUTTON_RTL
                            : SupportIconAlias::PAGINATION_LAST_BUTTON
                        ,'rel' => 'last','wire:click' => 'gotoPage(' . $paginator->lastPage() . ', \'' . $paginator->getPageName() . '\')','wire:key' => $this->getId() . '.pagination.last']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::pagination.item'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['aria-label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament::components/pagination.actions.last.label')),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isRtl ? Heroicon::ChevronDoubleLeft : Heroicon::ChevronDoubleRight),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
                            $isRtl
                            ? SupportIconAlias::PAGINATION_LAST_BUTTON_RTL
                            : SupportIconAlias::PAGINATION_LAST_BUTTON
                        ),'rel' => 'last','wire:click' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute('gotoPage(' . $paginator->lastPage() . ', \'' . $paginator->getPageName() . '\')'),'wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.pagination.last')]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $attributes = $__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__attributesOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121)): ?>
<?php $component = $__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121; ?>
<?php unset($__componentOriginalda20403f83ea19f113e3401e81e9bfdc73a0a93d854620b76ed04afd44955121); ?>
<?php endif; ?>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        </ol>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
</nav>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\support\resources\views/components/pagination/index.blade.php ENDPATH**/ ?>